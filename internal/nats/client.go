package nats

import (
	"github.com/meiti-x/snapp_chal/config"
	"github.com/meiti-x/snapp_chal/pkg/app_errors"
	"github.com/nats-io/nats.go"
	"log"
)

func MustInitNats(conf *config.Config) (*nats.Conn, error) {
	nc, err := nats.Connect(conf.Nats.ConnString)
	if err != nil {
		log.Fatalln(app_errors.ErrNATSConnectionFailed, err)
	}
	return nc, err
}
