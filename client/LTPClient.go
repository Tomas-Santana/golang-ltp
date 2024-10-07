package client

import (
	"net"

	"github.com/tomas-santana/ltp/conversion"
	"github.com/tomas-santana/ltp/types"
)

type LTPClient struct {
	ServerAddr string
}

func (c *LTPClient) SendRequest(req *types.LTPRequest) (*types.LTPResponse, error) {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return &types.LTPResponse{}, err
	}
	defer conn.Close()

	reqBytes := conversion.LTPRequestToBytes(req)

	if len(reqBytes) > 2048 {
		return &types.LTPResponse{}, types.ErrRequestTooLong
	}

	_, err = conn.Write(reqBytes)

	if err != nil {
		return &types.LTPResponse{}, err
	}

	buf := make([]byte, 2048)

	n, err := conn.Read(buf)

	if err != nil {
		return &types.LTPResponse{}, err
	}

	resCom, err := conversion.BytesToLTPResponse(buf[:n])

	if err != nil {
		return &types.LTPResponse{}, err
	}

	return resCom, nil
}
 