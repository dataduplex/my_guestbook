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
			// make sure to save to call SaveGuests before returning here, or in a defer statement
			// otherwise, we might end up losing data
			return

		case <-ticker.C:
			slog.InfoContext(ctx, "persisting guests to db")
			err := s.db.SaveGuests(ctx, s.guests)
			// Check for error here
			// Print success only when error is nil, else print the error message
			slog.InfoContext(ctx, "PERSISTED GUESTS TO DB", "error", err.Error())
		}
	}
}

func (s *Server) handlerGetRoot(rw http.ResponseWriter, req *http.Request) {
	// 1. unsafe operation as multiple threads can simultaneously increment s.visits
	// 2. this shall be moved to a function of it's own (example func (*Server) IncrementVisits() int) which is protected by a write lock
	// 3. save the visits to a local variable after increment since it's used to return the response
	s.visits++

	var buf bytes.Buffer

	buf.WriteString("Visits: ")
	// this s.visits may not be accurate as the increment step earlier and the read step below are not an atomic operation
	// after increment, save it a local var and use it here.
	buf.WriteString(strconv.Itoa(s.visits))
	buf.WriteString("\n\n")

	buf.WriteString("Guests:\n")

	// possible deadlock here as Guests and IsSpecial are read locking in succession without releasing
	// deadlock will happen when another thread tries to acquire a write lock between (1) and (2)
	for name := range s.guests.Guests { // (1) locking here
		if s.guests.IsSpecial(name) { // (2) locking again here without releasing the earlier one
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
	// validate "name" to make sure it is not empty
	name := req.FormValue("name")

	// Handle error here, return without updating the map if an invalid input was passed.
	special, _ := strconv.ParseBool(req.FormValue("special"))

	s.guests.Add(name, special)
}
