package grpc

import (
	"context"

	pb "github.com/tiendat-go/proto-service/gen/registry/v1"
	"github.com/tiendat-go/registry-service/internal/core/service"
)

type RegistryController struct {
	pb.UnimplementedDiscoveryServiceServer
	service *service.RegistryService
}

func NewRegistryController(service *service.RegistryService) *RegistryController {
	return &RegistryController{
		service: service,
	}
}

func (c *RegistryController) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	success, err := c.service.RegisterService(req.ServiceName, req.Address)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterServiceResponse{Success: success}, nil
}

func (c *RegistryController) DeregisterService(ctx context.Context, req *pb.DeregisterServiceRequest) (*pb.DeregisterServiceResponse, error) {
	success, err := c.service.DeregisterService(req.ServiceName, req.Address)
	if err != nil {
		return nil, err
	}
	return &pb.DeregisterServiceResponse{Success: success}, nil
}

func (c *RegistryController) GetServices(ctx context.Context, req *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	addresses, err := c.service.GetServices(req.ServiceName)
	if err != nil {
		return nil, err
	}
	return &pb.GetServicesResponse{Addresses: *addresses}, nil
}

func (c *RegistryController) GetRandService(ctx context.Context, req *pb.GetRandServiceRequest) (*pb.GetRandServiceResponse, error) {
	address, err := c.service.GetRandService(req.ServiceName)
	if err != nil {
		return nil, err
	}
	return &pb.GetRandServiceResponse{Address: *address}, nil
}

func (c *RegistryController) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	success, err := c.service.Heartbeat(req.ServiceName, req.Address)
	if err != nil {
		return nil, err
	}
	return &pb.HeartbeatResponse{Success: success}, nil
}
