package client

import (
	"net"

	"github.com/tomas-santana/ltp/conversion"
	"github.com/tomas-santana/ltp/types"
)

type LTPClient struct {
	ServerAddr string
}

func (c *LTPClient) SendRequest(req *types.Request) (*types.Response, error) {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return &types.Response{}, err
	}
	defer conn.Close()

	reqBytes := conversion.RequestToBytes(req)

	if len(reqBytes) > 2048 {
		return &types.Response{}, types.ErrRequestTooLong
	}

	_, err = conn.Write(reqBytes)

	if err != nil {
		return &types.Response{}, err
	}

	buf := make([]byte, 2048)

	n, err := conn.Read(buf)

	if err != nil {
		return &types.Response{}, err
	}

	resCom, err := conversion.BytesToResponse(buf[:n])

	if err != nil {
		return &types.Response{}, err
	}

	return resCom, nil
}
 