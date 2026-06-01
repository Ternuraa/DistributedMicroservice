package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ternuraa/DistributedMicroservice/services/listing_service/internal/config"
	proto "github.com/Ternuraa/DistributedMicroservice/services/listing_service/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Prometheus метрики для gRPC
var (
	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "listing_grpc_requests_total",
			Help: "Общее количество входящих gRPC запросов",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(grpcRequestsTotal)
}

// Interceptor для извлечения Correlation ID и сбора метрик
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Считаем метрику запроса
	grpcRequestsTotal.WithLabelValues(info.FullMethod).Inc()

	md, ok := metadata.FromIncomingContext(ctx)
	correlationID := "unknown"
	if ok && len(md["x-correlation-id"]) > 0 {
		correlationID = md["x-correlation-id"][0]
	}

	slog.Info("Получен gRPC запрос", "method", info.FullMethod, "correlation_id", correlationID)

	return handler(ctx, req)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Init("../../configs/listing_config.yaml")
	if err != nil {
		slog.Error("Ошибка конфигурации", "err", err)
		os.Exit(1)
	}

	// --- КЛАУД-НАЙТИВ НАСТРОЙКА (Для Docker / K8s) ---
	// Обязательное требование 4-го блока: переменные окружения для конфигурации.
	// Если переменная окружения задана (в .env, Docker Compose или Kubernetes), перезаписываем порт.
	if envPort := os.Getenv("LISTING_PORT"); envPort != "" {
		cfg.Service.Port = envPort
	}

	// Важно для контейнеризации: внутри Docker/K8s нельзя жестко привязываться к 127.0.0.1,
	// иначе другие контейнеры не увидят этот сервис.
	// Если запуск идет через окружение (Docker), слушаем на всех интерфейсах (например, ":50051").
	listenAddr := "127.0.0.1" + cfg.Service.Port
	if os.Getenv("LISTING_PORT") != "" {
		listenAddr = cfg.Service.Port // Становится видом ":50051", что открывает порт наружу внутри Docker-сети
	}
	// -------------------------------------------------

	// 2. Health Check + Prometheus метрики
	go func() {
		healthMux := http.NewServeMux()
		healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		healthMux.Handle("/metrics", promhttp.Handler()) // Добавили метрики

		slog.Info("Health check и Prometheus запущены", "port", ":8081")
		if err := http.ListenAndServe(":8081", healthMux); err != nil {
			slog.Error("Ошибка запуска HTTP сервера", "err", err)
		}
	}()

	// 3. gRPC сервер (теперь использует динамический listenAddr для поддержки Docker)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		slog.Error("Не удалось слушать порт", "err", err)
		os.Exit(1)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	proto.RegisterListingServiceServer(s, &ListingServer{})

	// 4. Graceful Shutdown
	go func() {
		slog.Info("Listing Service запущен", "port", cfg.Service.Port)
		if err := s.Serve(lis); err != nil {
			slog.Error("Ошибка gRPC сервера", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Завершение работы Listing Service...")
	s.GracefulStop()
	slog.Info("Сервер остановлен")
}
