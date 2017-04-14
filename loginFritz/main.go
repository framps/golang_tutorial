// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// See github.com/framps/golang_tutorial for latest code
//
// challenge response handling with AVM Fritz in go
// See https://www.linux-tips-and-tricks.de/en/programming/389-read-data-from-a-fritzbox-7390-with-python-and-bash
// for a python and bash implementation

// Note - still not working :-(

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"unicode/utf16"
)

var targetURL, password string

type loginResponse struct {
	SID       string
	Challenge string
	BlockTime string
	Rights    string
}

func utf16toString(b []uint8) (string, error) {
	if len(b)&1 != 0 {
		return "", errors.New("len(b) must be even")
	}

	// Check BOM
	var bom int
	if len(b) >= 2 {
		switch n := int(b[0])<<8 | int(b[1]); n {
		case 0xfffe:
			bom = 1
			fallthrough
		case 0xfeff:
			b = b[2:]
		}
	}

	w := make([]uint16, len(b)/2)
	for i := range w {
		w[i] = uint16(b[2*i+bom&1])<<8 | uint16(b[2*i+(bom+1)&1])
	}
	return string(utf16.Decode(w)), nil
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func retrieveChallenge() string {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", targetURL+"/login_sid.lua", nil)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "text/plain")
	res, err := client.Do(req)
	handleError(err)

	if res.StatusCode != 200 {
		fmt.Printf("%s %s", res.Status, res.StatusCode)
		fmt.Printf("*** Data ***\n%v\n", res)
		os.Exit(42)
	} else {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("Body: %+v\n", string(bodyBytes))

		response := &loginResponse{}
		if err := xml.Unmarshal(bodyBytes, response); err != nil {
			panic(err)
		}
		fmt.Printf("Response: %+v\n", response)

		if response.SID == "0000000000000000" {

			response.Challenge = "dd23b85f"
			fmt.Printf("Challenge: %s\n", response.Challenge)

			decode := response.Challenge + "-" + password
			decodeRune := []rune(decode)
			/*
					converter := latinx.Get(latinx.ISO_8859_1)
					// convert a stream of ISO_8859_1 bytes to UTF-8
					utf8bytes, err := converter.Decode([]byte(decode))
					handleError(err)

					fmt.Printf("Challenge + pwd decode: %s\n", utf8bytes)

					utf8Decoded := []rune{}
					for _, r := range utf8bytes {
						utf8Decoded = append(utf8Decoded, rune(r))
					}

					utf16Encoded := utf16.Encode(utf8Decoded)

					u := make([]byte, 0)
					for _, u16 := range utf16Encoded {
						b := make([]byte, 2)
						binary.BigEndian.PutUint16(b, u16)
						u = append(u, b...)
					}
				fmt.Printf("Challenge + pwd decode: %x\n", u)
			*/

			hasher := md5.New()
			hasher.Write([]byte(string(decodeRune)))
			enc := hex.EncodeToString(hasher.Sum(nil))

			fmt.Printf("M: %s\n", enc)

			response_bf := response.Challenge + "-" + enc

			return string(response_bf)

		} else {
			return response.SID
		}
	}
	return ""
}

func login(challenge string) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", targetURL+"/login_sid.lua&response="+challenge, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	handleError(err)

	if res.StatusCode != 200 {
		fmt.Printf("%s %s", res.Status, res.StatusCode)
		fmt.Printf("*** Data ***\n%v\n", res)
		os.Exit(42)
	}

}

func main() {

	targetURL = os.Getenv("FRITZ_URL")
	password = os.Getenv("FRITZ_PWD")

	if len(targetURL) == 0 || len(password) == 0 {
		fmt.Println("Environment variable FRITZ_URL and/or FRITZ_PWD not set")
		os.Exit(42)
	}

	challenge := retrieveChallenge()
	fmt.Printf("SID: %s\n", challenge)

	login(challenge)

}
