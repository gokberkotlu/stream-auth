package faceRecognition

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	imagedatacont "stream-auth-webserver/image-data-cont"
	"strings"
	"time"

	"github.com/Kagami/go-face"
	"github.com/gorilla/websocket"
)

// Recognize response
type IRecRes struct {
	Name string          `json:"name"`
	Rect image.Rectangle `json:"rect"`
	Id   int             `json:"id"`
}

const dataDir = "./"

var (
	modelsDir = filepath.Join(dataDir, "models")
	imagesDir = filepath.Join(dataDir, "images")
)

var Rec *face.Recognizer
var faces []face.Face = []face.Face{}
var labels []string = []string{}

func InitImgDb() {
	var err error
	Rec, err = face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}

	for _, imgName := range GetImageList() {
		refImage := filepath.Join(imagesDir, imgName)

		recognizedFaces, err := Rec.RecognizeFile(refImage)

		if err != nil {
			log.Fatalf("Can't recognize: %v", err)
		}

		faces = append(faces, recognizedFaces...)

		userName := strings.Split(imgName, imagedatacont.Salt)[0]

		labels = append(labels, userName)
	}

	var samples []face.Descriptor
	var ids []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		ids = append(ids, int32(i))
	}

	Rec.SetSamples(samples, ids)

	fmt.Println("ids", ids)
	fmt.Println("Rec val:", Rec)
}

func PerformFaceRecognition(imageData []byte, wsConn *websocket.Conn) {
	fmt.Println(time.Now())
	userFace, err := Rec.RecognizeSingleCNN(imageData)
	if err != nil {
		fmt.Printf("Can't recognize: %v\n", err)
	}
	fmt.Println(time.Now())

	if userFace == nil {
		noFaceDetectedMessage := "Not a single face on the image"
		fmt.Println(noFaceDetectedMessage)
		sendWsRes(wsConn, []byte(noFaceDetectedMessage))
	} else {
		ID := Rec.ClassifyThreshold(userFace.Descriptor, 0.4)
		fmt.Println("ID:", ID)

		if ID != -1 {
			fmt.Println(userFace.Rectangle)
			fmt.Println(userFace.Shapes)

			wsRes := IRecRes{
				Rect: userFace.Rectangle,
				Name: labels[ID],
			}

			// Convert the rectangle to JSON
			rectJSON, err := json.Marshal(wsRes)
			if err != nil {
				errMessage := fmt.Sprintf("Failed to convert image.Rectangle to JSON: %v", err)
				sendWsRes(wsConn, []byte(errMessage))
				return
			}

			// Send a response back to the client
			err = wsConn.WriteMessage(websocket.TextMessage, rectJSON)
			if err != nil {
				log.Printf("Failed to send response to WebSocket client: %v", err)
			}
		} else {
			notRecognizedText := "User not identified"
			fmt.Println(notRecognizedText)
			sendWsRes(wsConn, []byte(notRecognizedText))
		}
	}
}

func CheckFaceForRegistration(rawImageData []byte, encodedImageDataBuffer []byte, userName string, wsConn *websocket.Conn) {
	fmt.Println(time.Now())
	userFace, err := Rec.RecognizeSingleCNN(encodedImageDataBuffer)
	if err != nil {
		fmt.Printf("Can't recognize: %v\n", err)
	}
	fmt.Println(time.Now())

	if userFace == nil {
		noFaceDetectedMessage := "Not a single face on the image"
		fmt.Println(noFaceDetectedMessage)
		sendWsRes(wsConn, []byte(noFaceDetectedMessage))
	} else {
		ID := Rec.ClassifyThreshold(userFace.Descriptor, 0.4)
		fmt.Println("ID:", ID)

		if ID == -1 {
			// Convert the rectangle to JSON
			rectJSON, err := json.Marshal(userFace.Rectangle)
			if err != nil {
				errMessage := fmt.Sprintf("Failed to convert image.Rectangle to JSON: %v", err)
				sendWsRes(wsConn, []byte(errMessage))
				return
			}

			// Save user image to file system if user is not defined before
			registerResult := imagedatacont.RegisterUser(rawImageData, userName)
			// if registration done, add images to Rec
			if registerResult {
				InitImgDb()
			}

			// Send a response back to the client
			err = wsConn.WriteMessage(websocket.TextMessage, rectJSON)
			if err != nil {
				log.Printf("Failed to send response to WebSocket client: %v", err)
			}
		} else {
			notRecognizedText := "User already defined"
			fmt.Println(notRecognizedText)
			sendWsRes(wsConn, []byte(notRecognizedText))
		}
	}
}

func sendWsRes(wsConn *websocket.Conn, message []byte) {
	// Send a response back to the client
	err := wsConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Failed to send response to WebSocket client: %v", err)
	}
}

func GetImageList() []string {
	entries, err := os.ReadDir("./images")
	if err != nil {
		log.Fatal(err)
	}

	var fileList []string
	for _, e := range entries {
		fileList = append(fileList, e.Name())
	}

	return fileList
}
