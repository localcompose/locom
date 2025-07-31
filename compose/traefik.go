package compose

func GetTraefikCompose(network string) map[string]interface{} {
	return map[string]interface{}{
		"networks": map[string]interface{}{
			network: map[string]interface{}{
				"external": true,
			},
		},
		"services": map[string]interface{}{
			"traefik": map[string]interface{}{
				"image":          "traefik:v2.10",
				"container_name": "traefik",
				"restart":        "unless-stopped",
				"command": []string{
					"--api.dashboard=true",
					"--api.insecure=true",
					"--providers.docker=true",
					"--providers.docker.exposedbydefault=false",
					"--entrypoints.web.address=:80",
					"--entrypoints.websecure.address=:443",
					"--providers.file.directory=/etc/traefik/dynamic",
					"--providers.file.watch=true",
				},
				"ports": []string{
					"80:80",
					"443:443",
					"8080:8080",
				},
				"volumes": []string{
					"/var/run/docker.sock:/var/run/docker.sock:ro",
					"./config:/etc/traefik/dynamic",
				},
				"networks": []string{
					network,
				},
				"labels": map[string]string{
					"traefik.enable":                           "true",
					"traefik.http.routers.traefik.rule":        "Host(`proxy.locom.self`)",
					"traefik.http.routers.traefik.service":     "api@internal",
					"traefik.http.routers.traefik.entrypoints": "web",
				},
			},
		},
	}
}
