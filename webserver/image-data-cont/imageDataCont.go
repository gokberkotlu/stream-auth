package imagedatacont

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func ImageDataDecoder(imageData []byte) []byte {
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
