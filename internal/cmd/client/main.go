package main

import (
	"fmt"

	"github.com/tomas-santana/ltp/client"
	"github.com/tomas-santana/ltp/types"
)

func main() {
	c := client.LTPClient{ServerAddr: "localhost:8080"}

	req := &types.LTPRequest{
		Message: "Hello, server!",
		Level: types.Debug,
		Save: true,
	}

	res, err := c.SendRequest(req)

	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}