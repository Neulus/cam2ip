// Package handlers.
package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/Neulus/cam2ip/image"
	"github.com/Neulus/cam2ip/reader"
)

// MJPEG handler.
type MJPEG struct {
	reader reader.ImageReader
	delay  int
}

// NewMJPEG returns new MJPEG handler.
func NewMJPEG(reader reader.ImageReader, delay int) *MJPEG {
	return &MJPEG{reader, delay}
}

// ServeHTTP handles requests on incoming connections.
func (m *MJPEG) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	mimeWriter := multipart.NewWriter(w)
	mimeWriter.SetBoundary("--boundary")

	w.Header().Add("Connection", "close")
	w.Header().Add("Cache-Control", "no-store, no-cache")
	w.Header().Add("Content-Type", fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary()))

	cn := w.(http.CloseNotifier).CloseNotify()

loop:
	for {
		select {
		case <-cn:
			break loop

		default:
			partHeader := make(textproto.MIMEHeader)
			partHeader.Add("Content-Type", "image/jpeg")

			partWriter, err := mimeWriter.CreatePart(partHeader)
			if err != nil {
				log.Printf("mjpeg: createPart: %v", err)
				continue
			}

			img, err := m.reader.Read()
			if err != nil {
				log.Printf("mjpeg: read: %v", err)
				continue
			}

			err = image.NewEncoder(partWriter).Encode(img)
			if err != nil {
				log.Printf("mjpeg: encode: %v", err)
				continue
			}

			time.Sleep(time.Duration(m.delay) * time.Millisecond)
		}
	}

	mimeWriter.Close()
}
