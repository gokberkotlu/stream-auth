package wshandler

import (
	"fmt"
	"log"
	"net/http"

	"runtime"
	"sync"

	faceRecognition "stream-auth-webserver/face-recognition"
	imagedatacont "stream-auth-webserver/image-data-cont"

	"github.com/gorilla/websocket"
)

var mutex = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocket handler function
func WebsocketFaceRecHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

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

		var latestImageData []byte = imagedatacont.ImageDataEncodedBuffer(imageData)
		mutex.Unlock()

		go faceRecognition.PerformFaceRecognition(latestImageData, conn)
	}
}

func WebsocketFaceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("name")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	// on connection
	fmt.Println("Client connected for registration")
	imagedatacont.ClearCachedUserData(username)

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

		mutex.Unlock()

		// go imagedatacont.SaveImage(imageData, userName)
		go faceRecognition.CheckFaceForRegistration(imageData, imagedatacont.ImageDataEncodedBuffer(imageData), username, conn)
	}
}
