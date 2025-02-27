package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerPersist(t *testing.T) {
	db := NewMockDB()
	srv := NewServer(db)
	go srv.Run(context.Background())
	time.Sleep(2 * time.Minute)
	assert.NotNil(t, db.saved)
}

type MockDB struct {
	saved *Guests
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (m *MockDB) SaveGuests(ctx context.Context, g *Guests) error {
	m.saved = g
	return nil
}
