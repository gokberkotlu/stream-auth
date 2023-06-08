package main

import (
	"log"
	"net/http"

	faceRecognition "stream-auth-webserver/face-recognition"
	wshandler "stream-auth-webserver/ws-handler"
)

func main() {
	// init recognizer
	faceRecognition.InitImgDb()
	defer faceRecognition.Rec.Close()

	// init ws
	http.HandleFunc("/ws", wshandler.WebsocketHandler)

	// start server on port 8080
	log.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
