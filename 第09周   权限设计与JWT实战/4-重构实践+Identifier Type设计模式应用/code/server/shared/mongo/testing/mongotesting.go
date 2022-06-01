package mongotesting

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	image         = "mongo:latest"
	containerPort = "27017/tcp"
)

func RunWithMongoInDocker(m *testing.M, mongoURI *string) int {
	c, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image:        image,
		ExposedPorts: nat.PortSet{"27017/tcp": {}},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0", //会自动挑选本机一个端口，防止手写的端口27018被占用
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	containerID := resp.ID
	defer func() {
		err = c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err) // panic能让defer执行完，不能用log
		}
	}()
	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	inspRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	HostPort := inspRes.NetworkSettings.Ports[containerPort][0]

	*mongoURI = fmt.Sprintf("mongodb://%s:%s", HostPort.HostIP, HostPort.HostPort)
	return m.Run()
}
