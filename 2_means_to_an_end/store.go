package main

import (
	"fmt"
	"internal/runtime/exithook"
)

type Store interface {
	Query(minTime, maxTime int32) int32
	Insert(timestamp, price int32)
}

type InMemoryStore struct {
	store map[int32]int32 // stores: Timestamp -> Price
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[int32]int32),
	}
}

func (st *InMemoryStore) Query(minTime, maxTime int32) int32 {
	var sum int32 = 0
	var counter int32 = 0
	for k, v := range st.store {
		if k >= minTime && k <= maxTime {
			sum += v
			counter++
		}
	}

	if counter == 0 || sum == 0 {
		return 0
	}

	mean := sum / counter
	fmt.Printf("sum %d and counter %d = mean [%d]", sum, counter, mean)
	return mean
}

func (i *InMemoryStore) Insert(timestamp, price int32) {
	if _, exist := i.store[timestamp]; !exist {
		i.store[timestamp] = price
	}
}
