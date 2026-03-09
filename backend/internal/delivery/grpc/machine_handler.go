package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/usecase"
	pb "github.com/natapage/test_assignment/backend/pkg/gen/goldex/v1"
)

type MachineHandler struct {
	pb.UnimplementedMachineServiceServer
	uc *usecase.MachineUseCase
}

func NewMachineHandler(uc *usecase.MachineUseCase) *MachineHandler {
	return &MachineHandler{uc: uc}
}

func (h *MachineHandler) ListMachines(ctx context.Context, _ *pb.ListMachinesRequest) (*pb.ListMachinesResponse, error) {
	machines, err := h.uc.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list machines: %v", err)
	}
	resp := &pb.ListMachinesResponse{}
	for i := range machines {
		resp.Machines = append(resp.Machines, machineToProto(&machines[i]))
	}
	return resp, nil
}

func (h *MachineHandler) GetMachine(ctx context.Context, req *pb.GetMachineRequest) (*pb.GetMachineResponse, error) {
	m, err := h.uc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "machine not found: %v", err)
	}
	return &pb.GetMachineResponse{Machine: machineToProto(m)}, nil
}

func (h *MachineHandler) CreateMachine(ctx context.Context, req *pb.CreateMachineRequest) (*pb.CreateMachineResponse, error) {
	m := &domain.Machine{
		Name:         req.GetName(),
		SerialNumber: req.GetSerialNumber(),
		Enabled:      req.GetEnabled(),
		LocationID:   req.LocationId,
	}
	created, err := h.uc.Create(ctx, m)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create machine: %v", err)
	}
	return &pb.CreateMachineResponse{Machine: machineToProto(created)}, nil
}

func (h *MachineHandler) UpdateMachine(ctx context.Context, req *pb.UpdateMachineRequest) (*pb.UpdateMachineResponse, error) {
	m := &domain.Machine{
		ID:           req.GetId(),
		Name:         req.GetName(),
		SerialNumber: req.GetSerialNumber(),
		Enabled:      req.GetEnabled(),
	}
	updated, err := h.uc.Update(ctx, m)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update machine: %v", err)
	}
	return &pb.UpdateMachineResponse{Machine: machineToProto(updated)}, nil
}

func (h *MachineHandler) DeleteMachine(ctx context.Context, req *pb.DeleteMachineRequest) (*emptypb.Empty, error) {
	if err := h.uc.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete machine: %v", err)
	}
	return &emptypb.Empty{}, nil
}
