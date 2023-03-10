package main

type Storage struct {
	memory     map[string]string
	finishChan chan struct{}
	saveChan   chan struct {
		key   string
		value string
	}
	getIncomingChan  chan string
	getOutcomingChan chan string
}

func NewStorage() *Storage {
	s := &Storage{
		memory:     make(map[string]string),
		finishChan: make(chan struct{}, 1),
		saveChan: make(chan struct {
			key   string
			value string
		}, 1),
		getIncomingChan:  make(chan string, 1),
		getOutcomingChan: make(chan string, 1),
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
			s.memory[toSave.key] = toSave.value
		case key := <-s.getIncomingChan:
			s.getOutcomingChan <- s.memory[key]
		default:
		}
	}
}

func (s *Storage) Save(key, value string) {
	s.saveChan <- struct {
		key   string
		value string
	}{
		key:   key,
		value: value,
	}
}

func (s *Storage) Get(key string) string {
	s.getIncomingChan <- key
	return <-s.getOutcomingChan
}

func (s *Storage) Shutdown() {
	s.finishChan <- struct{}{}
}
