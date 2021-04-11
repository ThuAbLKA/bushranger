package util

import (
	"crypto/rand"
	"fmt"
)

// GenerateId
func GenerateId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	CheckError(err)

	uuid := fmt.Sprintf("%x-%x-%x-%x", b[0:4], b[4:8], b[8:12], b[12:])
	fmt.Println(uuid)
	return uuid
}

// CheckError
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
