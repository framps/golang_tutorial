package main

// Samples used in a small go tutorial
//
// Copyright (C) 2019 framp at linux-tips-and-tricks dot de
//
// Samples for go using http calls, files and line parsing
//
// Inspired by https://github.com/mschmitt/GeoLite2xtables

import (
	"archive/zip"
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
)

// URL Constants -
const (
	CountrCsvFileURL = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country-CSV.zip"
	CountryInfoURL   = "http://download.geonames.org/export/dump/countryInfo.txt"
	IPV4csv          = "GeoLite2-Country-CSV_20190219/GeoLite2-Country-Blocks-IPv4.csv"
	IPV6csv          = "GeoLite2-Country-CSV_20190219/GeoLite2-Country-Blocks-IPv6.csv"
)

// Country -
type Country struct {
	Code string
	Name string
}

// CountryBlock -
type CountryBlock struct {
	Network                     string
	GeonameID                   string
	RegisteredCountryGeonameID  string
	RepresentedCountryGeonameID string
	IsAnonymousProxy            bool
	IsSatelliteProvider         bool
}

// Countries -
type Countries map[string]Country

func main() {

	ipv4 := flag.Bool("4", true, "Use IPV4 addresses")
	ipv6 := flag.Bool("6", false, "Use IPV6 addresses")
	//debug := flag.Bool("d", false, "Write debug info")

	flag.Parse()

	cf, err := ioutil.TempFile("", "cf")
	HandleError(err)
	defer os.Remove(cf.Name())
	err = DownloadFile(CountrCsvFileURL, cf)
	HandleError(err)

	cfd, err := ioutil.TempDir("", "cfd")
	HandleError(err)
	defer os.RemoveAll(cfd)
	_, err = Unzip(cf, cfd)
	HandleError(err)

	cif, err := ioutil.TempFile("", "cif")
	HandleError(err)
	defer os.Remove(cif.Name())
	err = DownloadFile(CountryInfoURL, cif)
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
	countries := ParseCountries(cif)
	for k, v := range InitialCountries {
		countries[k] = v
	}

	if *ipv4 {
		ParseCountryBlocks(cfd+"/"+IPV4csv, countries)
	}
	if *ipv6 {
		ParseCountryBlocks(cfd+"/"+IPV6csv, countries)
	}
}

// HandleError -
func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// DownloadFile -
func DownloadFile(url string, file *os.File) error {
	log.Println("Retrieving", url)
	resp, err := http.Get(url)
	HandleError(err)
	defer resp.Body.Close()

	_, filename := path.Split(url)
	HandleError(err)
	log.Println("Creating", filename)
	_, err = io.Copy(file, resp.Body)
	HandleError(err)
	return err
}

// Unzip -
func Unzip(src *os.File, dest string) ([]string, error) {

	var filenames []string

	log.Printf("Unzipping %s into %s\n", src, dest)
	r, err := zip.OpenReader(src.Name())
	HandleError(err)
	defer r.Close()

	for _, f := range r.File {
		log.Printf("Unzipping %s\n", f.Name)
		rc, err := f.Open()
		HandleError(err)
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
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
func ParseCountries(file *os.File) Countries {

	log.Printf("Parsing %s\n", file.Name())
	countries := make(map[string]Country)

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

// ParseCountryBlocks -
func ParseCountryBlocks(fileName string, countries Countries) {

	log.Printf("Parsing %s\n", fileName)

	file, err := os.Open(fileName)
	HandleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "network") {
			continue
		}

		// network,geoname_id,registered_country_geoname_id,represented_country_geoname_id,is_anonymous_proxy,is_satellite_provider
		// 1.0.0.0/24,2077456,2077456,,0,0

		// => "1.0.0.0","1.0.0.255","16777216","16777471","AU","Australia"

		tokens := strings.SplitN(scanner.Text(), ",", -1)
		c := CountryBlock{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4] != "0", tokens[5] != "0"}

		_, ipnet, err := net.ParseCIDR(c.Network)
		HandleError(err)
		startIP, endIP := cidr.AddressRange(ipnet)

		cc := getCountry(c, countries)
		code := cc.Code
		name := cc.Name

		if isIPv4(startIP.String()) {
			startInt := ip4toint(startIP)
			endInt := ip4toint(endIP)
			fmt.Printf("\"%s\",\"%s\",\"%d\",\"%d\",\"%s\",\"%s\"\n",
				startIP, endIP, startInt, endInt, code, name)
		} else {
			startInt := ip6toint(startIP)
			endInt := ip6toint(endIP)
			fmt.Printf("\"%s\",\"%s\",\"%d\",\"%d\",\"%s\",\"%s\"\n",
				startIP, endIP, startInt, endInt, code, name)
		}
	}

	if err := scanner.Err(); err != nil {
		HandleError(err)
	}
}

func ip4toint(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// ipv6Decimal := IP6toInt(net.ParseIP("FE80::0202:B3FF:FE1E:8329"))
func ip6toint(IPv6Address net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPv6Address.To16())
	return IPv6Int
}

func isIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

func getCountry(cb CountryBlock, countries Countries) Country {

	switch {
	case cb.IsAnonymousProxy:
		return Country{"A1",
			"Anonymous Proxy"}
	case cb.IsSatelliteProvider:
		return Country{"A2",
			"Satellite Provider"}
	case len(cb.RepresentedCountryGeonameID) > 0:
		return Country{countries[cb.RepresentedCountryGeonameID].Code,
			countries[cb.RepresentedCountryGeonameID].Name}
	case len(cb.RegisteredCountryGeonameID) > 0:
		return Country{countries[cb.RegisteredCountryGeonameID].Code,
			countries[cb.RegisteredCountryGeonameID].Name}
	case len(cb.GeonameID) > 0:
		return Country{countries[cb.GeonameID].Code,
			countries[cb.GeonameID].Name}
	default:
		log.Printf("Unknown Geoname ID, panicking. This is a bug.\n")
		log.Printf("ID: %s\n", cb.GeonameID)
		log.Printf("ID Registered: %s\n", cb.RegisteredCountryGeonameID)
		log.Printf("ID Represented %s\n", cb.RepresentedCountryGeonameID)
		os.Exit(1)
		return Country{}
	}
}
