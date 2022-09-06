package nym

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("nym")

func GetSelfAddress(conn *websocket.Conn) string {
	selfAddressRequest, err := json.Marshal(map[string]string{"type": "selfAddress"})
	if err != nil {
		panic(err)
	}

	if err = conn.WriteMessage(websocket.TextMessage, []byte(selfAddressRequest)); err != nil {
		panic(err)
	}

	responseJSON := make(map[string]interface{})
	err = conn.ReadJSON(&responseJSON)
	if err != nil {
		panic(err)
	}

	return responseJSON["address"].(string)
}

func GetConn(uri string) (error, *websocket.Conn) {
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return err, nil
	} else {
		return nil, conn
	}

}
func SendText(conn *websocket.Conn, address string, message string, withReplySurb bool) error {
	sendRequest, err := json.Marshal(map[string]interface{}{
		"type":          "send",
		"recipient":     address,
		"message":       message,
		"withReplySurb": withReplySurb,
	})
	if err != nil {
		return err
	}

	log.Infof("sending '%v' over the mix network to %v...", message, address)
	return conn.WriteMessage(websocket.TextMessage, []byte(sendRequest))
}

func ReceiveText(conn *websocket.Conn) (error, string) {
	log.Infof("waiting to receive a message from the mix network...")
	_, receivedMessage, err := conn.ReadMessage()
	if err != nil {
		return err, ""
	}
	return nil, string(receivedMessage)
}

// request tags
const sendRequestTag = 0x00
const replyRequestTag = 0x01
const selfAddressRequestTag = 0x02

// response tags
const errorResponseTag = 0x00
const receivedResponseTag = 0x01
const selfAddressResponseTag = 0x02

func makeSelfAddressRequest() []byte {
	return []byte{selfAddressRequestTag}
}

func parseSelfAddressResponse(rawResponse []byte) []byte {
	if len(rawResponse) != 97 || rawResponse[0] != selfAddressResponseTag {
		panic("Received invalid response")
	}
	return rawResponse[1:]
}

func makeSendRequest(recipient []byte, message []byte, withReplySurb bool) []byte {
	messageLen := make([]byte, 8)
	binary.BigEndian.PutUint64(messageLen, uint64(len(message)))

	surbByte := byte(0)
	if withReplySurb {
		surbByte = 1
	}

	out := []byte{sendRequestTag, surbByte}
	out = append(out, recipient...)
	out = append(out, messageLen...)
	out = append(out, message...)

	return out
}

func makeReplyRequest(message []byte, replySURB []byte) []byte {
	messageLen := make([]byte, 8)
	binary.BigEndian.PutUint64(messageLen, uint64(len(message)))

	surbLen := make([]byte, 8)
	binary.BigEndian.PutUint64(surbLen, uint64(len(replySURB)))

	out := []byte{replyRequestTag}
	out = append(out, surbLen...)
	out = append(out, replySURB...)
	out = append(out, messageLen...)
	out = append(out, message...)

	return out
}

func parseReceived(rawResponse []byte) ([]byte, []byte) {
	if rawResponse[0] != receivedResponseTag {
		panic("Received invalid response!")
	}

	hasSurb := false
	if rawResponse[1] == 1 {
		hasSurb = true
	} else if rawResponse[1] == 0 {
		hasSurb = false
	} else {
		panic("malformed received response!")
	}

	data := rawResponse[2:]
	if hasSurb {
		surbLen := binary.BigEndian.Uint64(data[:8])
		other := data[8:]

		surb := other[:surbLen]
		msgLen := binary.BigEndian.Uint64(other[surbLen : surbLen+8])

		if len(other[surbLen+8:]) != int(msgLen) {
			panic("invalid msg len")
		}

		msg := other[surbLen+8:]
		return msg, surb
	} else {
		msgLen := binary.BigEndian.Uint64(data[:8])
		other := data[8:]

		if len(other) != int(msgLen) {
			panic("invalid msg len")
		}

		msg := other[:msgLen]
		return msg, nil
	}
}

func sendBinary(conn *websocket.Conn, address string, filename string) error {
	readData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	sendRequest := makeSendRequest([]byte(address), readData, false)
	fmt.Printf("sending content of 'dummy file' over the mix network...\n")
	return conn.WriteMessage(websocket.BinaryMessage, sendRequest)
}

func receiveBinary(conn *websocket.Conn) (error, []byte, []byte) {
	log.Infof("waiting to receive a message from the mix network...\n")
	_, receivedResponse, err := conn.ReadMessage()
	if err != nil {
		return err, nil, nil
	}
	fileData, replySURB := parseReceived(receivedResponse)
	return nil, fileData, replySURB
	//if replySURB != nil {
	//	panic("did not expect a replySURB!")
	//}
	//fmt.Printf("writing the file back to the disk!\n")
	//ioutil.WriteFile("received_file_no_reply", fileData, 0644)

}
