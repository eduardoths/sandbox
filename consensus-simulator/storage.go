package main

import "fmt"

type Storage struct {
	memory           map[string]StorageData
	finishChan       chan struct{}
	saveChan         chan StorageSaveStruct
	getIncomingChan  chan string
	getOutcomingChan chan StorageData
	deleteChan       chan string
}

type StorageData struct {
	Message  string
	Metadata map[string]interface{}
}

type StorageSaveStruct struct {
	Key   string
	Value StorageData
}

func NewStorage() *Storage {
	s := &Storage{
		memory:           make(map[string]StorageData),
		finishChan:       make(chan struct{}, 1),
		saveChan:         make(chan StorageSaveStruct, 1),
		getIncomingChan:  make(chan string, 1),
		getOutcomingChan: make(chan StorageData, 1),
		deleteChan:       make(chan string, 1),
	}
	go s.start()
	return s
}

func (s *Storage) start() {
	for {
		select {
		case <-s.finishChan:
			return
		case toSave := <-s.saveChan:
			s.memory[toSave.Key] = toSave.Value
		case key := <-s.getIncomingChan:
			s.getOutcomingChan <- s.memory[key]
		case key := <-s.deleteChan:
			delete(s.memory, key)
		default:
		}
	}
}

func (s *Storage) Save(sss StorageSaveStruct) {
	fmt.Println("saving", sss)
	s.saveChan <- sss
}

func (s *Storage) Get(key string) StorageData {
	s.getIncomingChan <- key
	return <-s.getOutcomingChan
}

func (s *Storage) Delete(key string) {
	s.deleteChan <- key
}

func (s *Storage) Shutdown() {
	s.finishChan <- struct{}{}
}
