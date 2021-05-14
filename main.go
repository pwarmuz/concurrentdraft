package main

import (
	"fmt"
	"sync"
)

type individual struct {
	Age    int
	Height int
	Weight int
}

type DataStore struct {
	Individuals []individual
}

func (i *DataStore) Capacity(max int) {
	i.Individuals = make([]individual, 0, max)
}

func (i *DataStore) Assign(index, age, height, weight int) {
	i.Individuals[index] = individual{age, height, weight}
}

// BuildDataSequential without concurrency
func BuildDataSequential(data *DataStore) {
	var gRand Randomizer
	gRand.Initialize()

	// Generate random individual data
	for i := 0; i < cap(data.Individuals); i++ {
		data.Individuals = append(data.Individuals, individual{gRand.Age(), gRand.Height(), gRand.Weight()})
	}
}

// BuildDataConcurrently simulates many users trying to create data
// This also builds the data that is used for the streaming, drafting and fan in functions
func BuildDataConcurrently(data *DataStore) {
	// make unbuffered channel of type individual
	ch := make(chan individual)
	// waitgroup used to help build data to capacity
	var wg sync.WaitGroup
	wg.Add(cap(data.Individuals))

	// Randomize individual data
	var gRand Randomizer
	gRand.Initialize()

	// Generate random individual data
	for i := 0; i < cap(data.Individuals); i++ {
		go gRand.Assign(ch)
	}

	// Append it to data.
	// data was initially created with 0 length so this appends from the beginning
	go func() {
		for details := range ch {
			data.Individuals = append(data.Individuals, details)
			wg.Done() // decrease waitgroup count
		}
	}()
	wg.Wait()
	close(ch)
}
func Drafted(want int, list <-chan individual) {
	complete := make(chan bool)
	go func(list <-chan individual) {
		var got int
		for drafted := range list {
			got++
			if got <= want {
				fmt.Println(drafted.Age, " with count of ", got)
			} else {
				break
			}
		}
		complete <- true
	}(list)
	<-complete
	close(complete)
}

// StreamIndividuals will output into a channel that accepts individual data
func StreamIndividuals(data *DataStore) <-chan individual {
	out := make(chan individual)
	go func() {
		for _, indi := range data.Individuals {
			out <- indi
		}
		close(out)
	}()
	return out
}

// Draft selects a specific age group and delivers it to a channel
// if filters out ignoring the age group that doesn't fit the criteria
func Draft(individuals <-chan individual) <-chan individual {
	out := make(chan individual)
	go func() {
		defer close(out)
		for draft := range individuals {
			if draft.Age >= 18 && draft.Age <= 35 {
				out <- draft
			}
		}
	}()
	return out
}

func FanInData(chans ...<-chan individual) <-chan individual {
	out := make(chan individual)
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))

		for _, ch := range chans {
			go func(ch <-chan individual) {
				for indi := range ch {
					out <- indi
				}
				wg.Done()
			}(ch)
		}

		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	var data DataStore
	want := 100
	data.Capacity(100)

	// Build initial data
	BuildDataConcurrently(&data)

	// channel data of individuals for processing
	ch := StreamIndividuals(&data)
	// Split into 2 channels - fan out
	d1 := Draft(ch)
	d2 := Draft(ch)
	// Merge channels - fan in
	list := FanInData(d1, d2)
	// output drafted data
	Drafted(want, list)

	for i := 0; i < 5; i++ {
		fmt.Printf("%d\n", data.Individuals[i])
	}
}
