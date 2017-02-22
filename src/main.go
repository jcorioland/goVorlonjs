package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

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
	var imageTag = "vorlonjs/dashboard:0.3.0"
	var vorlonjsPort = uint32(1337)
	var serviceName = r.URL.Query().Get("serviceName")
	if len(strings.TrimSpace(serviceName)) == 0 {
		fmt.Fprintln(w, "Service name cannot be empty")
		return
	}

	var randomPort = uint32(random(5000, 10000))
	log.Printf("New random port has been generated: %d\r\n", randomPort)

	result := createDockerService(imageTag, serviceName, vorlonjsPort, randomPort, "vorlonjs")
	log.Printf("New Vorlonjs container has been created: ID = %s\r\n", result.ID)

	fmt.Fprintf(w, "Vorlonjs is running at http://localhost:%d", randomPort)
}

// createDockerService creates a new Docker service in the Swarm cluster
func createDockerService(imageTag string, serviceName string, targetPort uint32, publishedPort uint32, networkName string) types.ServiceCreateResponse {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	var serviceSpec = swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: serviceName,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: imageTag,
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{
				swarm.PortConfig{
					PublishedPort: publishedPort,
					TargetPort:    targetPort,
					Protocol:      swarm.PortConfigProtocolTCP,
					PublishMode:   swarm.PortConfigPublishModeIngress,
				},
			},
		},
		Networks: []swarm.NetworkAttachmentConfig{
			swarm.NetworkAttachmentConfig{
				Target: networkName,
			},
		},
	}

	result, err := cli.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	return result
}

// random generates a random number between two range
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// pullDockerImage pulls an image using its tag
func pullDockerImage(imageTag string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cli.ImagePull(context.Background(), imageTag, types.ImagePullOptions{})
}
