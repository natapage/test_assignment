package grpc

import (
	"github.com/natapage/test_assignment/backend/internal/domain"
	pb "github.com/natapage/test_assignment/backend/pkg/gen/goldex/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func locationToProto(l *domain.Location) *pb.Location {
	if l == nil {
		return nil
	}
	loc := &pb.Location{
		Id:        l.ID,
		Address:   l.Address,
		PlaceName: l.PlaceName,
		CreatedAt: timestamppb.New(l.CreatedAt),
		UpdatedAt: timestamppb.New(l.UpdatedAt),
	}
	if l.Latitude != nil {
		loc.Latitude = l.Latitude
	}
	if l.Longitude != nil {
		loc.Longitude = l.Longitude
	}
	return loc
}

func machineToProto(m *domain.Machine) *pb.Machine {
	if m == nil {
		return nil
	}
	machine := &pb.Machine{
		Id:           m.ID,
		Name:         m.Name,
		SerialNumber: m.SerialNumber,
		Enabled:      m.Enabled,
		Location:     locationToProto(m.Location),
		CreatedAt:    timestamppb.New(m.CreatedAt),
		UpdatedAt:    timestamppb.New(m.UpdatedAt),
	}
	if m.LocationID != nil {
		machine.LocationId = m.LocationID
	}
	return machine
}

func movementToProto(m *domain.Movement) *pb.Movement {
	if m == nil {
		return nil
	}
	mv := &pb.Movement{
		Id:           m.ID,
		MachineId:    m.MachineID,
		ToLocationId: m.ToLocationID,
		FromLocation: locationToProto(m.FromLocation),
		ToLocation:   locationToProto(m.ToLocation),
		MovedAt:      timestamppb.New(m.MovedAt),
		CreatedAt:    timestamppb.New(m.CreatedAt),
	}
	if m.FromLocationID != nil {
		mv.FromLocationId = m.FromLocationID
	}
	return mv
}
