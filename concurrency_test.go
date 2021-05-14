package main

import (
	"testing"
)

func TestBuildConcurrentData(t *testing.T) {
	tests := []struct {
		want int
	}{
		{want: 1000},
		{want: 10000},
		{want: 100000},
	}

	for _, test := range tests {
		var data DataStore
		data.Capacity(test.want)
		BuildDataConcurrently(&data)

		gotCapacity := cap(data.Individuals)
		gotLength := len(data.Individuals)
		// length and capacity of data.Individuals should be the same since the data was supposed to fill it
		if gotCapacity != test.want && gotLength != test.want {
			t.Errorf("%s failed want{%d} capacity, got{%d} capacity, got{%d} length", t.Name(), test.want, gotCapacity, gotLength)
		}
	}
}

func TestFirstToDraft(t *testing.T) {
	tests := []struct {
		want          int
		capMultiplier int
	}{
		{want: 100, capMultiplier: 10},
		{want: 10000, capMultiplier: 0},
	}

	for _, test := range tests {
		var data DataStore
		data.Capacity(test.want * test.capMultiplier)

		// Build initial data
		BuildDataSequential(&data)

		// channel data of individuals for processing
		ch := StreamIndividuals(&data)
		// Split into 2 channels - fan out
		d1 := Draft(ch)
		d2 := Draft(ch)
		// Merge channels - fan in
		list := FanInData(d1, d2)
		// output drafted data

		complete := make(chan bool)
		go func(list <-chan individual) {
			var got int
			for drafted := range list {
				got++
				// check ages drafted are 18 through 35 inclusive
				if got <= test.want {
					if drafted.Age < 18 || drafted.Age > 35 {
						t.Errorf("%s failed wanted ages 18 through 35 inclusive, got{%d} drafted age", t.Name(), drafted.Age)
					}
				} else {
					break
				}
				if got > test.want {
					t.Errorf("%s failed got{%d} more than wanted{%d}", t.Name(), got, test.want)
				}
			}
			complete <- true
		}(list)
		<-complete
		close(complete)
	}
}
