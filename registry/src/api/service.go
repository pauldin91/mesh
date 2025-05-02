package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Network string `json:"network"`
	Address string `json:"address"`
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		http.Error(w, "unable to load client", http.StatusInternalServerError)
		return
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		http.Error(w, "unable to list local containers", http.StatusInternalServerError)
		return
	}

	var response []ContainerResponse
	for _, ctr := range containers {
		if strings.HasPrefix(ctr.Status, "Up") {
			inspect, _ := apiClient.ContainerInspect(context.Background(), ctr.ID)

			response = append(response, ContainerResponse{
				Id:      ctr.ID[:12],
				Name:    ctr.Image,
				Network: ctr.HostConfig.NetworkMode,
				Address: inspect.NetworkSettings.Networks[ctr.HostConfig.NetworkMode].IPAddress,
			})
			fmt.Printf("%s %s (status: %s) host %s\n", ctr.ID, ctr.Image, ctr.Status, inspect.NetworkSettings.Gateway)
		}
	}
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "unable to obtain json response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
