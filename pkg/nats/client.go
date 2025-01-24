package nats

import (
	"github.com/meiti-x/snapp_task/config"
	"github.com/meiti-x/snapp_task/pkg/app_errors"
	"github.com/nats-io/nats.go"
	"log"
)

// MustInitNats creates a new nats connection and panics if it fails
func MustInitNats(conf *config.Config) (*nats.Conn, error) {
	nc, err := nats.Connect(conf.Nats.ConnString)
	if err != nil {
		log.Fatalln(app_errors.ErrNATSConnectionFailed, err)
	}
	return nc, err
}
