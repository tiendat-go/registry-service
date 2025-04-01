package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	frt "github.com/tiendat-go/common-service/utils/format"
	httpController "github.com/tiendat-go/registry-service/internal/controller/http"
	"github.com/tiendat-go/registry-service/internal/core/service"
)

func main() {
	port := frt.GetString(os.Getenv("SERVICE_PORT"), "9999")
	service := service.NewRegistryService()
	server := httpController.NewRegistryController(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", server.RegisterService)
	mux.HandleFunc("/deregister", server.DeregisterService)
	mux.HandleFunc("/getservices", server.GetServices)
	mux.HandleFunc("/getrandservice", server.GetRandService)
	mux.HandleFunc("/heartbeat", server.Heartbeat)

	log.Printf("🔥 Service Registry is running on port:%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))
}
