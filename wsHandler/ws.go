package wsHandler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

// 1. client connects to game
// 2. server saves the socket
// 3. other clients are notified about new client
// 4. client sends request to start the game
// 5. other clients notified about game start
// 6. admin client sends first game state
// 7. others get that state
// 8. client notifies that game ended
// 9. on socket close, delete associated resources

type DataPacket struct {
	GameID string                 `json:"id"`
	Cmd    string                 `json:"cmd"`
	Data   map[string]interface{} `json:"data"`
}

func onConnect(packet DataPacket, conn *websocket.Conn) {
	log.Printf("Connect %+v\n", packet)

	usernamesFromCache, exists := UsernamesCache.Get(packet.GameID)
	username := packet.Data["username"].(string)

	var res []string
	if exists {
		res = append(usernamesFromCache.([]string), username)
	} else {
		res = make([]string, 1)
		res[0] = username
	}
	UsernamesCache.Set(packet.GameID, res)

	msg := &DataPacket{
		GameID: packet.GameID,
		Cmd:    "new_user",
		Data: map[string]interface{}{
			"users": res,
		},
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ERROR: couldn't parse data into the json: %+v\n", msg)
	}

	connsFromCache, exists := ConnsCache.Get(packet.GameID)
	var ips []*websocket.Conn
	if exists {
		ips = connsFromCache.([]*websocket.Conn)
		ips = append(ips, conn)
	} else {
		ips = make([]*websocket.Conn, 1)
		ips[0] = conn
	}
	ConnsCache.Set(packet.GameID, ips)

	for _, conn := range ips {
		if err = conn.Write(context.Background(), websocket.MessageText, b); err != nil {
			log.Printf("ERROR: couldn't send data to client: %+v\n", err)
		}
	}
}

func onUpdateState(packet DataPacket) {
	conns, exists := ConnsCache.Get(packet.GameID)
	if !exists {
		log.Printf("ERROR: invaild GameID %s\n", packet.GameID)
	}
	for _, conn := range conns.([]*websocket.Conn) {
		b, err := json.Marshal(packet)
		if err != nil {
			log.Printf("ERROR: could not create JSON: %+v\n", err)
		}
		err = conn.Write(context.TODO(), websocket.MessageText, b)
		if err != nil {
			log.Printf("ERROR: could not send JSON: %+v\n", err)
		}
	}
}

var WsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})

	if err != nil {
		log.Printf("ERROR: couldn't accept conn: %+v\n", err)
		return
	}
	log.Printf("Created new connection with %s\n", r.RemoteAddr)

	go func() {
	MAINLOOP:
		for {
			msgType, msg, err := c.Read(context.Background())

			if err == nil {
				var data DataPacket

				err = json.Unmarshal(msg, &data)
				if err != nil {
					log.Printf("ERROR: parsing json: %+v\n", err)
					log.Printf("json: %+v\n", string(msg))
					continue
				}

				log.Printf("Received msg: %s %+v\n", msgType.String(), data)
				switch data.Cmd {
				case "update":
					onUpdateState(data)
				case "connect":
					onConnect(data, c)
				case "start":
					onUpdateState(data)
				}
			} else {
				break MAINLOOP
			}
		}

		log.Println("Client disconnected")
	}()
})
