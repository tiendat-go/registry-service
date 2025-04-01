package http

import (
	"encoding/json"
	"net/http"

	"github.com/tiendat-go/registry-service/internal/service"
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
	serviceName := r.URL.Query().Get("ServiceName")
	address := r.URL.Query().Get("Address")
	success, _ := c.service.RegisterService(serviceName, address)
	w.Header().Set("Content-Type", "application/json")
	result := map[string]bool{"success": success}
	json.NewEncoder(w).Encode(result)
}

func (c *RegistryController) DeregisterService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("ServiceName")
	address := r.URL.Query().Get("Address")
	success, _ := c.service.DeregisterService(serviceName, address)
	w.Header().Set("Content-Type", "application/json")
	result := map[string]bool{"success": success}
	json.NewEncoder(w).Encode(result)
}

func (c *RegistryController) GetServices(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("ServiceName")
	addresses, _ := c.service.GetServices(serviceName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

func (c *RegistryController) GetRandService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("ServiceName")
	address, _ := c.service.GetRandService(serviceName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

func (c *RegistryController) Heartbeat(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("ServiceName")
	address := r.URL.Query().Get("Address")
	success, _ := c.service.Heartbeat(serviceName, address)
	w.Header().Set("Content-Type", "application/json")
	result := map[string]bool{"success": success}
	json.NewEncoder(w).Encode(result)
}
