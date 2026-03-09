package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	delivery "github.com/natapage/test_assignment/backend/internal/delivery/grpc"
	"github.com/natapage/test_assignment/backend/internal/repository/postgres"
	"github.com/natapage/test_assignment/backend/internal/usecase"
	"github.com/natapage/test_assignment/backend/migrations"
	pb "github.com/natapage/test_assignment/backend/pkg/gen/goldex/v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbURL := getEnv("DATABASE_URL", "postgres://goldex:goldex@localhost:5432/goldex?sslmode=disable")
	grpcPort := getEnv("GRPC_PORT", "50051")
	httpPort := getEnv("HTTP_PORT", "8080")

	// Run migrations
	if err := runMigrations(dbURL); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrations failed: %v", err)
	}
	log.Println("migrations applied successfully")

	// Connect to PostgreSQL
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	// Create repositories
	machineRepo := postgres.NewMachineRepo(pool)
	locationRepo := postgres.NewLocationRepo(pool)
	movementRepo := postgres.NewMovementRepo(pool)
	statisticsRepo := postgres.NewStatisticsRepo(pool)
	txManager := postgres.NewTxManager(pool)

	// Create use cases
	machineUC := usecase.NewMachineUseCase(machineRepo)
	locationUC := usecase.NewLocationUseCase(locationRepo)
	movementUC := usecase.NewMovementUseCase(machineRepo, locationRepo, movementRepo, txManager)
	statisticsUC := usecase.NewStatisticsUseCase(statisticsRepo)

	// Create handlers
	machineHandler := delivery.NewMachineHandler(machineUC)
	locationHandler := delivery.NewLocationHandler(locationUC)
	movementHandler := delivery.NewMovementHandler(movementUC)
	statisticsHandler := delivery.NewStatisticsHandler(statisticsUC)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterMachineServiceServer(grpcServer, machineHandler)
	pb.RegisterLocationServiceServer(grpcServer, locationHandler)
	pb.RegisterMovementServiceServer(grpcServer, movementHandler)
	pb.RegisterStatisticsServiceServer(grpcServer, statisticsHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	go func() {
		log.Printf("gRPC server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// Start gRPC-Gateway HTTP proxy
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%s", grpcPort)

	for _, reg := range []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
		pb.RegisterMachineServiceHandlerFromEndpoint,
		pb.RegisterLocationServiceHandlerFromEndpoint,
		pb.RegisterMovementServiceHandlerFromEndpoint,
		pb.RegisterStatisticsServiceHandlerFromEndpoint,
	} {
		if err := reg(ctx, mux, endpoint, opts); err != nil {
			log.Fatalf("failed to register gateway: %v", err)
		}
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)

	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: corsMiddleware(httpMux),
	}

	go func() {
		log.Printf("HTTP gateway listening on :%s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	grpcServer.GracefulStop()
	httpServer.Shutdown(ctx)
	log.Println("server stopped")
}

func runMigrations(dbURL string) error {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	return m.Up()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
