package app

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Tel3scop/auth/internal/closer"
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/interceptor"
	"github.com/Tel3scop/auth/internal/metrics"
	accessAPI "github.com/Tel3scop/auth/pkg/access_v1"
	authAPI "github.com/Tel3scop/auth/pkg/auth_v1"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/Tel3scop/helpers/logger"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	// Register statik for swagger UI
	_ "github.com/Tel3scop/auth/statik"
)

// App структура приложения с сервис-провайдером и GRPC-сервером
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

// NewApp вернуть новый экземпляр приложения с зависимостями
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запуск приложения
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		err := a.runHTTPServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runSwaggerServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initMetrics,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	_, err := config.New()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initMetrics(_ context.Context) error {
	err := metrics.Init(
		a.serviceProvider.Config().Metrics.Namespace,
		a.serviceProvider.Config().Metrics.AppName,
		a.serviceProvider.Config().Metrics.Subsystem,
		a.serviceProvider.Config().Metrics.BucketsStart,
		a.serviceProvider.Config().Metrics.BucketsFactor,
		a.serviceProvider.Config().Metrics.BucketsCount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.InitByParams(
		a.serviceProvider.Config().Log.FileName,
		a.serviceProvider.Config().Log.Level,
		a.serviceProvider.Config().Log.MaxSize,
		a.serviceProvider.Config().Log.MaxBackups,
		a.serviceProvider.Config().Log.MaxAge,
		a.serviceProvider.Config().Log.Compress,
	)
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			interceptor.ValidateInterceptor,
			interceptor.LogInterceptor,
			interceptor.MetricsInterceptor,
		)),
	)

	reflection.Register(a.grpcServer)
	userAPI.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	authAPI.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	accessAPI.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.Config().GRPC.Address)

	list, err := net.Listen("tcp", a.serviceProvider.Config().GRPC.Address)
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := userAPI.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.Config().GRPC.Address, opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.Config().HTTP.Address,
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 3 * time.Second,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api_user.swagger.json", serveSwaggerFile("/api_user.swagger.json"))
	mux.HandleFunc("/api_auth.swagger.json", serveSwaggerFile("/api_auth.swagger.json"))
	mux.HandleFunc("/api_access.swagger.json", serveSwaggerFile("/api_access.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.Config().Swagger.Address,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(a.serviceProvider.Config().Swagger.Timeout) * time.Second,
	}

	return nil
}

func (a *App) runPrometheus() error {
	logger.Info("Prometheus server is running on %s", zap.String("address", a.serviceProvider.Config().Metrics.Address), zap.String("port", a.serviceProvider.Config().Metrics.Port))

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:              net.JoinHostPort(a.serviceProvider.Config().Metrics.Address, a.serviceProvider.Config().Metrics.Port),
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
	}

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.httpServer.Addr)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.swaggerServer.Addr)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file http.File) {
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
