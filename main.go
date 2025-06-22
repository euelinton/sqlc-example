package main

import (
	"context"
	"net-http/internal/repository"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://postgres:12345@localhost:5432/testes")
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(ctx)

	repo := repository.New(conn)

	server := NewAPIServer(":8000", repo, ctx)
	server.Run()
}
