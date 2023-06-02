package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"

	faceRecognition "stream-auth-webserver/face-recognition"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

var latestImageData []byte
var mutex sync.Mutex

// WebSocket handler function
func websocketHandler(w http.ResponseWriter, r *http.Request) {
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
		// latestImageData = imageData
		latestImageData = imageDataDecoder(imageData)
		mutex.Unlock()

		go faceRecognition.PerformFaceRecognition(latestImageData)
	}
}

func imageDataDecoder(imageData []byte) []byte {
	encodedBase64StrImageData := string(imageData)
	// fmt.Println(imageData)

	// Split the data URI to extract the base64-encoded image data
	parts := strings.Split(encodedBase64StrImageData, ",")
	if len(parts) != 2 {
		fmt.Println("Invalid image data URI")
		// return
	}

	// Extract the base64-encoded image data
	base64Data := parts[1]

	// Decode the base64-encoded data
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Println("Error decoding base64 data:", err)
		// return
	}

	// Save the decoded data to a file for verification
	// err = ioutil.WriteFile("image.jpg", decodedData, 0644)
	// if err != nil {
	// 	fmt.Println("Error saving image:", err)
	// 	// return
	// }

	encodedBuffer := []byte(decodedData)

	// fmt.Println(encodedBuffer)

	return encodedBuffer
}

func main() {
	// init recognizer
	faceRecognition.InitImgDb()
	defer faceRecognition.Rec.Close()

	// init ws
	http.HandleFunc("/ws", websocketHandler)

	// start server on port 8080
	log.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
