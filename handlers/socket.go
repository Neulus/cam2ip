package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"

	"github.com/Neulus/cam2ip/image"
	"github.com/Neulus/cam2ip/reader"
)

// Socket handler.
type Socket struct {
	reader reader.ImageReader
	delay  int
}

// NewSocket returns new socket handler.
func NewSocket(reader reader.ImageReader, delay int) *Socket {
	return &Socket{reader, delay}
}

// ServeHTTP handles requests on incoming connections.
func (s *Socket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("socket: accept: %v", err)
		return
	}

	ctx := context.Background()

	for {
		img, err := s.reader.Read()
		if err != nil {
			log.Printf("socket: read: %v", err)
			break
		}

		w := new(bytes.Buffer)

		err = image.NewEncoder(w).Encode(img)
		if err != nil {
			log.Printf("socket: encode: %v", err)
			continue
		}

		b64 := image.EncodeToString(w.Bytes())

		err = conn.Write(ctx, websocket.MessageText, []byte(b64))
		if err != nil {
			break
		}

		time.Sleep(time.Duration(s.delay) * time.Millisecond)
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
