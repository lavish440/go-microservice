package main

import (
	"context"
	"log"
	"net/netip"

	"github.com/moby/moby/client"
)

func discoverFromDocker(service string) []Backend {
	cli, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := cli.ContainerList(context.Background(), client.ContainerListOptions{})
	if err != nil {
		log.Println("[docker] list error:", err)
		return nil
	}

	var backends []Backend

	for _, c := range containers.Items {
		if c.Labels["grpc.service"] != service {
			continue
		}

		port := c.Labels["grpc.port"]
		if port == "" {
			port = "50051"
		}

		inspect, err := cli.ContainerInspect(context.Background(), c.ID, client.ContainerInspectOptions{})
		if err != nil {
			log.Println("[docker] inspect error:", err)
			continue
		}

		if inspect.Container.NetworkSettings != nil && inspect.Container.NetworkSettings.Networks != nil {
			for _, net := range inspect.Container.NetworkSettings.Networks {
				if net.IPAddress != (netip.Addr{}) {
					backends = append(backends, Backend{
						Addr:    net.IPAddress.String() + ":" + port,
						Healthy: false,
					})
				}
			}
		}
	}

	return backends
}
