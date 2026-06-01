package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	proto "github.com/Ternuraa/DistributedMicroservice/listing_service/proto"
	"github.com/Ternuraa/DistributedMicroservice/services/booking_service/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Определение Prometheus метрик
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_http_requests_total",
			Help: "Общее количество HTTP запросов к booking_service",
		},
		[]string{"path", "status"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "booking_http_duration_seconds",
			Help:    "Длительность обработки HTTP запросов",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	// Регистрация метрик в Prometheus
	prometheus.MustRegister(httpRequestsTotal, httpDuration)
}

type BookRequest struct {
	ListingID string `json:"listing_id"`
	UserID    string `json:"user_id"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Init("../../configs/booking_config.yaml")
	if err != nil {
		slog.Error("Ошибка загрузки конфига", "err", err)
		os.Exit(1)
	}

	// --- КЛАУД-НАЙТИВ НАСТРОЙКА (Для Docker / K8s) ---
	if envPort := os.Getenv("BOOKING_PORT"); envPort != "" {
		cfg.Service.Port = envPort
	}
	if envListingAddr := os.Getenv("LISTING_SERVICE_ADDR"); envListingAddr != "" {
		cfg.ListingService.Address = envListingAddr
	}
	// -------------------------------------------------

	conn, err := grpc.Dial(cfg.ListingService.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("gRPC connection error", "err", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := proto.NewListingServiceClient(conn)

	// 3. Health Check + Метрики (на общем порту 8080)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	healthMux.Handle("/metrics", promhttp.Handler())

	slog.Info("Health check и Prometheus запущены", "port", ":8080")
	go http.ListenAndServe(":8080", healthMux)

	// 4. Основной сервер
	mux := http.NewServeMux()

	// --- FEATURE FLAG  ---
	mux.HandleFunc("/premium", func(w http.ResponseWriter, r *http.Request) {
		// Читаем флаг из переменной окружения
		isPremiumEnabled := os.Getenv("ENABLE_PREMIUM_FEATURES")

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if isPremiumEnabled == "true" {
			// Логика включена
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "success",
				"message": "✨ Премиум-функционал успешно работает! ✨",
				"feature": "early_access_booking",
			})
		} else {
			// Логика выключена
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "error",
				"message": "Премиум-функционал отключен (Feature Flag = false)",
			})
		}
	})
	// ---------------------------------------

	mux.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		statusCode := http.StatusOK

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Correlation-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Сбор метрики длительности в конце запроса
		defer func() {
			httpDuration.WithLabelValues("/book").Observe(time.Since(startTime).Seconds())
			httpRequestsTotal.WithLabelValues("/book", strconv.Itoa(statusCode)).Inc()
		}()

		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = "req-" + time.Now().Format("20060102150405")
		}

		var reqBook BookRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBook); err != nil {
			slog.Warn("Неверный JSON", "correlation_id", correlationID)
			statusCode = http.StatusBadRequest
			http.Error(w, "Invalid JSON", statusCode)
			return
		}

		slog.Info("Обработка бронирования", "listing_id", reqBook.ListingID, "correlation_id", correlationID)

		md := metadata.Pairs("x-correlation-id", correlationID)

		// --- Внедрение Retry Policy с Exponential Backoff & Timeouts ---
		var resp *proto.ListingResponse
		var grpcErr error

		maxRetries := 3
		backoff := 100 * time.Millisecond // Стартовая пауза

		for i := 0; i < maxRetries; i++ {
			retryCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			retryCtx = metadata.NewOutgoingContext(retryCtx, md)

			resp, grpcErr = client.GetListingInfo(retryCtx, &proto.ListingRequest{Id: reqBook.ListingID})
			cancel()

			if grpcErr == nil {
				break
			}

			if i < maxRetries-1 {
				slog.Warn("gRPC вызов завершился ошибкой, выполняем повторную попытку",
					"attempt", i+1, "backoff", backoff, "err", grpcErr, "correlation_id", correlationID)
				time.Sleep(backoff)
				backoff *= 3
			}
		}

		if grpcErr != nil {
			slog.Error("Все попытки gRPC вызова исчерпаны", "err", grpcErr, "correlation_id", correlationID)
			statusCode = http.StatusNotFound
			http.Error(w, "Listing not found after retries", statusCode)
			return
		}
		// -----------------------------------------------------------------

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "data": resp})
	})

	// Timeout Policy для всего HTTP сервера
	server := &http.Server{
		Addr:         cfg.Service.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 5. Graceful Shutdown
	go func() {
		slog.Info("Booking Service запущен", "port", cfg.Service.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка сервера", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Завершение работы Booking Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	slog.Info("Сервер остановлен")
}
