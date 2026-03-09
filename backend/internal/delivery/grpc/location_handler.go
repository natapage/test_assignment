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

type LocationHandler struct {
	pb.UnimplementedLocationServiceServer
	uc *usecase.LocationUseCase
}

func NewLocationHandler(uc *usecase.LocationUseCase) *LocationHandler {
	return &LocationHandler{uc: uc}
}

func (h *LocationHandler) ListLocations(ctx context.Context, _ *pb.ListLocationsRequest) (*pb.ListLocationsResponse, error) {
	locations, err := h.uc.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list locations: %v", err)
	}
	resp := &pb.ListLocationsResponse{}
	for i := range locations {
		resp.Locations = append(resp.Locations, locationToProto(&locations[i]))
	}
	return resp, nil
}

func (h *LocationHandler) GetLocation(ctx context.Context, req *pb.GetLocationRequest) (*pb.GetLocationResponse, error) {
	l, err := h.uc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "location not found: %v", err)
	}
	return &pb.GetLocationResponse{Location: locationToProto(l)}, nil
}

func (h *LocationHandler) CreateLocation(ctx context.Context, req *pb.CreateLocationRequest) (*pb.CreateLocationResponse, error) {
	l := &domain.Location{
		Address:   req.GetAddress(),
		PlaceName: req.GetPlaceName(),
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	created, err := h.uc.Create(ctx, l)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create location: %v", err)
	}
	return &pb.CreateLocationResponse{Location: locationToProto(created)}, nil
}

func (h *LocationHandler) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.UpdateLocationResponse, error) {
	l := &domain.Location{
		ID:        req.GetId(),
		Address:   req.GetAddress(),
		PlaceName: req.GetPlaceName(),
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	updated, err := h.uc.Update(ctx, l)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update location: %v", err)
	}
	return &pb.UpdateLocationResponse{Location: locationToProto(updated)}, nil
}

func (h *LocationHandler) DeleteLocation(ctx context.Context, req *pb.DeleteLocationRequest) (*emptypb.Empty, error) {
	if err := h.uc.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete location: %v", err)
	}
	return &emptypb.Empty{}, nil
}
