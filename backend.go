package main

import (
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"log"
	"net/http"
	"os"
)

// ENV VARIABLE WITH ROOT PATH (/path/to/dest/)
base_path := os.Getenv("BASE_PATH")

type M map[string]interface{}

func listHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(err)
	}

	var out_containers []M
	for _, container := range containers {
		out_container := M{
			"id":    container.ID,
			"name":  container.Names[0],
			"state": container.State,
		}
		out_containers = append(out_containers, out_container)
	}
	output_json, err := json.Marshal(out_containers)
	if err != nil {
		//error
	}
	message := string(output_json)
	w.Write([]byte(message))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}
	id := keys[0]

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	err = client.StartContainer(id, nil)
	if err != nil {
		//panic(err)
		message = "{\"result\":\"error\"}"
	} else {
		message = "{\"result\":\"started\"}"
	}
	w.Write([]byte(message))
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}
	id := keys[0]

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	//Timeout of 5 seconds
	err = client.StopContainer(id, 5)
	if err != nil {
		//panic(err)
		message = "{\"result\":\"error\"}"
	} else {
		message = "{\"result\":\"stopped\"}"
	}
	w.Write([]byte(message))
}

func redirect(w http.ResponseWriter, r *http.Request) {

    http.Redirect(w, r, base_path+"static/", 301)
}

func main() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle(base_path+"static/", http.StripPrefix(base_path+"static/", fs))
	http.HandleFunc(base_path+"list", listHandler)
	http.HandleFunc(base_path+"start", startHandler)
	http.HandleFunc(base_path+"stop", stopHandler)
	http.HandleFunc(base_path, redirect)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
