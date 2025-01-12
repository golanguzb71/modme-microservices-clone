package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TariffService struct {
	pb.UnimplementedTariffServiceServer
	companyRepo *repository.TariffRepository
}

func NewTariffService(repo *repository.TariffRepository) *TariffService {
	return &TariffService{
		companyRepo: repo,
	}
}

func (t *TariffService) Create(ctx context.Context, req *pb.Tariff) (*pb.Tariff, error) {
	return t.companyRepo.Create(ctx, req)
}
func (t *TariffService) Update(ctx context.Context, req *pb.Tariff) (*pb.Tariff, error) {
	return t.companyRepo.Update(ctx, req)
}
func (t *TariffService) Delete(ctx context.Context, req *pb.Tariff) (*pb.Tariff, error) {
	return t.companyRepo.Delete(ctx, req)
}
func (t *TariffService) Get(ctx context.Context, req *emptypb.Empty) (*pb.TariffList, error) {
	return t.companyRepo.Get(ctx, req)
}
