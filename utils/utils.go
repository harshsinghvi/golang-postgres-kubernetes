package utils

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

func GetEnv(varNameString string, defaultValue string) string {
	var varValue string
	if varValue = os.Getenv(varNameString); varNameString == "" {
		varValue = defaultValue
	}
	return varValue
}

func GenerateToken(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return ""
	}
	// fmt.Println("Hash to store:", string(hash))
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}
