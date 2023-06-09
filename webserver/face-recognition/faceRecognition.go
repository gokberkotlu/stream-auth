package faceRecognition

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/Kagami/go-face"
	"github.com/gorilla/websocket"
)

// Recognize channel
type IRecCh struct {
	id        int
	imageData []byte
	ts        int64
}

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

// var recCh chan []byte = make(chan []byte)
var recCh chan IRecCh = make(chan IRecCh, 1)

var counter int = 0

var lastImageProcessTimestamp int64 = 0

func QueueImageRec(imageData []byte) {
	fmt.Println("Add image to queue!")
	recData := IRecCh{
		id:        counter,
		imageData: imageData,
		ts:        time.Now().Unix(),
	}

	counter += 1

	recCh <- recData
}

func ConsumeImageRec(wsConn *websocket.Conn) {
	fmt.Println("Start to waiting new images!")
	for recData := range recCh {
		go PerformFaceRecognition(recData, wsConn)
	}
}

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

		labels = append(labels, imgName)
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

func PerformFaceRecognition(recData IRecCh, wsConn *websocket.Conn) {
	fmt.Println(time.Now())
	userFace, err := Rec.RecognizeSingleCNN(recData.imageData)
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
			// recCh <- imageData
			fmt.Println("name:", labels[ID])

			fmt.Println(userFace.Rectangle)
			fmt.Println(userFace.Shapes)

			fmt.Println("RECTANGLE", reflect.TypeOf(userFace.Rectangle))

			wsRes := IRecRes{
				Rect: userFace.Rectangle,
				Name: labels[ID],
				Id:   recData.id,
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

			lastImageProcessTimestamp = recData.ts
			fmt.Println("(*updated)lastImageProcessTimestamp", lastImageProcessTimestamp)
		} else {
			notRecognizedText := "Not a single face on the image"
			fmt.Println(notRecognizedText)
			sendWsRes(wsConn, []byte(notRecognizedText))
		}
	}
	fmt.Println("COUNTER:", recData.id)
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
