package deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"

	"log"
)

// TaskDefinition get a current task definition
func (d *Deploy) TaskDefinition() (*ecs.TaskDefinition, error) {
	taskArn, err := d.Service()
	if err != nil {
		return nil, err
	}

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(*taskArn),
	}
	resp, err := d.awsECS.DescribeTaskDefinition(params)
	if err != nil {
		return nil, err
	}

	return resp.TaskDefinition, nil
}

// Service get target service
func (d *Deploy) Service() (*string, error) {
	params := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(d.name),
		},
		Cluster: aws.String(d.cluster),
	}
	resp, err := d.awsECS.DescribeServices(params)
	if err != nil {
		return nil, err
	}

	return resp.Services[0].TaskDefinition, nil
}

// RegisterTaskDefinition register new task definition if needed
func (d *Deploy) RegisterTaskDefinition(baseDefinition *ecs.TaskDefinition) (*ecs.TaskDefinition, error) {
	var containerDefinitions []*ecs.ContainerDefinition
	for _, c := range baseDefinition.ContainerDefinitions {
		newDefinition, err := d.NewContainerDefinition(c)
		if err != nil {
			return nil, err
		}
		containerDefinitions = append(containerDefinitions, newDefinition)
	}
	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: containerDefinitions,
		Family:               baseDefinition.Family,
		NetworkMode:          baseDefinition.NetworkMode,
		PlacementConstraints: baseDefinition.PlacementConstraints,
		TaskRoleArn:          baseDefinition.TaskRoleArn,
		Volumes:              baseDefinition.Volumes,
	}

	log.Printf("[INFO] new task definition: %+v\n", params)
	resp, err := d.awsECS.RegisterTaskDefinition(params)
	if err != nil {
		return nil, err
	}

	return resp.TaskDefinition, nil
}

// NewContainerDefinition update image tag in a given container definition.
// If the container definition is not target container, return givien definition.
func (d *Deploy) NewContainerDefinition(baseDefinition *ecs.ContainerDefinition) (*ecs.ContainerDefinition, error) {
	if d.image == nil {
		return baseDefinition, nil
	}
	baseImage, _, err := divideImageAndTag(*baseDefinition.Image)
	if err != nil {
		return nil, err
	}
	if *d.image != *baseImage {
		return baseDefinition, nil
	}
	imageWithTag := (*d.image) + ":" + (*d.tag)
	baseDefinition.Image = &imageWithTag
	return baseDefinition, nil
}
