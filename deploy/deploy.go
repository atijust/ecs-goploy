package deploy

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

// Deploy have target ecs information
type Deploy struct {
	awsECS  *ecs.ECS
	cluster string
	name    string
	image   *string
	tag     *string
}

// NewDeploy return a new Deploy struct, and initialize aws ecs api client
func NewDeploy(cluster, name, profile, region, imageWithTag string) *Deploy {
	awsECS := ecs.New(session.New(), newConfig(profile, region))
	var image, tag *string
	if len(imageWithTag) > 0 {
		var err error
		image, tag, err = divideImageAndTag(imageWithTag)
		if err != nil {
			log.Fatalf("[ERROR] Can not parse --image parameter: %+v\n", err)
		}
	}
	return &Deploy{
		awsECS,
		cluster,
		name,
		image,
		tag,
	}
}

// Deploy run deploy commands
func (d *Deploy) Deploy() {
	task, err := d.TaskDefinition()
	if err != nil {
		log.Fatalf("[ERROR] Can not get current task definition: %+v\n", err)
	}
	_, err = d.RegisterTaskDefinition(task)
	if err != nil {
		log.Fatalf("[ERROR] Can not regist new task definition: %+v\n", err)
	}

}

func divideImageAndTag(imageWithTag string) (*string, *string, error) {
	res := strings.Split(imageWithTag, ":")
	if len(res) >= 3 {
		return nil, nil, errors.New("image format is wrong.")
	}
	return &res[0], &res[1], nil

}
