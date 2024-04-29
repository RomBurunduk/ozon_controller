package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"pvz_controller/internal/grpc/utils"
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
	tracer trace.Tracer
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
	ctx, span := s.tracer.Start(ctx, "CreatePVZ")
	defer span.End()
	defer customMetric.Add(1)
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
	ctx, span := s.tracer.Start(ctx, "GetPVZ")
	defer span.End()
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
	ctx, span := s.tracer.Start(ctx, "DeletePVZ")
	defer span.End()
	err := s.repo.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePVZResponse{Success: true}, nil
}

func (s *Server) UpdatePVZ(ctx context.Context, req *pb.UpdatePVZRequest) (*pb.UpdatePVZResponse, error) {
	ctx, span := s.tracer.Start(ctx, "UpdatePVZ")
	defer span.End()
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
	ctx, span := s.tracer.Start(ctx, "ListAllPVZ")
	defer span.End()
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
	ctx, span := s.tracer.Start(ctx, "DeleteListPVZ")
	defer span.End()
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

var (
	reg = prometheus.NewRegistry()

	customMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pvz_custom_metric",
		Help: "Total number og added PVZ",
	})
)

func init() {
	prometheus.MustRegister(customMetric)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pvzRepo := postgresql.NewPVZRepo(database)
	implementation := NewServer(pvzRepo, database.GetPool(ctx))

	shutdown, err := utils.InitProvider()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err = shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	implementation.tracer = otel.Tracer("pvz")

	fmt.Println("Starting server")

	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcMetrics)
	reg.MustRegister(customMetric)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpcMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcMetrics.StreamServerInterceptor(),
		),
	)

	pb.RegisterPVZServiceServer(grpcServer, implementation)
	grpcMetrics.InitializeMetrics(grpcServer)
	go http.ListenAndServe(":9091", promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	log.Fatal(grpcServer.Serve(lis))
}
