package server

import (
	"encoding/json"
	"fmt"
	"github.com/acneto/simple_cache/cache/domain"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	cache      domain.Cache
}

func NewServer(c domain.Cache, listenAddr string) *Server {
	return &Server{
		cache:      c,
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("listen error %s", err)
	}
	defer ln.Close()

	log.Printf("server starting on port [%s]\n", s.listenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("connection accept error %s\n", err)
			return err
		}

		log.Println("server is running...")
		go s.handleConn(conn) // one Goroutine per connection
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		data := make([]byte, 1024)
		size, err := conn.Read(data)
		if err != nil {
			log.Printf("Error reading/client disconnected %s", err)
			return
		}

		message, err := domain.ParseCommand(data[:size])
		if err != nil {
			log.Print(err)
			continue // don't break the connection if client sent a invalid command
		}

		if message.Type == domain.CMDGet {
			err := s.handleGetCmd(conn, *message)
			if err != nil {
				log.Print(err)
			}
			conn.Write([]byte("OK\n"))
		}

		if message.Type == domain.CMDSet {
			err := s.handleSetCmd(*message)
			if err != nil {
				log.Print(err)
			}
			conn.Write([]byte("OK\n"))
		}
	}
}

func (s *Server) handleSetCmd(msg domain.Command) error {
	if err := s.cache.Set(msg.Key, msg.Value); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleGetCmd(conn net.Conn, msg domain.Command) error {
	item, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}
	jsonItem, _ := json.Marshal(item)
	jsonItem = append(jsonItem, '\n')
	_, err = conn.Write(jsonItem)
	return nil
}
