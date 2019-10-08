package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func deployContainer(image string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	err = pullImage(cli, image)
	if err != nil {
		return err
	}

	ids, err := listContainersByImage(cli, image)
	if err != nil {
		return err
	}

	if len(ids) > 0 {
		err = stopContainer(cli, ids[0])

		if err != nil {
			return err
		}
	}

	_, err = startContainer(cli, image)
	if err != nil {
		return err
	}

	return nil
}

func stopContainer(cli *client.Client, contID string) error {
	fmt.Println("ohploy: stopping container...")
	err := cli.ContainerStop(context.Background(), contID, nil)
	if err != nil {
		return err
	}

	return err
}
func listContainersByImage(cli *client.Client, image string) ([]string, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	var result []string
	if len(containers) > 0 {
		for _, container := range containers {
			if strings.Contains(container.Image, image) {
				result = append(result, container.ID)
			}
		}
	}

	return result, nil

}

func pullImage(cli *client.Client, image string) error {
	fmt.Println("ohploy: pulling new image...")

	authConfig := types.AuthConfig{
		Username: "",
		Password: "",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, out)
	return nil
}

func startContainer(cli *client.Client, image string) (string, error) {
	fmt.Println("ohploy: starting container...")

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "5000",
	}

	containerPort, err := nat.NewPort("tcp", "5000")

	if err != nil {
		return "", err
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "")

	if err != nil {
		return "", err
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	if err != nil {
		return "", err
	}

	return cont.ID, nil
}

func main() {
	fmt.Println("ohploy: deploying container...")
	imageName := "docker.io/marcuslira/aspiratracker:latest"

	err := deployContainer(imageName)
	if err != nil {
		fmt.Printf("ohploy: Error - Deploying container: %v\n", err)
	}

	fmt.Println("ohploy: Done.")
}
