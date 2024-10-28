package main

import (
	"net"
	"sync"
)

type UserMap struct {
	userMap map[string]net.Conn
	mu      sync.Mutex
}

func NewUsersMap() *UserMap {
	return &UserMap{
		userMap: make(map[string]net.Conn),
	}
}

func (um *UserMap) AddUser(nickname string, conn net.Conn) {
	um.mu.Lock()
	um.userMap[nickname] = conn
	um.mu.Unlock()
}

func (um *UserMap) GetNicknames() []string {
	var roomMembers []string
	um.mu.Lock()
	for key := range um.userMap {
		roomMembers = append(roomMembers, key)
	}
	um.mu.Unlock()
	return roomMembers
}

func (um *UserMap) getConnection(nickname string) (conn net.Conn, b bool) {
	um.mu.Lock()
	conn, ok := um.userMap[nickname]
	um.mu.Unlock()
	return conn, ok
}
