package api

import (
	"github.com/meiti-x/cli_chat/config"
	"github.com/meiti-x/cli_chat/pkg/adapters/cache"
	"github.com/meiti-x/cli_chat/pkg/app_errors"
	c "github.com/meiti-x/cli_chat/pkg/cache"
	db2 "github.com/meiti-x/cli_chat/pkg/db"
	"github.com/meiti-x/cli_chat/pkg/logger"
	nats2 "github.com/meiti-x/cli_chat/pkg/nats"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type Server struct {
	Logger *logger.AppLogger
	Nats   *nats.Conn
	Rdb    c.Provider
	Db     *gorm.DB
}

func Setup(conf *config.Config) *Server {
	lo := logger.NewAppLogger(conf)
	lo.InitLogger(conf.Logger.Path)

	nc, err := nats2.MustInitNats(conf)
	if err != nil {
		lo.DPanic(app_errors.ErrNatsInit)
	}

	rdb, err := cache.NewRedisCache(conf)
	if err != nil {
		lo.DPanic(c.ErrRedisInit)
	}

	db, err := db2.InitDB(conf)
	if err != nil {
		lo.Error(app_errors.ErrInitDB)
		panic(app_errors.ErrInitDB)
	}

	return &Server{
		Logger: lo,
		Nats:   nc,
		Rdb:    rdb,
		Db:     db,
	}
}
