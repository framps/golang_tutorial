// Samples used in a small go tutorial
//
// Original code from https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
// Password Hash & Salt Using Golang from James Cox
// with some minor modifications
//
// See github.com/framps/golang_tutorial for latest code

package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	for {

		// Enter a password and generate a salted hash
		pwd := getPwd()
		hash := hashAndSalt(pwd)
		fmt.Println("Hash", hash)

		// Enter the same password again and compare it with the
		// first password entered
		pwd2 := getPwd()
		pwdMatch := comparePasswords(hash, pwd2)

		fmt.Println("Passwords Match?", pwdMatch)

	}

}

func getPwd() []byte {

	// Prompt the user to enter a password
	fmt.Println("Enter a password")

	// We will use this to store the users input
	var pwd string

	// Read the users input
	_, err := fmt.Scan(&pwd)
	if err != nil {
		log.Println(err)
	}

	// Return the users input as a byte slice which will save us
	// from having to do this conversion later on
	return []byte(pwd)

}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {

	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true

}
