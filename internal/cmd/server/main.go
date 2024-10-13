package main

import (
	"github.com/tomas-santana/ltp/server"
	"os"
	"io"
)

func main() {

	// create a writer to write logs to a file

	file, err := os.CreateTemp(".", "ltp-*.log")
	
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// need an io.Writer to pass to the server

	writer := io.MultiWriter(file, os.Stdout)

	// s := server.NewLTPServer("localhost:8080", writer, nil)
	// s.Start()

	u := server.NewUDPServer("localhost:8080", writer, nil)
	u.UDPStart()
}
