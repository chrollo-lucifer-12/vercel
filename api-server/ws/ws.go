package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var WsClients = make(map[*websocket.Conn]bool)
var WsMu sync.Mutex

func WsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	WsMu.Lock()
	WsClients[conn] = true
	WsMu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	WsMu.Lock()
	delete(WsClients, conn)
	WsMu.Unlock()
	conn.Close()
}
