package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func main() {
	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/createVorlonContainer", CreateVorlonContainer)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HandleHome handles the root route of the applicaiton
func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// CreateVorlonContainer creates a new service that runs a Vorlonjs Docker container on a Swarm cluster
func CreateVorlonContainer(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	var imageName = "vorlonjs/dashboard:0.3.0"
	PullDockerImage(imageName)

	var serviceName = r.URL.Query().Get("serviceName")

	var serviceSpec = swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: serviceName,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: imageName,
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{
				swarm.PortConfig{
					TargetPort:    1337,
					PublishedPort: 1337,
					Protocol:      swarm.PortConfigProtocolTCP,
					PublishMode:   swarm.PortConfigPublishModeIngress,
				},
			},
		},
	}

	cli.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions{})
}

// PullDockerImage pulls an image using its tag
func PullDockerImage(imageTag string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cli.ImagePull(context.Background(), imageTag, types.ImagePullOptions{})
}
