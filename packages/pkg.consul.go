package packages

import (
	"errors"
	"os"

	"github.com/hashicorp/consul/api"
)

type (
	IConsul interface {
		ServiceRegister(req *api.AgentServiceRegistration) error
		ServiceDeregister(serviceId string) error
		CheckRegister(check *api.AgentCheckRegistration) error
		CheckDeregister(serviceId string) error
		Services() (map[string]*api.AgentService, error)
		Service(serviceId string, q *api.QueryOptions) (*api.AgentService, *api.QueryMeta, error)
		HealthCheck(serviceName string, q *api.QueryOptions) (api.HealthChecks, *api.QueryMeta, error)
	}

	consul struct {
		client *api.Client
	}
)

func NewConsul() (IConsul, error) {
	config := new(api.Config)
	if err := validateConsulConfig(config); err != nil {
		return nil, err
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consul{client: client}, nil
}

func (h *consul) ServiceRegister(req *api.AgentServiceRegistration) error {
	return h.client.Agent().ServiceRegister(req)
}

func (h *consul) ServiceDeregister(serviceId string) error {
	return h.client.Agent().ServiceDeregister(serviceId)
}

func (h *consul) CheckRegister(req *api.AgentCheckRegistration) error {
	return h.client.Agent().CheckRegister(req)
}

func (h *consul) CheckDeregister(healthId string) error {
	return h.client.Agent().CheckDeregister(healthId)
}

func (h *consul) Services() (map[string]*api.AgentService, error) {
	return h.client.Agent().Services()
}

func (h *consul) Service(serviceId string, q *api.QueryOptions) (*api.AgentService, *api.QueryMeta, error) {
	return h.client.Agent().Service(serviceId, q)
}

func (h *consul) HealthCheck(serviceName string, q *api.QueryOptions) (api.HealthChecks, *api.QueryMeta, error) {
	return h.client.Health().Checks(serviceName, q)
}

func validateConsulConfig(config *api.Config) error {
	discoveryAddressVal, discoveryAddressOk := os.LookupEnv("DISCOVERY_ADDRESS")
	discoveryDataCenterVal, discoveryDataCenterOk := os.LookupEnv("DISCOVERY_DATA_CENTER")
	discoveryTokenVal, discoveryTokenOk := os.LookupEnv("DISCOVERY_TOKEN")

	if !discoveryAddressOk {
		return errors.New("Consul address is required")
	} else if !discoveryDataCenterOk {
		return errors.New("Consul datacenter is required")
	} else if !discoveryTokenOk {
		return errors.New("Consul token is required")
	}

	secure := true
	if os.Getenv("GO_ENV") != "development" {
		secure = false
	}

	config.Address = discoveryAddressVal
	config.Datacenter = discoveryDataCenterVal
	config.Datacenter = discoveryTokenVal
	config.TLSConfig = api.TLSConfig{InsecureSkipVerify: secure}

	return nil
}
