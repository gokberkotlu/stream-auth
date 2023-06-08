package wshandler

import (
	"log"
	"net/http"
	"runtime"
	"sync"

	faceRecognition "stream-auth-webserver/face-recognition"
	imagedatacont "stream-auth-webserver/image-data-cont"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

var mutex sync.Mutex

// WebSocket handler function
func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	// defer conn.Close()

	go faceRecognition.ConsumeImageRec(conn)

	for {
		// Read the message from the client
		_, imageData, err := conn.ReadMessage()

		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		// Update the latest image data
		mutex.Lock()
		runtime.LockOSThread()

		var latestImageData []byte = imagedatacont.ImageDataDecoder(imageData)
		faceRecognition.QueueImageRec(latestImageData)
		mutex.Unlock()

		// go faceRecognition.PerformFaceRecognition(latestImageData, conn)
		// faceRecognition.QueueImageRec(latestImageData)
	}
}
