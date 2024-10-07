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


func LTPRequestToBytes(req *types.LTPRequest) []byte {
	b64EncodedMessage := base64.StdEncoding.EncodeToString([]byte(req.Message))
	return []byte("LTP$1.0$" + 
	b64EncodedMessage + "$" + 
	string(req.Level) + "$" + 
	strconv.FormatBool(req.Save) + 
	"$LTP")
}

func BytesToLTPRequest(b []byte) (*types.LTPRequest, types.ResponseStatus) {
	splitted := strings.Split(string(b), "$")
	if len(splitted) != 6 {
		fmt.Println("Invalid request")
		return &types.LTPRequest{}, types.ParseError
	}

	if splitted[0] != "LTP" || splitted[1] != "1.0" || splitted[5] != "LTP" {
		return &types.LTPRequest{}, types.ParseError
	}

	level := splitted[3]
	if !helpers.Contains(types.AllLogLevels, types.LogLevel(level)) {
		return &types.LTPRequest{}, types.InvalidLevelError
	}

	save, err := strconv.ParseBool(splitted[4])
	if err != nil {
		return &types.LTPRequest{}, types.InvalidSaveError
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(splitted[2])
	if err != nil {
		return &types.LTPRequest{}, types.ParseError
	}

	return &types.LTPRequest{
		Message: string(decodedMessage),
		Level: types.LogLevel(level),
		Save: save,
	}, types.Success

}

func LTPResponseToBytes(res *types.LTPResponse) []byte {
	b64EncodedMessage := base64.StdEncoding.EncodeToString([]byte(res.Message))
	return []byte("LTP$1.0$" + 
	b64EncodedMessage + "$" + 
	strconv.Itoa(int(res.Status)) +
	"$LTP")

}

func BytesToLTPResponse(b []byte) (*types.LTPResponse, error) {
	invalidResponseError := errors.New("invalid response")

	splitted := strings.Split(string(b), "$")
	fmt.Println(splitted)
	if len(splitted) != 5 {
		return &types.LTPResponse{}, invalidResponseError
	}

	if splitted[0] != "LTP" || splitted[1] != "1.0" || splitted[4] != "LTP" {
		return &types.LTPResponse{}, invalidResponseError
	}

	status, err := strconv.Atoi(splitted[3])
	if err != nil {
		return &types.LTPResponse{}, invalidResponseError
	}

	if status < 0 || status > 5 {
		return &types.LTPResponse{}, invalidResponseError
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(splitted[2])
	if err != nil {
		return &types.LTPResponse{}, invalidResponseError
	}

	return &types.LTPResponse{
		Message: string(decodedMessage),
		Status: types.ResponseStatus(status),
	}, nil
}