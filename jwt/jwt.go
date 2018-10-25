package main

// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// See github.com/framps/golang_tutorial for latest code

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lestrrat/go-jwx/jwk"
	"github.com/lestrrat/go-jwx/jwt"
)

func ExampleJWT() *jwt.Token {

	const aLongLongTimeAgo = 233431200

	t := jwt.New()
	t.Set(jwt.SubjectKey, `https://github.com/lestrrat/go-jwx/jwt`)
	t.Set(jwt.AudienceKey, `Golang Users`)
	t.Set(jwt.IssuedAtKey, time.Unix(aLongLongTimeAgo, 0))
	t.Set(`privateClaimKey`, `Hello, World!`)

	buf, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		fmt.Printf("failed to generate JSON: %s\n", err)
		return nil
	}

	fmt.Printf("%s\n", buf)
	fmt.Printf("aud -> '%s'\n", t.Audience())
	fmt.Printf("iat -> '%s'\n", t.IssuedAt().Format(time.RFC3339))
	if v, ok := t.Get(`privateClaimKey`); ok {
		fmt.Printf("privateClaimKey -> '%s'\n", v)
	}
	fmt.Printf("sub -> '%s'\n", t.Subject())

	return t
}

func ExampleKey() string {

	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return ""
	}

	key, err := jwk.New(&privkey.PublicKey)
	if err != nil {
		log.Printf("failed to create JWK: %s", err)
		return ""
	}

	jsonbuf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		log.Printf("failed to generate JSON: %s", err)
		return ""
	}

	return string(jsonbuf)
}

func ExampleCheck(token *jwt.Token, key string) {

	set, err := jwk.Fetch("https://foobar.domain/jwk.json")
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return
	}

	// If you KNOW you have exactly one key, you can just
	// use set.Keys[0]
	keys := set.LookupKeyID("mykey")
	if len(keys) == 0 {
		log.Printf("failed to lookup key: %s", err)
		return
	}

	key, err := key.Materialize()
	if err != nil {
		log.Printf("failed to create public key: %s", err)
		return
	}

	// Use key for jws.Verify() or whatever
}

func main() {

	token := ExampleJWT()
	key := ExampleKey()

	token, err := jwt.Parse(token, getKey)
	if err != nil {
		panic(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range claims {
		fmt.Printf("%s\t%v\n", key, value)
	}
}
