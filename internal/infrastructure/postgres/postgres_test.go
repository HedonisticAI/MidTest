package postgres_test

import (
	"context"
	"log/slog"
	"midtest/internal/infrastructure/postgres"
	"testing"
)

var DSN = "postgresql://postgres:password@127.0.0.1:5432/test?sslmode=disable"

func TestSimple(t *testing.T) {
	ctx := context.Background()

	_, err := postgres.Open(ctx, DSN)
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
	}
}
