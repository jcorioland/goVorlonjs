package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// createDockerService creates a new Docker service in the Swarm cluster
func createDockerService(
	imageTag string,
	serviceName string,
	targetPort uint32,
	networkName string,
	environmentVariables []string,
	labels map[string]string) (types.ServiceCreateResponse, error) {

	// create a Docker client
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// create the Docker service specifications
	var serviceSpec = swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   serviceName,
			Labels: labels,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: imageTag,
				Env:   environmentVariables,
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{
				swarm.PortConfig{
					//PublishedPort: publishedPort,
					TargetPort:  targetPort,
					Protocol:    swarm.PortConfigProtocolTCP,
					PublishMode: swarm.PortConfigPublishModeIngress,
				},
			},
		},
		Networks: []swarm.NetworkAttachmentConfig{
			swarm.NetworkAttachmentConfig{
				Target: networkName,
			},
		},
	}

	return cli.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions{})
}

// remove a Docker service
func removeDockerService(nameOrIdentifier string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	rmError := cli.ServiceRemove(context.Background(), nameOrIdentifier)
	return rmError
}

func isServiceRunning(serviceName string) bool {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	serviceNameFilter := filters.NewArgs()
	serviceNameFilter.Add("name", serviceName)

	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filters: serviceNameFilter})
	if err != nil {
		panic(err)
	}

	return len(services) != 0
}

// pullDockerImage pulls an image using its tag
func pullDockerImage(imageTag string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cli.ImagePull(context.Background(), imageTag, types.ImagePullOptions{})
}
