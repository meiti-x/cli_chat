package api

import (
	"github.com/meiti-x/snapp_task/config"
	"github.com/meiti-x/snapp_task/pkg/adapters/cache"
	"github.com/meiti-x/snapp_task/pkg/app_errors"
	c "github.com/meiti-x/snapp_task/pkg/cache"
	"github.com/meiti-x/snapp_task/pkg/logger"
	nats2 "github.com/meiti-x/snapp_task/pkg/nats"
	"github.com/nats-io/nats.go"
)

type Server struct {
	logger *logger.AppLogger
	nats   *nats.Conn
	rdb    c.Provider
}

func Setup(conf *config.Config) *Server {
	lo := logger.NewAppLogger(conf)
	lo.InitLogger(conf.Logger.Path)

	nc, err := nats2.MustInitNats(conf)
	if err != nil {
		lo.DPanic(app_errors.ErrNatsInit)
	}

	rdb := cache.NewRedisCache(conf)

	return &Server{
		rdb:    rdb,
		nats:   nc,
		logger: lo,
	}
}
