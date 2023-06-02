package faceRecognition

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Kagami/go-face"
)

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
	// defer Rec.Close()

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

func PerformFaceRecognition(imageData []byte) {
	fmt.Println("REC", Rec)
	userFace, err := Rec.RecognizeSingleCNN(imageData)
	if err != nil {
		fmt.Printf("Can't recognize: %v\n", err)
	}

	if userFace == nil {
		fmt.Println("Not a single face on the image")
	} else {
		ID := Rec.ClassifyThreshold(userFace.Descriptor, 0.3)
		fmt.Println("ID:", ID)

		if ID != -1 {
			fmt.Println("name:", labels[ID])
		}
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
