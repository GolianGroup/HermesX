package app

import (
	"context"
	"hermesx/handler/routers"
	"hermesx/internal/config"

	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Application interface {
	Setup()
}

type application struct {
	ctx    context.Context
	config *config.Config
}

func NewApplication(ctx context.Context, config *config.Config) Application {
	return &application{ctx: ctx, config: config}
}

// bootstrap

func (a *application) Setup() {
	app := fx.New(
		fx.Provide(
			a.InitRouter,
			a.InitFramework,
			a.InitController,
			a.InitServices,
			a.InitRepositories,
			a.InitRedis,
			a.InitScyllaDB,
			a.InitDatabase,
			a.InitArangoDB,
			a.InitLogger,
			a.InitGRPCServer,
			a.InitNats,
		),

		fx.Invoke(func(lc fx.Lifecycle, grpcServer *grpc.Server, logger *zap.Logger) {
			logger.Info("Initializing gRPC server")
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					listener, err := net.Listen("tcp", a.config.GRPC.Host+":"+a.config.GRPC.Port)
					if err != nil {
						logger.Info("Failed to listen on gRPC port", zap.Error(err))
						return err
					}
					logger.Info("starting gRPC Server",
						zap.String("host", a.config.GRPC.Host),
						zap.String("port", a.config.GRPC.Port),
					)
					go func() {
						if err := grpcServer.Serve(listener); err != nil {
							logger.Info("Failed to serve gRPC", zap.Error(err))
						}
					}()
					log.Println("gRPC server started on", a.config.GRPC.Host+":"+a.config.GRPC.Port)
					return nil
				},
				OnStop: func(_ context.Context) error {
					grpcServer.Stop()
					logger.Info("gRPC server stopped")
					return nil
				},
			})
		}),

		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, logger *zap.Logger) {
			// Start Fiber server in a separate goroutine
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("Starting Fiber server")
					go func() {
						if err := app.Listen(a.config.Server.Host + ":" + a.config.Server.Port); err != nil {
							logger.Error("Failed to start Fiber server", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),

		fx.Invoke(func(app *fiber.App, router routers.Router) {
			router.AddRoutes(app.Group(""))
		}),
	)
	app.Run()
}
