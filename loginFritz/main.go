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

func retrieveData(sid string) {

	resp, err := http.Get(targetURL + "//internet/inetstat_counter.lua?sid=" + sid)
	defer resp.Body.Close() // close connection at end of func
	handleError(err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("RCV: %+v\n", string(bodyBytes))

}

func retrieveSID() string {

	resp, err := http.Get(targetURL + "/login_sid.lua")
	defer resp.Body.Close() // close connection at end of func
	handleError(err)

	if resp.StatusCode != 200 {
		fmt.Printf("%s %s", resp.Status, resp.StatusCode)
		fmt.Printf("*** Data ***\n%v\n", resp)
		os.Exit(42)
	} else {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Body: %+v\n", string(bodyBytes))

		response := &loginResponse{}
		if err := xml.Unmarshal(bodyBytes, response); err != nil {
			panic(err)
		}
		fmt.Printf("Response: %+v\n", response)

		if response.SID == "0000000000000000" {

			fmt.Printf("Challenge: %s\n", response.Challenge)

			decode := response.Challenge + "-" + password
			fmt.Printf("Challenge with PWD: %s\n", decode)

			hasher := md5.New()
			hasher.Write([]byte(string(decode)))
			enc := hex.EncodeToString(hasher.Sum(nil))

			fmt.Printf("Challenge with md5sum: %s\n", enc)

			response_bf := response.Challenge + "-" + enc

			return string(response_bf)

		} else {
			return response.SID
		}
	}
	return ""
}

func login(sid string) {

	resp, err := http.Get(targetURL + "/login_sid.lua?&response=" + sid)
	handleError(err)
	defer resp.Body.Close() // close connection at end of func

	if resp.StatusCode != 200 {
		fmt.Printf("%s %s", resp.Status, resp.StatusCode)
		fmt.Printf("*** Data ***\n%v\n", resp)
		os.Exit(42)
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body: %s\n", body)

}

func main() {

	targetURL = os.Getenv("FRITZ_URL")
	password = os.Getenv("FRITZ_PWD")

	if len(targetURL) == 0 || len(password) == 0 {
		fmt.Println("Environment variable FRITZ_URL and/or FRITZ_PWD not set")
		os.Exit(42)
	}

	sid := retrieveSID()
	fmt.Printf("SID: %s\n", sid)

	login(sid)

	retrieveData(sid)

}
