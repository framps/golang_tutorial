// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
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

	"github.com/bjarneh/latinx"
)

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

func main() {

	targetURL := os.Getenv("FRITZ_URL")
	password := os.Getenv("FRITZ_PWD")

	if len(targetURL) == 0 || len(password) == 0 {
		fmt.Println("Environment variable FRITZ_URL and/or FRITZ_PWD not set")
		os.Exit(42)
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", targetURL+"/login_sid.lua", nil)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "text/plain")
	res, _ := client.Do(req)

	var body string
	if res.StatusCode == 200 { // OK
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("Body: %+v\n", string(bodyBytes))

		response := &loginResponse{}
		if err := xml.Unmarshal(bodyBytes, response); err != nil {
			panic(err)
		}
		fmt.Printf("Response: %+v\n", response)

		decode := response.Challenge + "-" + password

		converter := latinx.Get(latinx.ISO_8859_1)
		// convert a stream of ISO_8859_1 bytes to UTF-8
		utf8bytes, _ := converter.Decode([]byte(decode))

		utf8Decoded := []rune{}
		for _, r := range utf8bytes {
			utf8Decoded = append(utf8Decoded, rune(r))
		}

		utf16Encoded := utf16.Encode(utf8Decoded)

		hasher := md5.New()
		hasher.Write([]byte(utf16Encoded))
		hex.EncodeToString(hasher.Sum(nil))

		fmt.Printf("%s", utf16.Decode(utf16Encoded))

	} else {
		fmt.Printf("%v\n%v\n", res.StatusCode, body)
		os.Exit(42)
	}
}
