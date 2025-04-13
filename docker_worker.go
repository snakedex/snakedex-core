package main

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
)

func docker(ctx context.Context) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	reader, err := cli.ImagePull(ctx, "hashicorp/terraform:1.11", image.PullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      "hashicorp/terraform:1.11",
		Entrypoint: []string{"sh", "-c"},
		Cmd:        []string{"terraform version"},
		Tty:        false,
	}, nil, nil, nil, "terraform-worker")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	//exec, err := cli.ContainerExecCreate(ctx, resp.ID, container.ExecOptions{
	//	Cmd: []string{"terraform version"},
	//})
	//
	//err = cli.ContainerExecStart(ctx, exec.ID, container.ExecStartOptions{})
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	err = cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
	if err != nil {
		panic(err)
	}
}
