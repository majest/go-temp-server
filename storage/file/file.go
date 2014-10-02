package file

import (
	log "github.com/cihub/seelog"
)

type Storage struct {
	messages chan *Message
}

type Message struct {
	Data     string "json:`data`"
	Location string "json:`location`"
}

func New() *Storage {
	log.Infof("New storage")
	s := &Storage{
		messages: make(chan *Message),
	}

	return s
}

func (s *Storage) Save(data, location string) {
	log.Infof("Msg %+s ,location : %s", data, location)
}

func (s *Storage) Run() {
	<-s.messages
}
