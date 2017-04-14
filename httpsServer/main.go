package main

// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// See github.com/framps/golang_tutorial for latest code
//
// Samples for go http support - simple hello world which redirects unsecure port
// http://localhost:8080 to https://localhost:8443. In addition a TLS certificate is
// generated if it doesn't exist already

// Part of code from https://github.com/denji/golang-tls
// Part of code from https://github.com/golang/go/blob/master/src/crypto/tls/generate_cert.go

// Open http://localhost:8080 which will be redirected to https://localhost:8443

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	httpPort     = ":8080"
	httpsPort    = ":8443"
	endpoint     = "/"
	certFileName = "cert.pem"
	keyFileName  = "key.pem"
)

// create TLS certificate

func createCertificate() {

	fmt.Printf("Creating certificate...")
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serNum, _ := rand.Int(rand.Reader, max)

	sub := pkix.Name{
		Organization:       []string{"Sample organization"},
		OrganizationalUnit: []string{"Sample organizational unit"},
		CommonName:         "Sample common name",
	}

	template := x509.Certificate{
		SerialNumber: serNum,
		Subject:      sub,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // one year
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	reportWhenFailure(err)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	reportWhenFailure(err)

	certOut, err := os.Create("cert.pem")
	reportWhenFailure(err)

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.Create("key.pem")
	reportWhenFailure(err)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()

	fmt.Printf(" done\n")

}

// helper for error handling
func reportWhenFailure(err error) {
	if err != nil {
		panic(fmt.Errorf("Caught error %v\n", err))
	}
}

// HelloServer -
func HelloServer(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello world from a TLS secured server.\n"))
}

// redirector from http port to https port
func redirect(w http.ResponseWriter, req *http.Request) {
	hostURLParts := strings.Split(req.Host, ":")
	redirectURL := "https://" + hostURLParts[0] + httpsPort + req.URL.String()
	fmt.Printf("Redirecting to %s\n", redirectURL)
	func() {
		http.Redirect(w, req,
			redirectURL,
			http.StatusMovedPermanently)
	}()
}

func main() {

	if _, err := os.Stat(certFileName); os.IsNotExist(err) {
		createCertificate()
	} else if _, err := os.Stat(keyFileName); os.IsNotExist(err) {
		createCertificate()
	}

	// redirect every http request to https
	go http.ListenAndServe(httpPort, http.HandlerFunc(redirect))

	fmt.Printf("Listening on https port %s and endpoint %s\n", httpsPort, endpoint)

	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, HelloServer)

	err := http.ListenAndServeTLS(httpsPort, certFileName, keyFileName, mux)
	reportWhenFailure(err)
}
