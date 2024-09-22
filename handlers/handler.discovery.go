package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"

	"github.com/restuwahyu13/discovery-api/dtos"
	"github.com/restuwahyu13/discovery-api/helpers"
	"github.com/restuwahyu13/discovery-api/interfaces"
)

type HandlerDiscovery struct {
	Service interfaces.IServiceDiscovery
}

func NewHandlerDiscovery(options *HandlerDiscovery) interfaces.IHandlerDiscovery {
	return &HandlerDiscovery{Service: options.Service}
}

func (h *HandlerDiscovery) Register(rw http.ResponseWriter, r *http.Request) {
	var (
		req *helpers.Request[dtos.DRegisterBody] = new(helpers.Request[dtos.DRegisterBody])
		res *helpers.Response                    = new(helpers.Response)
	)

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
		res.StatCode = http.StatusUnprocessableEntity
		res.ErrCode = "REQUEST_FAILED"
		res.ErrMsg = err.Error()

		helpers.ApiResponse(rw, res)
		return
	}

	if err := h.Service.Register(req, res); err != nil {
		helpers.ApiResponse(rw, res)
		return
	}

	helpers.ApiResponse(rw, res)
}

func (h *HandlerDiscovery) Deregister(rw http.ResponseWriter, r *http.Request) {
	var (
		req *helpers.Request[dtos.DDeregisterParam] = new(helpers.Request[dtos.DDeregisterParam])
		res *helpers.Response                       = new(helpers.Response)
	)

	req.Param.ServiceID = chi.URLParam(r, "serviceId")

	if err := h.Service.Deregister(req, res); err != nil {
		helpers.ApiResponse(rw, res)
		return
	}

	helpers.ApiResponse(rw, res)
}

func (h *HandlerDiscovery) CheckRegister(rw http.ResponseWriter, r *http.Request) {
	var (
		req *helpers.Request[dtos.DCheckRegisterBody] = new(helpers.Request[dtos.DCheckRegisterBody])
		res *helpers.Response                         = new(helpers.Response)
	)

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
		res.StatCode = http.StatusUnprocessableEntity
		res.ErrCode = "REQUEST_FAILED"
		res.ErrMsg = err.Error()

		helpers.ApiResponse(rw, res)
		return
	}

	if err := h.Service.CheckRegister(req, res); err != nil {
		helpers.ApiResponse(rw, res)
		return
	}

	helpers.ApiResponse(rw, res)
}

func (h *HandlerDiscovery) CheckDeregister(rw http.ResponseWriter, r *http.Request) {
	var (
		req *helpers.Request[dtos.DCCheckDeregisterParam] = new(helpers.Request[dtos.DCCheckDeregisterParam])
		res *helpers.Response                             = new(helpers.Response)
	)

	req.Param.HealthID = chi.URLParam(r, "checkId")

	if err := h.Service.CheckDeregister(req, res); err != nil {
		helpers.ApiResponse(rw, res)
		return
	}

	helpers.ApiResponse(rw, res)
}

func (h *HandlerDiscovery) ListDiscovery(rw http.ResponseWriter, r *http.Request) {
	var (
		req *helpers.Request[http.Request] = new(helpers.Request[http.Request])
		res *helpers.Response              = new(helpers.Response)
	)

	if err := h.Service.ListDiscovery(req, res); err != nil {
		helpers.ApiResponse(rw, res)
		return
	}

	helpers.ApiResponse(rw, res)
}

func (h *HandlerDiscovery) DetailDiscovery(rw http.ResponseWriter, r *http.Request) {
	var (
		req *helpers.Request[dtos.DDiscoveryParam] = new(helpers.Request[dtos.DDiscoveryParam])
		res *helpers.Response                      = new(helpers.Response)
	)

	req.Param.ServiceID = chi.URLParam(r, "serviceId")

	if err := h.Service.DetailDiscovery(req, res); err != nil {
		helpers.ApiResponse(rw, res)
		return
	}

	helpers.ApiResponse(rw, res)
}

func (h *HandlerDiscovery) HealthCheck(rw http.ResponseWriter, r *http.Request) {
	var (
		req *helpers.Request[dtos.DHealthCheckParam] = new(helpers.Request[dtos.DHealthCheckParam])
		res *helpers.Response                        = new(helpers.Response)
	)

	req.Param.ServiceName = chi.URLParam(r, "serviceName")

	if err := h.Service.HealthCheck(req, res); err != nil {
		helpers.ApiResponse(rw, res)
		return
	}

	helpers.ApiResponse(rw, res)
}
