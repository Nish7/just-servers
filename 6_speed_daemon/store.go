package main

import "sync"

type Store interface {
	AddPlateRecord(Camera, Plate)
}

type InMemoryStore struct {
	PlateRecords map[Camera][]Plate
	mu           sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		PlateRecords: make(map[Camera][]Plate),
	}
}

func (db *InMemoryStore) AddPlateRecord(cam Camera, record Plate) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.PlateRecords[cam] = append(db.PlateRecords[cam], record)
}
