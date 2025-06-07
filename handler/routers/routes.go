package routers

import (
	"golang_template/handler/controllers"
	"golang_template/internal/producers"

	"golang_template/handler/middlewares"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type Router interface {
	AddRoutes(router fiber.Router)
}

type router struct {
	systemRouter SystemRouter
	redisClient  producers.RedisClient
	tracer       trace.Tracer
}

func NewRouter(controllers controllers.Controllers, redisClient producers.RedisClient, tracer trace.Tracer) Router {

	systemRouter := NewSystemRouter(controllers.SystemController())
	return &router{
		systemRouter: systemRouter,
		redisClient:  redisClient,
		tracer:       tracer,
	}
}

func (r router) AddRoutes(router fiber.Router) {

	// router
	// init user router, etc ...
	// rate limiter
	// CORS
	router.Use(middlewares.TracingMiddleware(r.tracer))

	r.systemRouter.AddRoutes(router)
}
