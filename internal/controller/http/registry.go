package http

import (
	"net/http"

	"github.com/tiendat-go/common-service/model"
	"github.com/tiendat-go/registry-service/internal/core/service"
	"github.com/tiendat-go/registry-service/internal/core/utils"
)

type RegistryController struct {
	service *service.RegistryService
}

func NewRegistryController(service *service.RegistryService) *RegistryController {
	return &RegistryController{
		service: service,
	}
}

func (c *RegistryController) RegisterService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("serviceName")
	address := r.URL.Query().Get("address")
	success, err := c.service.RegisterService(serviceName, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, model.Response{
		Code:    model.Success,
		Success: success,
		Message: model.SuccessMsg,
	})
}

func (c *RegistryController) DeregisterService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("serviceName")
	address := r.URL.Query().Get("address")
	success, err := c.service.DeregisterService(serviceName, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, model.Response{
		Code:    model.Success,
		Success: success,
		Message: model.SuccessMsg,
	})
}

func (c *RegistryController) GetServices(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("serviceName")
	addresses, err := c.service.GetServices(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, model.Response{
		Code:    model.Success,
		Data:    addresses,
		Success: true,
		Message: model.SuccessMsg,
	})
}

func (c *RegistryController) GetRandService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("serviceName")
	address, err := c.service.GetRandService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, model.Response{
		Code:    model.Success,
		Data:    address,
		Success: true,
		Message: model.SuccessMsg,
	})
}

func (c *RegistryController) Heartbeat(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("serviceName")
	address := r.URL.Query().Get("address")
	success, err := c.service.Heartbeat(serviceName, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, model.Response{
		Code:    model.Success,
		Data:    success,
		Success: true,
		Message: model.SuccessMsg,
	})
}
