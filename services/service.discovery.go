package services

import (
	"net/http"

	"github.com/hashicorp/consul/api"

	"github.com/restuwahyu13/discovery-api/configs"
	dto "github.com/restuwahyu13/discovery-api/dtos"
	"github.com/restuwahyu13/discovery-api/helpers"
	inf "github.com/restuwahyu13/discovery-api/interfaces"
	"github.com/restuwahyu13/discovery-api/packages"
)

type ServiceDiscovery struct {
	Env    *configs.Environtment
	Consul packages.IConsul
}

func NewServiceDiscovery(options *ServiceDiscovery) inf.IServiceDiscovery {
	return &ServiceDiscovery{Env: options.Env, Consul: options.Consul}
}

func (s *ServiceDiscovery) Register(req *helpers.Request[dto.DRegisterBody], res *helpers.Response) error {
	if err := s.Consul.ServiceRegister(&req.Body.AgentServiceRegistration); err != nil {
		res.StatCode = http.StatusPreconditionFailed
		res.ErrCode = "DISCOVERY_FAILED"
		res.ErrMsg = "Registered a new service to consul failed"

		defer packages.Logrus("error", err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMsg = "Registered a new service to consul success"

	return nil
}

func (s *ServiceDiscovery) Deregister(req *helpers.Request[dto.DDeregisterParam], res *helpers.Response) error {
	if err := s.Consul.ServiceDeregister(req.Param.ServiceID); err != nil {
		res.StatCode = http.StatusPreconditionFailed
		res.ErrCode = "DISCOVERY_FAILED"
		res.ErrMsg = "Deregister a old service to consul failed"

		defer packages.Logrus("error", err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMsg = "Deregister a old service to consul success"

	return nil
}

func (s *ServiceDiscovery) CheckRegister(req *helpers.Request[dto.DCheckRegisterBody], res *helpers.Response) error {
	err := s.Consul.CheckRegister(&req.Body.AgentCheckRegistration)
	if err != nil {
		res.StatCode = http.StatusNotFound
		res.ErrCode = "DISCOVERY_FAILED"
		res.ErrMsg = "Registered a new health to consul failed"

		defer packages.Logrus("error", err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMsg = "Registered a new health to consul success"

	return nil
}

func (s *ServiceDiscovery) CheckDeregister(req *helpers.Request[dto.DCCheckDeregisterParam], res *helpers.Response) error {
	err := s.Consul.CheckDeregister(req.Param.HealthID)
	if err != nil {
		res.StatCode = http.StatusNotFound
		res.ErrCode = "DISCOVERY_NOT_FOUND"
		res.ErrMsg = "Deregister a old health to consul failed"

		defer packages.Logrus("error", err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMsg = "Deregister a old health to consul success"

	return nil
}

func (s *ServiceDiscovery) ListDiscovery(req *helpers.Request[http.Request], res *helpers.Response) error {
	services, err := s.Consul.Services()
	if err != nil {
		res.StatCode = http.StatusNotFound
		res.ErrCode = "DISCOVERY_NOT_FOUND"
		res.ErrMsg = "No service registered in consul"

		defer packages.Logrus("error", err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMsg = "Success"
	res.Data = services

	return nil
}

func (s *ServiceDiscovery) DetailDiscovery(req *helpers.Request[dto.DDiscoveryParam], res *helpers.Response) error {
	services, _, err := s.Consul.Service(req.Param.ServiceID, &api.QueryOptions{Datacenter: s.Env.DATA_CENTER, Token: s.Env.TOKEN})

	if err != nil {
		res.StatCode = http.StatusNotFound
		res.ErrCode = "DISCOVERY_NOT_FOUND"
		res.ErrMsg = "No service registered in consul"

		defer packages.Logrus("error", err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMsg = "Success"
	res.Data = services

	return nil
}

func (s *ServiceDiscovery) HealthCheck(req *helpers.Request[dto.DHealthCheckParam], res *helpers.Response) error {
	services, _, err := s.Consul.HealthCheck(req.Param.ServiceName, &api.QueryOptions{Datacenter: s.Env.DATA_CENTER, Token: s.Env.TOKEN})

	if err != nil {
		res.StatCode = http.StatusNotFound
		res.ErrCode = "DISCOVERY_NOT_FOUND"
		res.ErrMsg = "No service registered in consul"

		defer packages.Logrus("error", err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMsg = "Success"
	res.Data = services

	return nil
}
