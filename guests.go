package main

import (
	"maps"
	"sync"
)

type Guests struct {
	mu sync.RWMutex

	guests map[string]bool // guest name -> special status
}

func NewGuests() *Guests {
	return &Guests{
		guests: make(map[string]bool),
	}
}

func (g *Guests) Clone() *Guests {
	g.mu.Lock()
	defer g.mu.Unlock()

	guests := maps.Clone(g.guests)
	return &Guests{guests: guests}
}

// potential performance bottleneck
func (g *Guests) Add(name string, special bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.guests[name] = special
}

// potential performance bottleneck especially if we have large number of guests
// write operation on guests map will be blocking until this method finishes
func (g *Guests) Guests(cb func(name string) bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for name := range g.guests {
		if !cb(name) {
			break
		}
	}
}

func (g *Guests) IsSpecial(name string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.guests[name]
}
