package imagedatacont

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type ICachedUserImages struct {
	FileName  string
	ImageData []byte
}

var cachedUsers = map[string][]ICachedUserImages{}

var Salt string = "_#^!?_"

func ImageDataDecoder(imageData []byte) []byte {
	encodedBase64StrImageData := string(imageData)

	// Split the data URI to extract the base64-encoded image data
	parts := strings.Split(encodedBase64StrImageData, ",")
	if len(parts) != 2 {
		fmt.Println("Invalid image data URI")
	}

	// Extract the base64-encoded image data
	base64Data := parts[1]

	// Decode the base64-encoded data
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Println("Error decoding base64 data:", err)
	}

	return decodedData
}

func ImageDataEncodedBuffer(imageData []byte) []byte {
	decodedData := ImageDataDecoder(imageData)

	encodedBuffer := []byte(decodedData)

	return encodedBuffer
}

func RegisterUser(imageData []byte, userName string) bool {
	fileName := fmt.Sprintf("./images/%s%s%s.jpg", userName, Salt, strconv.FormatInt(time.Now().Unix(), 10))

	cachedUsers[userName] = append(cachedUsers[userName], ICachedUserImages{
		FileName:  fileName,
		ImageData: imageData,
	})

	if len(cachedUsers[userName]) == 5 {
		for _, userImage := range cachedUsers[userName] {
			SaveImage(userImage)
		}

		return true
	}

	return false
}

func SaveImage(userImage ICachedUserImages) {
	decodedData := ImageDataDecoder(userImage.ImageData)

	// Save the decoded data to a file for verification
	err := ioutil.WriteFile(userImage.FileName, decodedData, 0644)
	if err != nil {
		fmt.Println("Error saving image:", err)
	}
}

func CreateImagesDirectory() {
	imagesFolderExists, _ := imagesDirectoryExists()

	if !imagesFolderExists {
		os.Mkdir("./images", 0755)
	}
}

func imagesDirectoryExists() (bool, error) {
	path := "./images"
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ClearCachedUserData(username string) {
	delete(cachedUsers, username)
}
