package client

import (
	"net"
	"fmt"

	"github.com/tomas-santana/ltp/conversion"
	"github.com/tomas-santana/ltp/types"
)

type UDPClient struct {
	Address string
}

func (uc *UDPClient) Send(req *types.Request) (*types.Response, error) {
	conn, err := net.Dial("udp", uc.Address)
	if err!= nil {
		fmt.Println(err.Error())
    return &types.Response{}, err
  }

	defer conn.Close()

	reqByte := conversion.RequestToBytes(req)

	if len(reqByte) > 2048 {
		return &types.Response{}, types.ErrRequestTooLong
	}

	_, err = conn.Write(reqByte)

	if err!= nil {
		fmt.Println(err.Error())
    return &types.Response{}, err
  }

	buffer := make([]byte, 2048)

	n, err := conn.Read(buffer)

	if err!= nil {
		fmt.Println(err.Error())
    return &types.Response{}, err
  }

	resMessage, err := conversion.BytesToResponse(buffer[:n])

	if err!= nil {
		fmt.Println(err.Error())
    return &types.Response{}, err
  }
	return resMessage, nil
}