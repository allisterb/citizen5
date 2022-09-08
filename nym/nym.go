package nym

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log/v2"
)

type Message struct {
	RawMessage []byte
	Binary     []byte
	Text       string
	Json       map[string]interface{}
	Surb       []byte
	Error      string
}

type Command struct {
}

type SignedCommand struct {
	Command   Command
	Identity  string
	Signature []byte
}

var log = logging.Logger("nym")

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

func parseBinaryResponse(rawResponse []byte) ([]byte, []byte) {
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

func ResponseIsError(rawResponse []byte) bool {
	return rawResponse[0] == errorResponseTag
}

func ResponseIsBinary(rawResponse []byte) bool {
	return rawResponse[0] == receivedResponseTag
}

func ResponseIsSelfAddress(rawResponse []byte) bool {
	return rawResponse[0] == selfAddressResponseTag
}

func GetSelfAddressText(conn *websocket.Conn) string {
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

func GetSelfAddressBinary(conn *websocket.Conn) []byte {
	selfAddressRequest := makeSelfAddressRequest()
	if err := conn.WriteMessage(websocket.BinaryMessage, selfAddressRequest); err != nil {
		panic(err)
	}
	_, receivedResponse, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	return parseSelfAddressResponse(receivedResponse)

}
func GetConn(uri string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	} else {
		return conn, nil
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

func SendBinary(conn *websocket.Conn, address []byte, data []byte) error {
	sendRequest := makeSendRequest(address, data, false)
	log.Infof("sending binary data of length %s over mix network...", len(data))
	return conn.WriteMessage(websocket.BinaryMessage, sendRequest)
}

func SendBinaryFile(conn *websocket.Conn, address []byte, filename string) error {
	readData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	sendRequest := makeSendRequest(address, readData, false)
	log.Infof("sending content of file %s over mix network...", filename)
	return conn.WriteMessage(websocket.BinaryMessage, sendRequest)
}

func Receive(conn *websocket.Conn) ([]byte, error) {
	log.Infof("waiting to receive a message from the mix network...")
	_, receivedMessage, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return receivedMessage, nil
}

func ReceiveMessage(conn *websocket.Conn) (Message, error) {
	msg := Message{}
	r, err := Receive(conn)
	if err != nil {
		log.Errorf("error receiving message from Nym WebSocket connection: %v", err)
		return msg, err
	}
	msg.RawMessage = r
	if ResponseIsError(r) {
		msg.Error = string(r[1:])
		log.Infof("received error response from Nym mix network: %s", msg.Error)
	} else if ResponseIsBinary(r) {
		payload, surb := parseBinaryResponse(r)
		msg.Binary = payload
		msg.Text = string(payload)
		msg.Surb = surb
		if len(payload) == 4 && string(payload) == "ping" {
			reply := makeReplyRequest([]byte("ping"), surb)
			if err = conn.WriteMessage(websocket.BinaryMessage, reply); err != nil {
				log.Errorf("Could not reply to ping message")
			}
		}
	} else {
		var data map[string]interface{}
		if err = json.Unmarshal(r, &data); err != nil {
			log.Infof("message is not JSON")
			return msg, nil
		} else {
			msg.Json = data
			log.Info("received JSON message from Nym mix network.")
		}
	}
	return msg, nil
}

func ReceiveCommand(conn *websocket.Conn) (bool, Command, error) {
	c := Command{}
	msg, err := ReceiveMessage(conn)
	if err != nil {
		return false, c, err
	}
	if msg.Json == nil {
		return false, c, nil
	} else {
		return false, c, nil
	}
}
