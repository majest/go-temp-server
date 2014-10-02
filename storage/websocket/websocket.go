package websocket

import (
	"code.google.com/p/go.net/websocket"
	log "github.com/cihub/seelog"
	"net/http"
)

type Storage struct {
	messages chan *Message
}

type Message struct {
	Data     string "json:`data`"
	Location string "json:`location`"
}

func (s *Storage) webHandler(ws *websocket.Conn) {

	log.Infof("starting handler")

	defer func() {
		log.Infof("Closing websocket")
		ws.Close()
	}()

	for {
		select {
		case msg := <-s.messages:
			log.Infof("websocket msg: %v", ws)
			websocket.JSON.Send(ws, msg)
		}
	}
}

func (s *Storage) listenHandler(w http.ResponseWriter, req *http.Request) {
	serv := websocket.Server{Handler: websocket.Handler(s.webHandler)}
	serv.ServeHTTP(w, req)
}

func New() *Storage {
	log.Infof("New storage")
	s := &Storage{
		messages: make(chan *Message),
	}

	return s
}

func (s *Storage) Run() {
	http.HandleFunc("/listen", s.listenHandler)
	log.Infof("Starting websocket listener on port 9003")
	err := http.ListenAndServe("localhost:9003", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func (s *Storage) Save(data, location string) {
	log.Infof("websocket msg %+s", data)
	m := &Message{Data: data, Location: location}
	s.messages <- m
}
