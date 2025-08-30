package globalFunc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fullstack-journal/app/config"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Hash(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyHash(hashed, plainText string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plainText))
}

func GetJwt(payload interface{}) (string, error) {
	env, err := config.GetEnv()
	if err != nil {
		log.Fatal("Error load env", err)
	}

	key, ok := env["appKey"].(string)
	if !ok {
		log.Fatal("appKey value not found", err)
	}

	// Convert struct/interface{} to map[string]interface{}
	var claims jwt.MapClaims
	bytes, _ := json.Marshal(payload)
	json.Unmarshal(bytes, &claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}


func DecodeJwt(token string)(jwt.MapClaims, error){
	env, err := config.GetEnv()
	if err != nil{
		log.Fatal("Error load env", err)
	}

	key, ok := env["appKey"].(string)
	if !ok{
		log.Fatal("appKey value not found", err)
	}
	jwtKey := []byte(key)

	decodeToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid Signing Method")
		}
		return jwtKey, nil
	})

	if err != nil{
		return nil, errors.New(err.Error())
	}

	if claims, ok := decodeToken.Claims.(jwt.MapClaims); ok && decodeToken.Valid{
		return claims, nil
	}

	return nil, errors.New("Invalid Token or Expired Token")
}

func Encrypt(plainText string)(string, error){
	env, err := config.GetEnv()

	if err != nil{
		log.Fatal("Error load env", err)
	}
	
	key, ok := env["appKey"].(string)
	if !ok{
		log.Fatal("appKey value not found", err)
	}

	if len(key) != 32{
		log.Fatal("appKey harus 32 digits")
	}

	var secretKey = []byte(key)

	block , err := aes.NewCipher(secretKey)
	if err != nil{
		return "", err
	}

	chipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := chipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil{
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(chipherText[aes.BlockSize:], []byte(plainText))
	return base64.URLEncoding.EncodeToString(chipherText), nil
}

func Decrypt(encryptedText string)(interface{}, error){
	env, err := config.GetEnv()

	if err != nil{
		log.Fatal("Error load env", err)
	}
	
	key, ok := env["appKey"].(string)
	if !ok{
		log.Fatal("appKey value not found", err)
	}

	if len(key) != 32{
		log.Fatal("appKey harus 32 digits")
	}

	var secretKey = []byte(key)
	cipherText, err := base64.URLEncoding.DecodeString(encryptedText)
	if err != nil{
		return "",err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil{
		return "",err
	}

	if len(cipherText) < aes.BlockSize{
		return "", errors.New("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	// if number is id and must decode to int64
	id, err := strconv.Atoi(string(cipherText))
	if err != nil{
		// is decode string
		return string(cipherText), nil
	}
	// is number
	return int64(id), nil
}

func UploadFile(dstDir string, file *multipart.FileHeader)(string, error){
	ext := filepath.Ext(file.Filename)
	randomFileName := uuid.New().String() + ext
	dstPath := filepath.Join(dstDir, randomFileName)

	src, err := file.Open()
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", errors.New("failed copy file")
	}

	return randomFileName, nil
}

func BuildSet(arr []string) map[string]bool{
	set := make(map[string]bool)
	for _, v := range arr {
		set[v] = true
	}

	return set
}

func Contains(arr []int, target int) bool {
    for _, v := range arr {
        if v == target {
            return true
        }
    }
    return false
}