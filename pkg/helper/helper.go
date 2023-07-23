package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
)

func Encrypt(key, plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	bPlainText := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(bPlainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(bPlainText))

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(key, cryptoText string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText), nil
}

func EchoBindErrorTranslator(err error) string {
	fmt.Println(err)
	switch t := err.(*echo.HTTPError).Unwrap().(type) {
	case *json.UnmarshalTypeError:
		switch t.Type.Kind() {
		case reflect.String:
			return fmt.Sprintf("%s must be a string", t.Field)
		case reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%s must be a number", t.Field)
		default:
			return fmt.Sprintf("%s must be a %s", t.Field, t.Type.Name())
		}
	case *time.ParseError:
		return fmt.Sprintf("%s is not a valid time format", t.Value)
	default:
		return err.Error()
	}
}
