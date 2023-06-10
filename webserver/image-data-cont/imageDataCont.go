package imagedatacont

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type ICachedUserImages struct {
	FileName  string
	ImageData []byte
}

var CachedUsers = map[string][]ICachedUserImages{}

var Salt string = "_#^!?_"

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

func RegisterUser(imageData []byte, userName string) bool {
	fileName := fmt.Sprintf("./images/%s%s%s.jpg", userName, Salt, strconv.FormatInt(time.Now().Unix(), 10))

	CachedUsers[userName] = append(CachedUsers[userName], ICachedUserImages{
		FileName:  fileName,
		ImageData: imageData,
	})

	if len(CachedUsers[userName]) == 5 {
		for _, userImage := range CachedUsers[userName] {
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
