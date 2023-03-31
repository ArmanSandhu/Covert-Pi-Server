package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
	"time"
	"fmt"
)

const (
	SecretKey = "Th!s!5aS3cre+K3y"
)


func Encrypt(jsonString []byte) string {
	// Generate IV
	iv := make([]byte, aes.BlockSize)
	rand.Seed(time.Now().UnixNano())
	rand.Read(iv)

	// Create AES cipher
	block, err := aes.NewCipher([]byte(SecretKey))
	if err != nil {
		fmt.Println("Error Creating Cipher: ", err.Error())
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	paddedString := pad(jsonString, aes.BlockSize)
	
	ciphertext := make([]byte, aes.BlockSize + len(paddedString))
	copy(ciphertext[:aes.BlockSize], iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedString)
	encodedCipherText := base64.StdEncoding.EncodeToString(ciphertext)

	fmt.Printf("Encoded CipherText: %s\n", encodedCipherText)
	return encodedCipherText
}

func Decrypt(encCmdString string) string {
	decodedCipherText, err := base64.StdEncoding.DecodeString(encCmdString)
	if err != nil {
		fmt.Println("Error Decrypting encrypted Cmd String: ", err.Error())
		panic(err)
	}

	extractedIV := decodedCipherText[:aes.BlockSize]
	block, err := aes.NewCipher([]byte(SecretKey))
	if err != nil {
		fmt.Println("Error Creating Cipher: ", err.Error())
		panic(err)
	}
	
	mode := cipher.NewCBCDecrypter(block, extractedIV)

	plnTxt := make([]byte, len(decodedCipherText) - aes.BlockSize)
	mode.CryptBlocks(plnTxt, decodedCipherText[aes.BlockSize:])
	plnTxt = unpad(plnTxt)
	decryptedString := string(plnTxt)
	fmt.Println("Decrypted String: ", decryptedString)
	return decryptedString
}

func pad(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext) % blockSize
	paddedPlnTxt := make([]byte, len(plaintext) + padding)
	copy(paddedPlnTxt, plaintext)
	for i := len(plaintext); i < len(paddedPlnTxt); i++ {
		paddedPlnTxt[i] = byte(padding)
	}
	return paddedPlnTxt
}

func unpad(paddedPlnTxt []byte) []byte {
	padding := int(paddedPlnTxt[len(paddedPlnTxt) - 1])
	return paddedPlnTxt[:len(paddedPlnTxt) - padding]
}
