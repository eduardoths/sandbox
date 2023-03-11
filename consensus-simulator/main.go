package main

import "fmt"

func main() {
	storage := NewStorage()
	storage.Save(StorageSaveStruct{
		Key: "key1",
		Value: StorageData{
			Message: "value1",
		},
	})
	storage.Save(StorageSaveStruct{
		Key: "key2",
		Value: StorageData{
			Message: "value2",
		},
	})
	fmt.Println(storage.Get("key1"))
	storage.Shutdown()
}
