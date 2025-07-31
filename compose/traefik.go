package compose

func GetTraefikCompose(networkName string) ComposeFile {
	return ComposeFile{
		Services: map[string]Service{
			"traefik": {
				Image:         "traefik:v2.10",
				ContainerName: "traefik",
				Restart:       "unless-stopped",
				Command: []string{
					"--api.dashboard=true",
					"--api.insecure=true",
					"--providers.docker=true",
					"--providers.docker.exposedbydefault=false",
					"--entrypoints.web.address=:80",
					"--entrypoints.websecure.address=:443",
					"--providers.file.directory=/etc/traefik/dynamic",
					"--providers.file.watch=true",
				},
				Ports: []string{
					"80:80",
					"443:443",
					"8080:8080",
				},
				Volumes: []string{
					"/var/run/docker.sock:/var/run/docker.sock:ro",
					"./config:/etc/traefik/dynamic",
				},
				Networks: []string{networkName},
				Labels: map[string]string{
					"traefik.enable":                           "true",
					"traefik.http.routers.traefik.rule":        "Host(`proxy.locom.self`)",
					"traefik.http.routers.traefik.service":     "api@internal",
					"traefik.http.routers.traefik.entrypoints": "web",
				},
			},
		},
		Networks: map[string]ExternalNetwork{
			networkName: {External: true},
		},
	}
}
