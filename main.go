package main

import (
	"context"
	"log"
	db "net-http/internal/repository"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://postgres:12345@localhost:5432/testes")
	if err != nil {
		log.Fatal("Erro ao se conectar com o banco de dados")
	}
	defer conn.Close(ctx)

	store := db.New(conn)

	server := NewAPIServer(":8000", store)
	server.Run()
}
