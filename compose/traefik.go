package compose

func GetTraefikCompose(networkName string) ComposeFile {
	composeFIle := ComposeFile{
		Networks: map[string]ExternalNetwork{
			networkName: {External: true},
		},
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
					"./certs:/certs:ro",
				},
				Networks: []string{networkName},
				Labels:   nil,
			},
		},
	}

	isHttps := true
	s := composeFIle.Services["traefik"]
	if isHttps {
		s.Labels = map[string]string{
			// https
			// # Unsecured HTTP router (redirect to HTTPS)
			// "traefik.http.routers.traefik.rule": `&traefik-rule >-
			// 	Host(` + "`proxy.locom.self`)",
			"traefik.http.routers.traefik.rule":        "Host(`proxy.locom.self`)",
			"traefik.http.routers.traefik.entrypoints": "web",
			"traefik.http.routers.traefik.middlewares": "redirect-to-https",

			// # HTTPS router (certresolver + secure dashboard)
			"traefik.http.middlewares.redirect-to-https.redirectscheme.scheme": "https",
			// "traefik.http.routers.traefik-secure.rule":                         "*traefik-rule",
			"traefik.http.routers.traefik-secure.rule":        "Host(`proxy.locom.self`)",
			"traefik.http.routers.traefik-secure.entrypoints": "websecure",
			"traefik.http.routers.traefik-secure.service":     "api@internal",
			"traefik.http.routers.traefik-secure.tls":         "true",
		}
	} else {
		s.Labels = map[string]string{
			// http
			"traefik.enable":                           "true",
			"traefik.http.routers.traefik.rule":        "Host(`proxy.locom.self`)",
			"traefik.http.routers.traefik.service":     "api@internal",
			"traefik.http.routers.traefik.entrypoints": "web",
		}
	}
	composeFIle.Services["traefik"] = s

	return composeFIle
}
