package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/natapage/test_assignment/backend/internal/usecase"
	pb "github.com/natapage/test_assignment/backend/pkg/gen/goldex/v1"
)

type StatisticsHandler struct {
	pb.UnimplementedStatisticsServiceServer
	uc *usecase.StatisticsUseCase
}

func NewStatisticsHandler(uc *usecase.StatisticsUseCase) *StatisticsHandler {
	return &StatisticsHandler{uc: uc}
}

func (h *StatisticsHandler) GetLocationDurations(ctx context.Context, req *pb.GetLocationDurationsRequest) (*pb.GetLocationDurationsResponse, error) {
	durations, err := h.uc.GetLocationDurations(ctx, req.GetMachineId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get location durations: %v", err)
	}
	resp := &pb.GetLocationDurationsResponse{}
	for _, d := range durations {
		resp.Durations = append(resp.Durations, &pb.LocationDuration{
			LocationId: d.LocationID,
			Location:   locationToProto(d.Location),
			Days:       int32(d.Days),
		})
	}
	return resp, nil
}

func (h *StatisticsHandler) GetMovementsCount(ctx context.Context, req *pb.GetMovementsCountRequest) (*pb.GetMovementsCountResponse, error) {
	from := req.GetFrom().AsTime()
	to := req.GetTo().AsTime()
	counts, err := h.uc.GetMovementsCount(ctx, from, to)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get movements count: %v", err)
	}
	resp := &pb.GetMovementsCountResponse{}
	for _, c := range counts {
		mc := &pb.MovementsCount{
			MachineId: c.MachineID,
			Count:     int32(c.Count),
		}
		if c.Machine != nil {
			mc.MachineName = c.Machine.Name
		}
		resp.Counts = append(resp.Counts, mc)
	}
	return resp, nil
}

func (h *StatisticsHandler) GetMachineTimeline(ctx context.Context, req *pb.GetMachineTimelineRequest) (*pb.GetMachineTimelineResponse, error) {
	entries, err := h.uc.GetMachineTimeline(ctx, req.GetMachineId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get machine timeline: %v", err)
	}
	resp := &pb.GetMachineTimelineResponse{}
	for _, e := range entries {
		resp.Entries = append(resp.Entries, &pb.TimelineEntry{
			MovedAt:      timestamppb.New(e.MovedAt),
			FromLocation: locationToProto(e.FromLocation),
			ToLocation:   locationToProto(e.ToLocation),
		})
	}
	return resp, nil
}
