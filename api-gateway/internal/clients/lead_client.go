package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LidClient struct {
	leadClient     pb.LeadServiceClient
	expectClient   pb.ExpectServiceClient
	setClient      pb.SetServiceClient
	leadDataClient pb.LeadDataServiceClient
}

// NewLidClient creates a new gRPC client for LidService
func NewLidClient(addr string) (*LidClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	leadClient := pb.NewLeadServiceClient(conn)
	expectClient := pb.NewExpectServiceClient(conn)
	setClient := pb.NewSetServiceClient(conn)
	leadDataClient := pb.NewLeadDataServiceClient(conn)

	return &LidClient{leadClient: leadClient, expectClient: expectClient, setClient: setClient, leadDataClient: leadDataClient}, nil
}

// LeadService methods
func (lc *LidClient) CreateLead(ctx context.Context, title string) (*pb.AbsResponse, error) {
	req := &pb.CreateLeadRequest{
		Title: title,
	}

	resp, err := lc.leadClient.CreateLead(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) GetLeadCommon(ctx context.Context, req *pb.GetLeadCommonRequest) (*pb.GetLeadCommonResponse, error) {
	resp, err := lc.leadClient.GetLeadCommon(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead common: %w", err)
	}
	return resp, nil
}

func (lc *LidClient) UpdateLead(ctx context.Context, id, title string) (*pb.AbsResponse, error) {
	req := &pb.UpdateLeadRequest{
		Id:    id,
		Title: title,
	}

	resp, err := lc.leadClient.UpdateLead(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update lead: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) DeleteLead(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := &pb.DeleteAbsRequest{
		Id: id,
	}

	resp, err := lc.leadClient.DeleteLead(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete lead: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) GetAllLead(ctx context.Context) (*pb.GetLeadListResponse, error) {
	emptyRequest := &emptypb.Empty{}
	return lc.leadClient.GetListSection(ctx, emptyRequest)
}

// ExpectService methods
func (lc *LidClient) CreateExpect(ctx context.Context, title string) (*pb.AbsResponse, error) {
	req := &pb.CreateExpectRequest{
		Title: title,
	}

	resp, err := lc.expectClient.CreateExpect(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create expectation: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) UpdateExpect(ctx context.Context, id, title string) (*pb.AbsResponse, error) {
	req := &pb.UpdateExpectRequest{
		Id:    id,
		Title: title,
	}

	resp, err := lc.expectClient.UpdateExpect(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update expectation: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) DeleteExpect(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := &pb.DeleteAbsRequest{
		Id: id,
	}

	resp, err := lc.expectClient.DeleteExpect(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete expectation: %w", err)
	}

	return resp, nil
}

// SetService methods
func (lc *LidClient) CreateSet(ctx context.Context, req *pb.CreateSetRequest) (*pb.AbsResponse, error) {
	resp, err := lc.setClient.CreateSet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create set: %w", err)
	}
	return resp, nil
}

func (lc *LidClient) UpdateSet(ctx context.Context, req *pb.UpdateSetRequest) (*pb.AbsResponse, error) {
	resp, err := lc.setClient.UpdateSet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update set: %w", err)
	}
	return resp, nil
}

func (lc *LidClient) DeleteSet(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := &pb.DeleteAbsRequest{
		Id: id,
	}

	resp, err := lc.setClient.DeleteSet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete set: %w", err)
	}

	return resp, nil
}

// LeadDataService methods
func (lc *LidClient) CreateLeadData(ctx context.Context, req *pb.CreateLeadDataRequest) (*pb.AbsResponse, error) {
	resp, err := lc.leadDataClient.CreateLeadData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead data: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) UpdateLeadData(ctx context.Context, req *pb.UpdateLeadDataRequest) (*pb.AbsResponse, error) {
	resp, err := lc.leadDataClient.UpdateLeadData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update lead data: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) DeleteLeadData(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := &pb.DeleteAbsRequest{
		Id: id,
	}

	resp, err := lc.leadDataClient.DeleteLeadData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete lead data: %w", err)
	}

	return resp, nil
}

func (lc *LidClient) ChangeLeadPlace(ctx context.Context, p *pb.ChangeLeadPlaceRequest) (*pb.AbsResponse, error) {
	resp, err := lc.leadDataClient.ChangeLeadPlace(ctx, p)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (lc *LidClient) ChangeSetToGroup(ctx context.Context, p *pb.ChangeToSetRequest) (*pb.AbsResponse, error) {
	resp, err := lc.setClient.ChangeToSet(ctx, p)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (lc *LidClient) GetLeadReports(from string, till string, ctx context.Context) (*pb.GetLeadReportsResponse, error) {
	return lc.leadClient.GetLeadReports(ctx, &pb.GetLeadReportsRequest{
		StartYear: from,
		EndYear:   till,
	})
}

func (lc *LidClient) GetActiveLeadCount(ctx context.Context) int {
	count, err := lc.leadClient.GetActiveLeadCount(ctx, &emptypb.Empty{})
	if err != nil {
		return 0
	}
	return int(count.ActiveLeadCount)
}

func (lc *LidClient) GetByIdSet(setId string, ctx context.Context) (*pb.SetDataResponse, error) {
	return lc.setClient.GetById(ctx, &pb.DeleteAbsRequest{Id: setId})
}
