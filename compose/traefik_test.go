package compose_test

import (
	"testing"

	"github.com/localcompose/locom/compose"
)

func TestGetTraefikCompose(t *testing.T) {
	compose := compose.GetTraefikCompose("locom-net")

	if _, ok := compose["services"]; !ok {
		t.Fatal("expected services key in compose")
	}

	networks := compose["networks"].(map[string]interface{})
	if _, ok := networks["locom-net"]; !ok {
		t.Error("expected network 'locom-net' to be defined")
	}
}
