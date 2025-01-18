package tests

import (
	"messenger/internal/bootstrap"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	client = &http.Client{
		Timeout: time.Second * 5,
	}

	app = bootstrap.NewApp()
	resetTestDb()

	defer app.Cleanup()
	m.Run()
}
