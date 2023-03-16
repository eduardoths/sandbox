package main

import (
	"os"
	"strings"
)

func main() {
	args := os.Args
	serverType := args[1]
	port := args[2]

	if strings.ToLower(serverType) == "follower" {
		newFollower(port)
	}

	if strings.ToLower(serverType) == "leader" {
		newLeader(port)
	}

}

func newFollower(port string) {
	f := NewFollower()
	if err := f.ServeHTTP(port); err != nil {
		panic(err)
	}
}

func newLeader(port string) {
	l := NewLeader()
	if err := l.ServeHTTP(port); err != nil {
		panic(err)
	}
}
