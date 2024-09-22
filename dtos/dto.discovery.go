package dtos

import (
	"github.com/hashicorp/consul/api"
)

type (
	DRegisterBody struct {
		api.AgentServiceRegistration
	}

	DDeregisterParam struct {
		ServiceID string `json:"serviceId"`
	}

	DCheckRegisterBody struct {
		api.AgentCheckRegistration
	}

	DCCheckDeregisterParam struct {
		HealthID string `json:"healthId"`
	}

	DDiscoveryParam struct {
		ServiceID string `json:"serviceId"`
	}

	DHealthCheckParam struct {
		ServiceName string `json:"serviceName"`
	}
)
