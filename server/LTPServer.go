package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/tomas-santana/ltp/conversion"
	"github.com/tomas-santana/ltp/types"
)

type LTPServer struct {
	ListenAddr string
	Format func(*types.Request) []byte
	WriteStream io.Writer
}

func (s *LTPServer) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	fmt.Println("Listening on", s.ListenAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.HandleNewMessage(conn)
	}
}


func (s *LTPServer) HandleNewMessage(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from", conn.RemoteAddr().String())

	buf := make([]byte, 2048)

	n, err := conn.Read(buf)

	if err != nil {
		return
	}

	req, status := conversion.BytesToRequest(buf[:n])
	

	if status != types.Success {
		invalidRequest := types.Request{
			Message: "Invalid request: reason: " + status.String(),
			Level: types.Error,
			Save: false,
		}
		message := s.Format(&invalidRequest)
		fmt.Println(string(message))
		conn.Write(conversion.ResponseToBytes(&types.Response{
			Message: status.String(),
			Status: status,
		}))
		return
	}

	message := s.Format(req)

	if req.Save {
		s.WriteMessage(message)
	} else {
		fmt.Println(string(message))
	}

	res := types.Response{
		Message: "Message received",
		Status: types.Success,
	}

	conn.Write(conversion.ResponseToBytes(&res))
}

func (s *LTPServer) WriteMessage(message []byte) (int, error) {
	return s.WriteStream.Write(message)
}

func defaultFormat(req *types.Request) []byte{
	fullLog := time.Now().UTC().String() + " - " + string(req.Level) + " - " + req.Message + "\n"
	return []byte(fullLog)
}

func NewLTPServer(listenAddr string, writeStream io.Writer, format func(*types.Request) []byte) *LTPServer {
	if format == nil {
		format = defaultFormat
	}
	if writeStream == nil {
		writeStream = io.Writer(os.Stdout)
	}

	return &LTPServer{
		ListenAddr: listenAddr,
		WriteStream: writeStream,
		Format: format,
	}
}

