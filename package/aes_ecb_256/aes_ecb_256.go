package aesecb256

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

func EncryptAes256Ecb(secret string, payload string) (string, error) {
	payloadByte := []byte(payload)
	secretByte := []byte(secret)
	enCode, err := encrypt(payloadByte, secretByte)
	if err != nil {
		return "", fmt.Errorf("key_is_not_valid")
	}
	ciphertext := base64.StdEncoding.EncodeToString(enCode)
	return ciphertext, nil
}

func DecryptAes256Ecb(secret string, encryptStr string) (string, error) {
	secretByte := []byte(secret)
	sDec, err := base64.StdEncoding.DecodeString(encryptStr)
	if err != nil {
		return "", fmt.Errorf("key_is_not_valid")
	}
	plaintextByte, err := decrypt(sDec, secretByte)
	if err != nil {
		return "", fmt.Errorf("key_is_not_valid")
	}
	return string(plaintextByte), nil
}

func GenSecret256() (secret string, err error) {
	bytes := make([]byte, 16) 
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("key_is_not_valid")
	}
	key := hex.EncodeToString(bytes) 
	return key, nil
}

func encrypt(pt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt)  
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct, nil
}

func decrypt(ct, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt) 
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	return pt, nil
}
