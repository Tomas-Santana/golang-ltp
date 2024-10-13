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

type UDPServer struct {
	ListenAddr string
	Format func(*types.Request) []byte
	WriteStream io.Writer
}

func (u *UDPServer) UDPStart() error {

	addr, err := net.ResolveUDPAddr("udp", u.ListenAddr)
	if err!= nil {
    return err
  }

	conn, err := net.ListenUDP("udp", addr)
	fmt.Println("Listening on ", addr)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err!= nil {
      fmt.Println(err.Error())
      return err
    }
		go u.HandleNewMessage(conn, clientAddr, buffer[:n])
	}
}

func (u *UDPServer) HandleNewMessage(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {

	req, status := conversion.BytesToRequest(buf)
	fmt.Printf("Received message from %s: %s\n", addr.String(), req.Message)
	

	if status != types.Success {
		invalidRequest := types.Request{
			Message: "Invalid request: reason: " + status.String(),
			Level: types.Error,
			Save: false,
		}
		message := u.Format(&invalidRequest)
		fmt.Println(string(message))
		conn.Write(conversion.ResponseToBytes(&types.Response{
			Message: status.String(),
			Status: status,
		}))
		return
	}

	message := u.Format(req)
	fmt.Println(string(message))

	if req.Save {
		u.WriteMessage(message)
	} else {
		fmt.Println(string(message))
	}

	res := types.Response{
		Message: "Message received",
		Status: types.Success,
	}

	response := conversion.ResponseToBytes(&res)
	_, err := conn.WriteToUDP(response, addr)
	if err!= nil {
    fmt.Println(err.Error())
    return
  }
}

func (u *UDPServer) WriteMessage(message []byte) (int, error){
	return u.WriteStream.Write(message)
}

func defaultFormater(req *types.Request) []byte{
	fullLog := time.Now().UTC().String() + " - " + string(req.Level) + " - " + req.Message + "\n"
	return []byte(fullLog)
}

func NewUDPServer(listenAddr string, writeStream io.Writer, format func(*types.Request) []byte) *UDPServer {
	if format == nil {
		format = defaultFormater
	}
	if writeStream == nil {
		writeStream = io.Writer(os.Stdout)
	}

	return &UDPServer{
		ListenAddr: listenAddr,
		WriteStream: writeStream,
		Format: format,
	}
}