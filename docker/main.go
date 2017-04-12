package main

import (
	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	"fmt"
	"strings"
)

// State is the container internal state
type State int

const (
	// Created state
	Created State = iota
	// Running state
	Running
	// Stopped state
	Stopped
	// Destroyed state
	Destroyed
)

// DockerClient internal data structure
type DockerClient struct {
	client *docker.Client
}

func main() {

	dc, err := docker.NewClient("unix:///var/run/docker.sock")

	if err != nil {
		fmt.Errorf("failed to create a docker client: %v", err)
	}

	id := "a51ba7d410134226a577f0f9517e5b8eff6f3a43763f31ffe7f4b1300d7218ca"
	//ctx := context.Background()
	container, err := dc.InspectContainer(id);

	if err != nil {
		fmt.Errorf("failed to inspect container with context: %v", err)
	}

	log.Info(container.Name) // k8s_user_user-468431046-ktrt4_default

	pod, err := parsePodName(container.Name); // user-468431046-ktrt4
	if err != nil {
		fmt.Errorf("failed to ParseDockerName: %v", err)
	}

	log.Printf("pod is %s", pod)

}
/*
 convert container name to pod name:
  k8s_user_user-468431046-ktrt4_default -> user-468431046-ktrt4
 */

// Unpacks a container name, returning the pod full name.
// If we are unable to parse the name, an error is returned.
// https://github.com/kubernetes/kubernetes/blob/cda109d22480bc6dea3c06cef21bd4c4fca6fca2/pkg/kubelet/dockertools/docker.go
func parsePodName(name string) (podName string, err error) {
	// For some reason docker appears to be appending '/' to names.
	// If it's there, strip it.
	var containerNamePrefix = "k8s"
	name = strings.TrimPrefix(name, "/")
	parts := strings.Split(name, "_")
	if len(parts) == 0 || parts[0] != containerNamePrefix {
		err = fmt.Errorf("failed to parse Docker container name %q into parts", name)
		return "", err
	}
	if len(parts) < 6 {
		// We have at least 5 fields.  We may have more in the future.
		// Anything with less fields than this is not something we can
		// manage.
		log.Warningf("found a container with the %q prefix, but too few fields (%d): %q", containerNamePrefix, len(parts), name)
		err = fmt.Errorf("Docker container name %q has less parts than expected %v", name, parts)
		return "", err
	}

	nameParts := strings.Split(parts[1], ".")
	containerName := nameParts[0]
	log.Printf(containerName)
	if len(nameParts) > 1 {
		if err != nil {
			log.Warningf("invalid container hash %q in container %q", nameParts[1], name)
		}
	}

	podFullName := parts[2]

	return podFullName, nil
}

func (c *DockerClient) getContainerState(containerID string) (*docker.State, error) {
	dockerContainer, err := c.getContainer(containerID)
	if dockerContainer == nil {
		return nil, err
	}
	return &dockerContainer.State, nil
}

func (c *DockerClient) getContainer(containerID string) (*docker.Container, error) {
	dockerContainer, err := c.client.InspectContainer(containerID)
	if err != nil {
		if _, ok := err.(*docker.NoSuchContainer); !ok {
			return nil, err
		}
		return nil, nil
	}
	return dockerContainer, nil
}