package main

import (
	"context"
	"log/slog"
	"net/http"
)

func main() {
	ctx := context.Background()

	db, err := NewFileDB("/tmp")
	if err != nil {
		panic(err)
	}

	srv := NewServer(db)

	go srv.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", srv.handlerGetRoot)
	mux.HandleFunc("POST /sign", srv.handlerPostSign)

	err = http.ListenAndServe(":8089", mux)
	slog.ErrorContext(ctx, "listen and serve of %s server finished: %w", "GuestBook", err)
}
