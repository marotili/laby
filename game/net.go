package game

import (
// "bytes"
// "encoding/gob"
// "fmt"
)

// func (v ActionType) MarshalBinary() ([]byte, error) {
// 	// A simple encoding: plain text.
// 	var b bytes.Buffer
// 	fmt.Fprintln(&b, v.x, v.y, v.z)
// 	return b.Bytes(), nil
// }

type ServerResponse int

const (
	ServerActionWait ServerResponse = iota
	ServerActionOk
	ServerActionDenied
)

type ClientRequest int

const (
	ClientReqSendAction ClientRequest = iota
	ClientReqUpdate
)
