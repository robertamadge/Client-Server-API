package main

import (
	"github.com/google/uuid"
)

type Bid struct {
	ID   string
	Name string
	Rate float64
}

func NewExchangeRate(name string, rate float64) *Bid {
	return &Bid{
		ID:   uuid.New().String(),
		Name: name,
		Rate: rate,
	}
}
