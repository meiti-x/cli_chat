package test

import (
	"fmt"
	"github.com/meiti-x/snapp_chal/api/http"
	"log"
	"os"
	"testing"
	"time"

	"github.com/meiti-x/snapp_chal/config"
	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	conf, err := config.LoadConfig("../test_config.yml")
	if err != nil {
		log.Panic(fmt.Errorf("load config error: %w", err))
	}

	s, err := http.NewServer(conf)
	if err != nil {
		log.Panic(fmt.Errorf("could not start server: %w", err))
	}

	testDB = s.DB

	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	code := m.Run()

	os.Exit(code)
}
