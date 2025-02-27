package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGuest(t *testing.T) {
	guests := NewGuests()
	expected := Guests{}
	assert.Equal(t, reflect.TypeOf(expected).String(), reflect.TypeOf(*guests).String(), "Guests struct returned from NewGuests must be of type Guests")
}

func TestIsSpecial(t *testing.T) {
	guests := NewGuests()
	guests.Add("Normal Guest", false)
	guests.Add("Special Guest", true)
	assert.True(t, guests.IsSpecial("Special Guest"))
	assert.False(t, guests.IsSpecial("Normal Guest"))
	assert.False(t, guests.IsSpecial("Not A. Guest"))
}
