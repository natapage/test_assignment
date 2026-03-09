package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/natapage/test_assignment/backend/internal/usecase"
	pb "github.com/natapage/test_assignment/backend/pkg/gen/goldex/v1"
)

type MovementHandler struct {
	pb.UnimplementedMovementServiceServer
	uc *usecase.MovementUseCase
}

func NewMovementHandler(uc *usecase.MovementUseCase) *MovementHandler {
	return &MovementHandler{uc: uc}
}

func (h *MovementHandler) MoveMachine(ctx context.Context, req *pb.MoveMachineRequest) (*pb.MoveMachineResponse, error) {
	movement, err := h.uc.MoveMachine(ctx, req.GetMachineId(), req.GetToLocationId())
	if err != nil {
		if errors.Is(err, usecase.ErrSameLocation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to move machine: %v", err)
	}
	return &pb.MoveMachineResponse{Movement: movementToProto(movement)}, nil
}

func (h *MovementHandler) GetMovementHistory(ctx context.Context, req *pb.GetMovementHistoryRequest) (*pb.GetMovementHistoryResponse, error) {
	movements, err := h.uc.GetHistory(ctx, req.GetMachineId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get movement history: %v", err)
	}
	resp := &pb.GetMovementHistoryResponse{}
	for i := range movements {
		resp.Movements = append(resp.Movements, movementToProto(&movements[i]))
	}
	return resp, nil
}
