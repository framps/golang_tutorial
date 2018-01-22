// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// See github.com/framps/golang_tutorial for latest code
//
// Logon to AVM Fritz and challenge response handling in go
// This sample code retrieves the byte counters
//
// See https://www.linux-tips-and-tricks.de/en/programming/389-read-data-from-a-fritzbox-7390-with-python-and-bash
// for a python and bash implementation

package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"unicode/utf16"
)

var (
	hostname,
	password,
	fritzURL string
)

type loginResponse struct {
	SID       string
	Challenge string
	BlockTime string
	Rights    string
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func retrieveData(target, url, sid string) []byte {

	fmt.Printf("Retrieving statistic data from %s ...\n", target)

	endpoint := target + url + "?sid=" + sid
	resp, err := http.Get(endpoint)

	defer resp.Body.Close() // close connection at end of func
	handleError(err)

	response, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	return response
}

func createChallengeResponse(server, password string) string {

	resp, err := http.Get("http://" + hostname + "/login_sid.lua")
	handleError(err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s %d", resp.Status, resp.StatusCode)
		handleError(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	var response loginResponse
	err = xml.Unmarshal(bodyBytes, &response)
	handleError(err)

	var zeroSID = regexp.MustCompile(`^0+$`)
	if zeroSID.MatchString(response.SID) {

		combined := response.Challenge + "-" + password

		decoded := []rune(combined)
		encoded := utf16.Encode(decoded)

		buf := new(bytes.Buffer)
		for i := range encoded {
			err = binary.Write(buf, binary.LittleEndian, encoded[i])
			handleError(err)
		}

		hasher := md5.New()
		hasher.Write(buf.Bytes())
		enc := hex.EncodeToString(hasher.Sum(nil))
		responseBf := response.Challenge + "-" + enc

		return string(responseBf)

	}

	if v, e := strconv.Atoi(response.BlockTime); e == nil && v > 0 {
		err = fmt.Errorf("Logon blocked for %s seconds", response.BlockTime)
		handleError(err)
	}

	err = fmt.Errorf("Unknown error occured. Response: %+v", response)
	handleError(err)

	return ""
}

func retrieveSID(targetURL, responseBf string) string {

	resp, err := http.Get(targetURL + "/login_sid.lua?&response=" + responseBf)
	handleError(err)
	defer resp.Body.Close() // close connection at end of func

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s %d", resp.Status, resp.StatusCode)
		handleError(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	response := &loginResponse{}
	err = xml.Unmarshal(body, response)
	handleError(err)

	return response.SID

}

func main() {

	hostname = os.Getenv("FRITZ_HOSTNAME")
	password = os.Getenv("FRITZ_PASSWORD")
	if len(hostname) == 0 || len(password) == 0 {
		fmt.Println("Environment variable FRITZ_HOSTNAME and/or FRITZ_PASSWORD not set")
		os.Exit(1)
	}

	fritzURL = "http://" + hostname

	responseBf := createChallengeResponse(fritzURL, password)
	sid := retrieveSID(fritzURL, responseBf)
	data := retrieveData(fritzURL, "/internet/inetstat_counter.lua", sid)

	// fmt.Printf("%v\n", string(data))

	/*
	   <td datalabel="Online-Zeit (hh:mm)" class="time">17:16</td><td datalabel="Datenvolumen gesamt(MB)" class="vol">3876</td><td datalabel="Datenvolumen gesendet(MB)" class="vol">193</td><td datalabel="Datenvolumen empfangen(MB)" class="vol">3683</td><td datalabel="Verbindungen" class="conn">1</td></tr><tr><td datalabel="" class="first_col">Gestern</td>
	   <td datalabel="Online-Zeit (hh:mm)" class="time">24:00</td><td datalabel="Datenvolumen gesamt(MB)" class="vol">3836</td><td datalabel="Datenvolumen gesendet(MB)" class="vol">229</td><td datalabel="Datenvolumen empfangen(MB)" class="vol">3607</td><td datalabel="Verbindungen" class="conn">1</td></tr><tr><td datalabel="" class="first_col">Aktuelle Woche</td>
	   <td datalabel="Online-Zeit (hh:mm)" class="time">17:16</td><td datalabel="Datenvolumen gesamt(MB)" class="vol">3876</td><td datalabel="Datenvolumen gesendet(MB)" class="vol">193</td><td datalabel="Datenvolumen empfangen(MB)" class="vol">3683</td><td datalabel="Verbindungen" class="conn">1</td></tr><tr><td datalabel="" class="first_col">Aktueller Monat</td>
	   <td datalabel="Online-Zeit (hh:mm)" class="time">401:07</td><td datalabel="Datenvolumen gesamt(MB)" class="vol">448103</td><td datalabel="Datenvolumen gesendet(MB)" class="vol">278255</td><td datalabel="Datenvolumen empfangen(MB)" class="vol">169848</td><td datalabel="Verbindungen" class="conn">39</td></tr><tr><td datalabel="" class="first_col">Vormonat</td>
	   <td datalabel="Online-Zeit (hh:mm)" class="time">742:55</td><td datalabel="Datenvolumen gesamt(MB)" class="vol">354100</td><td datalabel="Datenvolumen gesendet(MB)" class="vol">27789</td><td datalabel="Datenvolumen empfangen(MB)" class="vol">326311</td><td datalabel="Verbindungen" class="conn">36</td></tr></table>
	*/

	var volRegex = regexp.MustCompile(`class="vol">(?P<sum>[\d]+)</td>.+class="vol">(?P<sent>[\d]+)</td>.+class="vol">(?P<received>[\d]+)</td>.+`)

	type usage struct {
		summb  int64
		sentmb int64
		recvmb int64
	}

	statRange := [...]string{"Today", "Yesterday", "Current week", "Current month", "Last month"}

	stats := make([]usage, 0)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		if vol := volRegex.FindStringSubmatch(scanner.Text()); len(vol) > 0 {
			sum, err := strconv.ParseInt(vol[1], 10, 64)
			handleError(err)
			s, err := strconv.ParseInt(vol[2], 10, 64)
			handleError(err)
			r, err := strconv.ParseInt(vol[3], 10, 64)
			handleError(err)
			stats = append(stats, usage{sum, s, r})
		}
	}

	if len(stats) != 5 {
		err := fmt.Errorf("Missing stats. Expected: 5 - Detected: %d", len(stats))
		handleError(err)
	}

	for i, stat := range stats {
		fmt.Printf("%15s: Sum: %10d MB - Sent: %10d MB - Received: %10d MB\n", statRange[i], stat.summb, stat.sentmb, stat.recvmb)
	}
}
