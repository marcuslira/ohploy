package lib

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

// ContainerMgmt defines a manament object
type ContainerMgmt struct {
	config Config
	cli    *client.Client
}

// NewContainerMgmt creates a new ContainerMgmt instance
func NewContainerMgmt(config Config) (*ContainerMgmt, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	mgmt := new(ContainerMgmt)
	mgmt.config = config
	mgmt.cli = cli

	return mgmt, nil
}

// DeployContainer pulls and deploys the image
func (c *ContainerMgmt) DeployContainer() error {

	err := c.pullImage(c.cli, c.config.Container.ImageName)
	if err != nil {
		return err
	}

	ids, err := c.listContainersByImage(c.cli, c.config.Container.ImageName)
	if err != nil {
		return err
	}

	if len(ids) > 0 {
		err = c.stopContainer(c.cli, ids[0])

		if err != nil {
			return err
		}
	}

	_, err = c.startContainer(c.cli, c.config.Container.ImageName)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContainerMgmt) stopContainer(cli *client.Client, contID string) error {
	fmt.Println("ohploy: stopping container...")
	err := cli.ContainerStop(context.Background(), contID, nil)
	if err != nil {
		return err
	}

	return err
}
func (c *ContainerMgmt) listContainersByImage(cli *client.Client, image string) ([]string, error) {
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

func (c *ContainerMgmt) pullImage(cli *client.Client, image string) error {
	fmt.Println("ohploy: pulling new image...")

	authConfig := types.AuthConfig{
		Username: c.config.DeployServer.RegistryUser,
		Password: c.config.DeployServer.RegistryPass,
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

func (c *ContainerMgmt) startContainer(cli *client.Client, image string) (string, error) {
	fmt.Println("ohploy: starting container...")

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: c.config.DeployServer.HostPort,
	}

	containerPort, err := nat.NewPort("tcp", c.config.Container.ContainerPort)

	if err != nil {
		return "", err
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
			Env:   c.config.Container.ContainerEnv,
		},
		&container.HostConfig{
			PortBindings:  portBinding,
			RestartPolicy: container.RestartPolicy{Name: c.config.Container.RestartPolicy},
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
