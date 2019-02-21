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
)

// URL Constants -
const (
	CountrCsvFileURL = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country-CSV.zip"
	CountryInfoURL   = "http://download.geonames.org/export/dump/countryInfo.txt"
)

// Country -
type Country struct {
	Code string
	Name string
}

// Countries -
type Countries map[string]Country

func main() {

	csvFile, err := DownloadFile(CountrCsvFileURL)
	HandleError(err)

	files, err := Unzip(csvFile, ".")
	HandleError(err)

	countryFileName, err := DownloadFile(CountryInfoURL)
	HandleError(err)

	InitialCountries := Countries{
		"6255146": Country{"AF", "Africa"},
		"6255147": Country{"AS", "Asia"},
		"6255148": Country{"EU", "Europe"},
		"6255149": Country{"NA", "North America"},
		"6255150": Country{"SA", "South America"},
		"6255151": Country{"OC", "Oceania"},
		"6255152": Country{"AN", "Antarctica"},
	}
	countries := ParseCountries(countryFileName)
	for k, v := range InitialCountries {
		countries[k] = v
	}

	for _, c := range countries {
		fmt.Printf("%v\n", c)
	}

	_ = files
}

// HandleError -
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

// DownloadFile -
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

// Unzip -
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

// ParseCountries -
func ParseCountries(fileName string) Countries {

	fmt.Printf("Parsing %s\n", fileName)
	countries := make(map[string]Country)

	file, err := os.Open(fileName)
	HandleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
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
