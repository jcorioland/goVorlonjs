package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func main() {
	http.HandleFunc("/govorlonjs/api/createVorlonContainer", CreateVorlonContainer)
	log.Fatal(http.ListenAndServe(":82", nil))
}

// CreateVorlonContainer creates a new service that runs a Vorlonjs Docker container on a Swarm cluster
func CreateVorlonContainer(w http.ResponseWriter, r *http.Request) {
	var imageTag = "vorlonjs/dashboard:0.5.4"
	var vorlonjsPort = uint32(1337)
	var serviceName = r.URL.Query().Get("serviceName")
	if len(strings.TrimSpace(serviceName)) == 0 {
		fmt.Fprintln(w, "Service name cannot be empty")
		return
	}

	var randomPort = uint32(random(5000, 10000))
	log.Printf("New random port has been generated: %d\r\n", randomPort)

	labels := map[string]string{
		"com.df.notify":      "true",
		"com.df.distribute":  "true",
		"com.df.servicePath": "/" + serviceName,
		"com.df.port":        strconv.Itoa(int(vorlonjsPort)),
	}

	env := []string{
		"BASE_URL=/" + serviceName,
	}

	result := createDockerService(imageTag, serviceName, vorlonjsPort, randomPort, "vorlonjs", env, labels)
	log.Printf("New Vorlonjs container has been created: ID = %s\r\n", result.ID)

	fmt.Fprintf(w, "Vorlonjs is running at http://localhost:%d", randomPort)
}

// createDockerService creates a new Docker service in the Swarm cluster
func createDockerService(imageTag string, serviceName string, targetPort uint32, publishedPort uint32, networkName string, environmentVariables []string, labels map[string]string) types.ServiceCreateResponse {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

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
