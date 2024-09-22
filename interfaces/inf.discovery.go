package interfaces

import (
	"net/http"

	dto "github.com/restuwahyu13/discovery-api/dtos"
	"github.com/restuwahyu13/discovery-api/helpers"
)

type (
	IHandlerDiscovery interface {
		Register(rw http.ResponseWriter, r *http.Request)
		Deregister(rw http.ResponseWriter, r *http.Request)
		CheckRegister(rw http.ResponseWriter, r *http.Request)
		CheckDeregister(rw http.ResponseWriter, r *http.Request)
		ListDiscovery(rw http.ResponseWriter, r *http.Request)
		DetailDiscovery(rw http.ResponseWriter, r *http.Request)
		HealthCheck(rw http.ResponseWriter, r *http.Request)
	}

	IServiceDiscovery interface {
		Register(req *helpers.Request[dto.DRegisterBody], res *helpers.Response) error
		Deregister(req *helpers.Request[dto.DDeregisterParam], res *helpers.Response) error
		CheckRegister(req *helpers.Request[dto.DCheckRegisterBody], res *helpers.Response) error
		CheckDeregister(req *helpers.Request[dto.DCCheckDeregisterParam], res *helpers.Response) error
		ListDiscovery(req *helpers.Request[http.Request], res *helpers.Response) error
		DetailDiscovery(req *helpers.Request[dto.DDiscoveryParam], res *helpers.Response) error
		HealthCheck(req *helpers.Request[dto.DHealthCheckParam], res *helpers.Response) error
	}
)
