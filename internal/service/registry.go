package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	pb "github.com/tiendat-go/proto-service/gen/registry/v1"
)

type ServiceInstance struct {
	Address   string
	LastCheck time.Time
}

type RegistryServer struct {
	pb.UnimplementedDiscoveryServiceServer
	services map[string]map[string]*ServiceInstance
	mu       sync.RWMutex
}

func NewRegistryServer() *RegistryServer {
	server := &RegistryServer{
		services: make(map[string]map[string]*ServiceInstance),
	}
	go server.cleanup(1*time.Second, 3*time.Second)
	return server
}

func (s *RegistryServer) cleanup(interval time.Duration, timeout time.Duration) {
	for {
		time.Sleep(interval)
		s.mu.Lock()
		for service, instances := range s.services {
			for addr, instance := range instances {
				if time.Since(instance.LastCheck) > timeout {
					log.Printf("❌ Removing inactive service: %s at %s", service, addr)
					delete(s.services[service], addr)
				}
			}
		}
		s.mu.Unlock()
	}
}

func (s *RegistryServer) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[req.ServiceName]; !exists {
		s.services[req.ServiceName] = make(map[string]*ServiceInstance)
	}

	s.register(req.ServiceName, req.Address)
	return &pb.RegisterServiceResponse{Success: true}, nil
}

func (s *RegistryServer) DeregisterService(ctx context.Context, req *pb.DeregisterServiceRequest) (*pb.DeregisterServiceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if instances, exists := s.services[req.ServiceName]; exists {
		if _, found := instances[req.Address]; found {
			delete(instances, req.Address)
			log.Printf("🛑 Deregistered: %s at %s", req.ServiceName, req.Address)
			return &pb.DeregisterServiceResponse{Success: true}, nil
		}
	}

	return &pb.DeregisterServiceResponse{Success: false}, fmt.Errorf("service instance not found: %s at %s", req.ServiceName, req.Address)
}

func (s *RegistryServer) GetServices(ctx context.Context, req *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := []string{}
	if instances, exists := s.services[req.ServiceName]; exists {
		for addr := range instances {
			addresses = append(addresses, addr)
		}
	}

	return &pb.GetServicesResponse{Addresses: addresses}, nil
}

func (s *RegistryServer) GetRandService(ctx context.Context, req *pb.GetRandServiceRequest) (*pb.GetRandServiceResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := []string{}
	if instances, exists := s.services[req.ServiceName]; exists {
		for addr := range instances {
			addresses = append(addresses, addr)
		}
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no available instances for service: %s", req.ServiceName)
	}

	return &pb.GetRandServiceResponse{Address: addresses[rand.Intn(len(addresses))]}, nil
}

func (s *RegistryServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if instances, exists := s.services[req.ServiceName]; exists {
		if instance, found := instances[req.Address]; found {
			instance.LastCheck = time.Now()
			log.Printf("💓 Heartbeat received: %s at %s\n", req.ServiceName, req.Address)
			return &pb.HeartbeatResponse{Success: true}, nil
		}
	}

	s.register(req.ServiceName, req.Address)
	return &pb.HeartbeatResponse{Success: false}, nil
}

func (s *RegistryServer) register(serviceName, address string) {
	if _, exists := s.services[serviceName]; !exists {
		s.services[serviceName] = make(map[string]*ServiceInstance)
	}
	s.services[serviceName][address] = &ServiceInstance{
		Address:   address,
		LastCheck: time.Now(),
	}
	log.Printf("🔥 Registered: %s at %s\n", serviceName, address)
}
