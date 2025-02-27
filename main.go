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

	// though the main thread never returns because it's running a web server, it is possible that
	// this goroutine (srv.Run) is writing to a file while we shut down the application, in which case, we'll lose data
	// make sure to stop this thread gracefully when application shuts down.
	go srv.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", srv.handlerGetRoot)
	mux.HandleFunc("POST /sign", srv.handlerPostSign)

	err = http.ListenAndServe(":8089", mux)
	// use info log to print success, and Error log in case err is not nil
	slog.ErrorContext(ctx, "listen and serve of %s server finished: %w", "GuestBook", err)
}
