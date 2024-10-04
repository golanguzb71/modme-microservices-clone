package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type LidClient struct {
	client pb.LidServiceClient
}

// NewLidClient creates a new gRPC client for LidService
func NewLidClient(addr string) (*LidClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewLidServiceClient(conn)
	return &LidClient{client: client}, nil
}

// CreateLead sends a request to create a new lead
func (lc *LidClient) CreateLead(ctx context.Context, title string) (*pb.AbsResponse, error) {
	req := &pb.CreateLeadRequest{
		Title: title,
	}

	resp, err := lc.client.CreateLead(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}

	return resp, nil
}

// GetLeadCommon retrieves common lead information by type and id
func (lc *LidClient) GetLeadCommon(ctx context.Context, leadType, id string) (*pb.GetLeadCommonResponse, error) {
	req := &pb.GetLeadCommonRequest{
		Type: leadType,
		Id:   id,
	}

	resp, err := lc.client.GetLeadCommon(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead common: %w", err)
	}

	return resp, nil
}

// UpdateLead sends a request to update a lead's title
func (lc *LidClient) UpdateLead(ctx context.Context, id, title string) (*pb.AbsResponse, error) {
	req := &pb.UpdateLeadRequest{
		Id:    id,
		Title: title,
	}

	resp, err := lc.client.UpdateLead(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update lead: %w", err)
	}

	return resp, nil
}

// DeleteLead sends a request to delete a lead by id
func (lc *LidClient) DeleteLead(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := &pb.DeleteAbsRequest{
		Id: id,
	}

	resp, err := lc.client.DeleteLead(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete lead: %w", err)
	}

	return resp, nil
}
