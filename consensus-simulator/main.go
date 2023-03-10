package main

import "fmt"

func main() {
	storage := NewStorage()
	storage.Save("key1", "value1")
	storage.Save("key2", "value2")
	fmt.Println(storage.Get("key1"))
	storage.Shutdown()
}
