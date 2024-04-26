package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"

	"pvz_controller/internal/model"
	"pvz_controller/internal/pkg/db"
	"pvz_controller/internal/pkg/pb"
	"pvz_controller/internal/pkg/repository"
	"pvz_controller/internal/pkg/repository/in_memory"
	"pvz_controller/internal/pkg/repository/postgresql"
	"pvz_controller/internal/pkg/repository/redis"
	"pvz_controller/internal/pkg/repository/transaction_manager"
	"pvz_controller/internal/service"
)

type Server struct {
	pb.UnimplementedPVZServiceServer

	repo repository.PVZRepo

	cache     service.InMemory
	redis     service.Redis
	txManager service.TransactionManager
}

func NewServer(repo repository.PVZRepo, db *pgxpool.Pool) *Server {
	c := in_memory.NewInMemoryCache()
	r := redis.NewRedis()
	tx := transaction_manager.NewTransactionManager(db)
	return &Server{
		repo:      repo,
		cache:     c,
		redis:     r,
		txManager: tx}
}

func (s *Server) CreatePVZ(ctx context.Context, req *pb.CreatePVZRequest) (*pb.CreatePVZResponse, error) {
	pvz := &model.Pickups{
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
	}
	id, err := s.repo.Add(ctx, pvz)
	if err != nil {
		return nil, err
	}
	response := &pb.CreatePVZResponse{
		Success: true,
		Id:      int32(id),
	}
	return response, nil
}

func (s *Server) GetPVZ(ctx context.Context, req *pb.GetPVZRequest) (*pb.GetPVZResponse, error) {
	item, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetPVZResponse{Pvz: &pb.PVZ{
		Id:      item.Id,
		Name:    item.Name,
		Address: item.Address,
	}}, nil
}

func (s *Server) DeletePVZ(ctx context.Context, req *pb.DeletePVZRequest) (*pb.DeletePVZResponse, error) {
	err := s.repo.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePVZResponse{Success: true}, nil
}

func (s *Server) UpdatePVZ(ctx context.Context, req *pb.UpdatePVZRequest) (*pb.UpdatePVZResponse, error) {
	err := s.repo.Update(ctx, &model.Pickups{
		Name:    req.Pvz.Name,
		Address: req.Pvz.Address,
		Contact: req.Pvz.Contact,
	}, req.Pvz.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePVZResponse{Success: true}, nil
}

func (s *Server) ListAllPVZ(ctx context.Context, req *pb.ListAllPVZRequest) (*pb.ListAllPVZResponse, error) {
	all, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]*pb.PVZ, 0, len(all))
	for _, item := range all {
		response = append(response, &pb.PVZ{
			Id:      item.Id,
			Name:    item.Name,
			Address: item.Address,
			Contact: item.Contact,
		})
	}
	return &pb.ListAllPVZResponse{Pvzs: response}, nil
}

func (s *Server) DeleteListPVZ(ctx context.Context, req *pb.DeleteListPVZRequest) (*pb.DeleteListPVZResponse, error) {
	err := s.txManager.RunSerializable(ctx, func(ctxTX context.Context) error {
		for _, id := range req.Ids {
			err := s.repo.DeleteByID(ctx, id)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteListPVZResponse{Success: true}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pvzRepo := postgresql.NewPVZRepo(database)
	implementation := NewServer(pvzRepo, database.GetPool(ctx))
	fmt.Println("Starting server")
	pb.RegisterPVZServiceServer(s, implementation)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}