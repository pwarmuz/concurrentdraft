package main

import (
	"math/rand"
	"sync"
	"time"
)

type Randomizer struct {
	random *rand.Rand
	mu     sync.Mutex
}

func (x *Randomizer) Initialize() {
	x.random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Age randomly created
func (x *Randomizer) Age() int {
	MaxAge := 115
	MinAge := 1
	t := x.random.Intn(MaxAge-MinAge) + MinAge
	return t
}

// Height randomly created
func (x *Randomizer) Height() int {
	MaxHeight := 84
	MinHeight := 6
	t := x.random.Intn(MaxHeight-MinHeight) + MinHeight
	return t
}

// Weight randomly created
func (x *Randomizer) Weight() int {
	MaxWeight := 750
	MinWeight := 4
	t := x.random.Intn(MaxWeight-MinWeight) + MinWeight
	return t
}

// Assign random data to channel recieving individual data
func (x *Randomizer) Assign(ch chan<- individual) {
	x.mu.Lock()
	a := x.Age()
	h := x.Height()
	w := x.Weight()
	x.mu.Unlock()
	ch <- individual{a, h, w}
}
