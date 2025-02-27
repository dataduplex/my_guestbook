package main

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type DB interface {
	SaveGuests(context.Context, *Guests) error
}

type Server struct {
	db     DB
	visits int
	guests *Guests
}

func NewServer(db DB) *Server {
	return &Server{
		db:     db,
		guests: NewGuests(),
	}
}

func (s *Server) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			slog.InfoContext(ctx, "persisting guests to db")
			err := s.db.SaveGuests(ctx, s.guests)
			slog.InfoContext(ctx, "PERSISTED GUESTS TO DB", "error", err.Error())
		}
	}
}

func (s *Server) handlerGetRoot(rw http.ResponseWriter, req *http.Request) {
	s.visits++

	var buf bytes.Buffer

	buf.WriteString("Visits: ")
	buf.WriteString(strconv.Itoa(s.visits))
	buf.WriteString("\n\n")

	buf.WriteString("Guests:\n")

	for name := range s.guests.Guests {
		if s.guests.IsSpecial(name) {
			buf.WriteString("* ")
		} else {
			buf.WriteString("- ")
		}
		buf.WriteString(name)
		buf.WriteString("\n")
	}

	rw.Write(buf.Bytes())
}

func (s *Server) handlerPostSign(rw http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	special, _ := strconv.ParseBool(req.FormValue("special"))

	s.guests.Add(name, special)
}
