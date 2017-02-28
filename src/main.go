package main

import (
	"encoding/json"
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
	http.HandleFunc("/api/instance/create", CreateVorlonInstance)

	// handle the remove vorlon service action
	http.HandleFunc("/api/instance/remove", RemoveVorlonInstance)

	// start the http server
	log.Printf("The Vorlon.js API has started on the port %d", 82)
	log.Fatal(http.ListenAndServe(":82", nil))
}

// CreateVorlonInstance creates a new service that runs a Vorlonjs Docker container on a Swarm cluster
func CreateVorlonInstance(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Usage: POST /api/instance/create {\"serviceName\": \"SERVICE_NAME\"}")
		return
	}

	// create a JSON decoder to parse the request body
	decoder := json.NewDecoder(r.Body)
	var requestBody VorlonInstanceRequestBody
	err := decoder.Decode(&requestBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Something went wrong with your request: "+err.Error())
		return
	}

	var vorlonjsPort = uint32(1337)
	var networkName = "vorlonjs"

	// if the service name has not been specified
	if len(strings.TrimSpace(requestBody.ServiceName)) == 0 {
		// return HTTP 400 -> BAD REQUEST
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Service name cannot be empty")
		return
	}

	labels := map[string]string{
		"com.df.notify":      "true",
		"com.df.distribute":  "true",
		"com.df.servicePath": "/" + requestBody.ServiceName,
		"com.df.port":        strconv.Itoa(int(vorlonjsPort)),
	}

	env := []string{
		"BASE_URL=/" + requestBody.ServiceName,
	}

	result, err := createDockerService(vorlonjsImageTag, requestBody.ServiceName, vorlonjsPort, networkName, env, labels)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Something went wrong with your request: "+err.Error())
		return
	}

	// return HTTP 201 -> CREATED
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Vorlonjs is running at /"+requestBody.ServiceName)

	log.Printf("New Vorlonjs container has been created: ID = %s\r\n", result.ID)
}

// RemoveVorlonInstance removes a Vorlonjs service that is running in the Swarm cluster
func RemoveVorlonInstance(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Usage: POST /api/instance/remove {\"serviceName\": \"SERVICE_NAME\"}")
		return
	}

	// create a JSON decoder to parse the request body
	decoder := json.NewDecoder(r.Body)
	var requestBody VorlonInstanceRequestBody
	err := decoder.Decode(&requestBody)

	// if the service name has not been specified
	if len(strings.TrimSpace(requestBody.ServiceName)) == 0 {
		// return HTTP 400 -> BAD REQUEST
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Service name cannot be empty")
		return
	}

	// remove the service
	err = removeDockerService(requestBody.ServiceName)
	if err != nil {
		// return HTTP 400 -> BAD REQUEST
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Something went wrong with your request: "+err.Error())
		return
	}

	// return HTTP 200 -> OK
	w.WriteHeader(http.StatusOK)
	log.Printf("The Vorlonjs instance %s has been removed\r\n", requestBody.ServiceName)
}
