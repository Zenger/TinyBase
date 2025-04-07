/*
 * SPDX-License-Identifier: MIT
 * Author: Zenger (https://github.com/Zenger)
 */

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string, salt string) (string, error) {
	saltAndPepper := pwd + ":" + salt
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(saltAndPepper), bcrypt.DefaultCost)
	return string(hashBytes), err
}

func CheckPassword(pwd string, hash string, salt string) bool {
	saltAndPepper := pwd + ":" + salt

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltAndPepper))
	return err == nil
}

func GenerateHash() string {
	return rand.Text()
}

func GenerateSalt() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
