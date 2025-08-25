package compose

import (
	"gopkg.in/yaml.v3"
)

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
			},
		},
	}

	isHttps := true
	s := composeFIle.Services["traefik"]

	if isHttps {
		httpRule := "Host(`proxy.locom.self`)"
		s.LabelsNode = &yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				// proxy
				{Kind: yaml.ScalarNode, Value: "traefik.enable"},
				{Kind: yaml.ScalarNode, Value: "true"},

				// http
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik.rule"},
				{Kind: yaml.ScalarNode, Value: httpRule},
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik.entrypoints"},
				{Kind: yaml.ScalarNode, Value: "web"},

				// redirect
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik.middlewares"},
				{Kind: yaml.ScalarNode, Value: "redirect-to-https"},

				// https
				{Kind: yaml.ScalarNode, Value: "traefik.http.middlewares.redirect-to-https.redirectscheme.scheme"},
				{Kind: yaml.ScalarNode, Value: "https"},

				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik-secure.rule"},
				{Kind: yaml.ScalarNode, Value: httpRule},

				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik-secure.entrypoints"},
				{Kind: yaml.ScalarNode, Value: "websecure"},
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik-secure.service"},
				{Kind: yaml.ScalarNode, Value: "api@internal"},
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik-secure.tls"},
				{Kind: yaml.ScalarNode, Value: "true"},
			},
		}
	} else {
		s.LabelsNode = &yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				// proxy
				{Kind: yaml.ScalarNode, Value: "traefik.enable"},
				{Kind: yaml.ScalarNode, Value: "true"},

				// http
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik.rule"},
				{Kind: yaml.ScalarNode, Value: "Host(`proxy.locom.self`)"},
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik.service"},
				{Kind: yaml.ScalarNode, Value: "api@internal"},
				{Kind: yaml.ScalarNode, Value: "traefik.http.routers.traefik.entrypoints"},
				{Kind: yaml.ScalarNode, Value: "web"},
			},
		}
	}
	composeFIle.Services["traefik"] = s

	return composeFIle
}
