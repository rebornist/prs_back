package common

import (
	"encoding/base32"

	"prs/configs"
)

func Signing(s string) (string, error) {
	byteStr, err := configs.EncryptionAESKey(s)
	if err != nil {
		return "", err
	}
	signedtext := base32.StdEncoding.EncodeToString(byteStr)
	return signedtext, nil
}

func Unsigning(s string) (string, error) {
	unsignedByte, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	byteStr, err := configs.DecryptionAESKey(string(unsignedByte))
	if err != nil {
		return "", err
	}

	var rebyte []byte
	for _, b := range byteStr {
		if b != 0 {
			rebyte = append(rebyte, b)
		}
	}

	return string(rebyte), nil
}
