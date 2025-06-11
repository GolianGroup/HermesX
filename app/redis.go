package app

import "hermesx/internal/producers"

func (a *application) InitRedis() producers.RedisClient {
	return producers.NewRedis(&a.config.Redis)
}
