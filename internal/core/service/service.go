package service

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type ServiceInstance struct {
	Address   string
	LastCheck time.Time
}

type RegistryService struct {
	services map[string]map[string]*ServiceInstance
	mu       sync.RWMutex
}

func NewRegistryService() *RegistryService {
	service := &RegistryService{
		services: make(map[string]map[string]*ServiceInstance),
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			for service, instances := range service.services {
				log.Printf("%v: %v", service, len(instances))
			}
		}
	}()
	go service.cleanup(1*time.Second, 3*time.Second)
	return service
}

func (s *RegistryService) cleanup(interval time.Duration, timeout time.Duration) {
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

func (s *RegistryService) register(serviceName, address string) {
	if _, exists := s.services[serviceName]; !exists {
		s.services[serviceName] = make(map[string]*ServiceInstance)
	}
	s.services[serviceName][address] = &ServiceInstance{
		Address:   address,
		LastCheck: time.Now(),
	}
	log.Printf("🔥 Registered: %s at %s\n", serviceName, address)
}

func (s *RegistryService) RegisterService(serviceName, address string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[serviceName]; !exists {
		s.services[serviceName] = make(map[string]*ServiceInstance)
	}

	s.register(serviceName, address)
	return true, nil
}

func (s *RegistryService) DeregisterService(serviceName, address string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if instances, exists := s.services[serviceName]; exists {
		if _, found := instances[address]; found {
			delete(instances, address)
			log.Printf("🛑 Deregistered: %s at %s", serviceName, address)
			return true, nil
		}
	}

	return false, fmt.Errorf("service instance not found: %s at %s", serviceName, address)
}

func (s *RegistryService) GetServices(serviceName string) (*[]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := []string{}
	if instances, exists := s.services[serviceName]; exists {
		for addr := range instances {
			addresses = append(addresses, addr)
		}
	}

	return &addresses, nil
}

func (s *RegistryService) GetRandService(serviceName string) (*string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := []string{}
	if instances, exists := s.services[serviceName]; exists {
		for addr := range instances {
			addresses = append(addresses, addr)
		}
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no available instances for service: %s", serviceName)
	}

	serviceAddress := addresses[rand.Intn(len(addresses))]
	return &serviceAddress, nil
}

func (s *RegistryService) Heartbeat(serviceName, address string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if instances, exists := s.services[serviceName]; exists {
		if instance, found := instances[address]; found {
			instance.LastCheck = time.Now()
			// log.Printf("💓 Heartbeat received: %s at %s\n", serviceName, address)
			return true, nil
		}
	}

	s.register(serviceName, address)
	return false, nil
}
