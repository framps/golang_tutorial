package main

// Inspired by https://github.com/mschmitt/GeoLite2xtables

// Work in progress

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
)

const (
	COUNTRY_CSV_FILE = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country-CSV.zip"
	COUNTRY_INFO_URL = "http://download.geonames.org/export/dump/countryInfo.txt"
)

type Country struct {
	Code string
	Name string
}

type Countries map[string]Country

func main() {

	csvFile, err := DownloadFile(COUNTRY_CSV_FILE)
	HandleError(err)

	files, err := Unzip(csvFile, ".")
	HandleError(err)

	countryFileName, err := DownloadFile(COUNTRY_INFO_URL)
	HandleError(err)

	InitialCountries := Countries{
		"6255146", Country{"Africa", "AF"},
		"6255147": Country{"Asia", "AS"},
		"6255148": Country{"Europe", "EU"},
		"6255149", Country{"North America", "NA"},
		"6255150", Country{"South America", "SA"},
		"6255151", Country{"Oceania", "OC"},
		"6255152": Country{"Antarctica", "AN"},
	}
	countries := ParseCountries(countryFileName)
	copier.Copy(&countries, &InitialCountries)

	for _, c := range countries {
		fmt.Printf("%v\n", c)
	}

	_ = files
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func DownloadFile(url string) (string, error) {
	fmt.Println("Retrieving", url)
	resp, err := http.Get(url)
	HandleError(err)
	defer resp.Body.Close()

	_, filename := path.Split(url)
	HandleError(err)
	fmt.Println("Creating", filename)
	out, err := os.Create(filename)
	HandleError(err)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	HandleError(err)
	return filename, err
}

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	fmt.Printf("Unzipping %s into %s\n", src, dest)
	r, err := zip.OpenReader(src)
	HandleError(err)
	defer r.Close()

	for _, f := range r.File {
		fmt.Printf("Unzipping %s\n", f.Name)
		rc, err := f.Open()
		HandleError(err)
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				HandleError(err)
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			HandleError(err)
			_, err = io.Copy(outFile, rc)
			HandleError(err)
			outFile.Close()
		}
	}
	return filenames, nil
}

func ParseCountries(fileName string) Countries {

	fmt.Printf("Parsing %s\n", fileName)
	countries := make(map[string]Country)

	file, err := os.Open(fileName)
	HandleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasSuffix(l, "#") {
			continue
		}
		tokens := strings.SplitN(scanner.Text(), "\t", -1)
		id := tokens[16]
		c := Country{Code: tokens[0], Name: tokens[4]}
		countries[id] = c
	}

	if err := scanner.Err(); err != nil {
		HandleError(err)
	}
	return countries
}
