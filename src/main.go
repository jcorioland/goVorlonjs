package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var vorlonjsImageTag = "vorlonjs/dashboard:0.5.4"

func main() {
	// get image version from environment variable
	imageVersion := os.Getenv("VORLONJS_DOCKER_IMAGE_VERSION")
	if len(strings.TrimSpace(imageVersion)) > 0 {
		vorlonjsImageTag = imageVersion
	}

	// handle the create vorlon service action
	http.HandleFunc("/govorlonjs/api/create", CreateVorlonInstance)

	// handle the remove vorlon service action
	http.HandleFunc("/govorlonjs/api/remove", RemoveVorlonInstance)

	// start the http server
	log.Printf("The Vorlon.js API has started on the port %d", 82)
	log.Fatal(http.ListenAndServe(":82", nil))
}

// CreateVorlonInstance creates a new service that runs a Vorlonjs Docker container on a Swarm cluster
func CreateVorlonInstance(w http.ResponseWriter, r *http.Request) {
	var vorlonjsPort = uint32(1337)
	var serviceName = r.URL.Query().Get("serviceName")

	// if the service name has not been specified
	if len(strings.TrimSpace(serviceName)) == 0 {
		// return HTTP 400 -> BAD REQUEST
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Service name cannot be empty")
		return
	}

	// generate a random port
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

	result, err := createDockerService(vorlonjsImageTag, serviceName, vorlonjsPort, randomPort, "vorlonjs", env, labels)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Something went wrong with your request: " + err.Error())
		return
	}

	log.Printf("New Vorlonjs container has been created: ID = %s\r\n", result.ID)

	// return HTTP 201 -> CREATED
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Vorlonjs is running at /"+serviceName)
}

// RemoveVorlonInstance removes a Vorlonjs service that is running in the Swarm cluster
func RemoveVorlonInstance(w http.ResponseWriter, r *http.Request) {
	var serviceName = r.URL.Query().Get("serviceName")

	// if the service name has not been specified
	if len(strings.TrimSpace(serviceName)) == 0 {
		// return HTTP 400 -> BAD REQUEST
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Service name cannot be empty")
		return
	}

	// remove the service
	err := removeDockerService(serviceName)
	if err != nil {
		// return HTTP 400 -> BAD REQUEST
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Something went wrong with your request: "+err.Error())
		return
	}

	// return HTTP 200 -> OK
	w.WriteHeader(http.StatusOK)
}
