package main

import (
	"context"
	"fmt"
	"io"
	"os"

	dt "github.com/docker/docker/api/types"
	df "github.com/docker/docker/api/types/filters"
	d "github.com/docker/docker/client"
	flag "github.com/spf13/pflag"
)

var (
	ctx = context.Background()
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprint(os.Stderr, "I need exactly *1* argument to work.\n")
		os.Exit(1)
	}
	requestedContainer := flag.Arg(0)

	docker, err := d.NewEnvClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to docker: %s\n", err.Error())
		os.Exit(2)
	}
	defer docker.Close()
	pong, err := docker.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not ping docker: %v\n", err)
		os.Exit(3)
	}
	fmt.Printf("Connected to Docker v %s\n", pong.APIVersion)

	listFilters := df.NewArgs()
	listFilters.Add("name", requestedContainer)
	containers, err := docker.ContainerList(ctx, dt.ContainerListOptions{
		Filters: listFilters,
		All:     true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not check if container is running: %s", err.Error())
		os.Exit(4)
	}
	if len(containers) == 1 {
		go StreamContainerLogs(os.Stdout, docker, containers[0].ID, requestedContainer)
	}
	eventFilters := df.NewArgs()
	eventFilters.Add("container", requestedContainer)
	msgC, errC := docker.Events(ctx, dt.EventsOptions{
		Filters: eventFilters,
	})
	for {
		select {
		case msg := <-msgC:
			switch msg.Status {
			case "start":
				containerName := msg.Actor.Attributes["name"]
				fmt.Printf("%s>> started. stream log...\n\n", requestedContainer)
				go StreamContainerLogs(os.Stdout, docker, msg.ID, containerName)
			default:
			}

		case err := <-errC:
			fmt.Fprintf(os.Stderr, "received err: %#v\n", err)
			os.Exit(4)
		}
	}
}

func StreamContainerLogs(to *os.File, docker *d.Client, containerID, containerName string) {
	rc, err := docker.ContainerLogs(ctx, containerID, dt.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
	})
	defer rc.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error trying to access logs of %s/%s", containerName, containerID)
		return
	}
	io.Copy(os.Stdout, rc)
	fmt.Printf("\n%s>> stream ended.\n", containerName)
}
