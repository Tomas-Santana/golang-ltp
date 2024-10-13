package conversion

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tomas-santana/ltp/helpers"
	"github.com/tomas-santana/ltp/types"
)


func RequestToBytes(req *types.Request) []byte {
	b64EncodedMessage := base64.StdEncoding.EncodeToString([]byte(req.Message))
	return []byte("LTP$1.0$" + 
	b64EncodedMessage + "$" + 
	string(req.Level) + "$" + 
	strconv.FormatBool(req.Save) + 
	"$LTP")
}

func BytesToRequest(b []byte) (*types.Request, types.ResponseStatus) {
	splitted := strings.Split(string(b), "$")
	if len(splitted) != 6 {
		fmt.Println("Invalid request")
		return &types.Request{}, types.ParseError
	}

	if splitted[0] != "LTP" || splitted[1] != "1.0" || splitted[5] != "LTP" {
		return &types.Request{}, types.ParseError
	}

	level := splitted[3]
	if !helpers.Contains(types.AllLogLevels, types.LogLevel(level)) {
		return &types.Request{}, types.InvalidLevelError
	}

	save, err := strconv.ParseBool(splitted[4])
	if err != nil {
		return &types.Request{}, types.InvalidSaveError
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(splitted[2])
	if err != nil {
		return &types.Request{}, types.ParseError
	}

	return &types.Request{
		Message: string(decodedMessage),
		Level: types.LogLevel(level),
		Save: save,
	}, types.Success

}

func ResponseToBytes(res *types.Response) []byte {
	b64EncodedMessage := base64.StdEncoding.EncodeToString([]byte(res.Message))
	return []byte("LTP$1.0$" + 
	b64EncodedMessage + "$" + 
	strconv.Itoa(int(res.Status)) +
	"$LTP")

}

func BytesToResponse(b []byte) (*types.Response, error) {
	invalidResponseError := errors.New("invalid response")

	splitted := strings.Split(string(b), "$")
	fmt.Println(splitted)
	if len(splitted) != 5 {
		return &types.Response{}, invalidResponseError
	}

	if splitted[0] != "LTP" || splitted[1] != "1.0" || splitted[4] != "LTP" {
		return &types.Response{}, invalidResponseError
	}

	status, err := strconv.Atoi(splitted[3])
	if err != nil {
		return &types.Response{}, invalidResponseError
	}

	if status < 0 || status > 5 {
		return &types.Response{}, invalidResponseError
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(splitted[2])
	if err != nil {
		return &types.Response{}, invalidResponseError
	}

	return &types.Response{
		Message: string(decodedMessage),
		Status: types.ResponseStatus(status),
	}, nil
}