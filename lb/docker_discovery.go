package main

import (
	"context"
	"log"
	"net/netip"
	"strconv"

	"github.com/moby/moby/client"
)

func discoverFromDocker(service string, existingBackends []Backend) []Backend {
	cli, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := cli.ContainerList(context.Background(), client.ContainerListOptions{})
	if err != nil {
		log.Println("[docker] Error listing containers:", err)
		return nil
	}

	var discoveredBackends []Backend
	for _, c := range containers.Items {
		if c.Labels["grpc.service"] != service {
			continue
		}

		port := c.Labels["grpc.port"]
		if port == "" {
			port = "50051"
		}

		weight := 1
		if w, ok := c.Labels["grpc.weight"]; ok {
			weight = toInt(w, 1)
		}

		inspect, err := cli.ContainerInspect(context.Background(), c.ID, client.ContainerInspectOptions{})
		if err != nil {
			log.Println("[docker] Error inspecting container:", err)
			continue
		}

		if inspect.Container.NetworkSettings != nil && inspect.Container.NetworkSettings.Networks != nil {
			for _, net := range inspect.Container.NetworkSettings.Networks {
				if net.IPAddress != (netip.Addr{}) {
					discoveredBackends = append(discoveredBackends, Backend{
						Addr:    net.IPAddress.String() + ":" + port,
						Healthy: false,
						Weight:  weight,
					})
				}
			}
		}
	}

	for i := range discoveredBackends {
		for _, existing := range existingBackends {
			if existing.Addr == discoveredBackends[i].Addr {
				discoveredBackends[i].Healthy = existing.Healthy
				break
			}
		}
	}

	log.Printf("[docker] Discovered %d backend(s): %v\n", len(discoveredBackends), discoveredBackends)
	return discoveredBackends
}

func toInt(str string, fallback int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return fallback
	}
	return val
}
