package imagedatacont

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func ImageDataDecoder(imageData []byte) []byte {
	encodedBase64StrImageData := string(imageData)
	// fmt.Println(imageData)

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

func SaveImage(imageData []byte, userName string) {
	decodedData := ImageDataDecoder(imageData)

	salt := "_#^!?_"
	fileName := fmt.Sprintf("./images/%s%s%s.jpg", userName, salt, strconv.FormatInt(time.Now().Unix(), 10))

	// Save the decoded data to a file for verification
	err := ioutil.WriteFile(fileName, decodedData, 0644)
	if err != nil {
		fmt.Println("Error saving image:", err)
	}
}
