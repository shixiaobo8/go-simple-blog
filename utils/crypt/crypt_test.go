package crypt

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"
)

func TestDesCbcEnc(t *testing.T)  {
	data := []byte("123456")
	key := []byte("abcdefgh")
	res, err := DesCbcEnc(data, key, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("hex:", hex.EncodeToString(res))
	encode := base64.StdEncoding.EncodeToString(res)
	t.Log("base64:", encode)
}

func TestDesCbcDec(t *testing.T) {
	b64Str := "opfVWWpnmcY="
	data, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		log.Fatal(err)
	}
	key := []byte("abcdefgh")
	res, err := DesCbcDec(data, key, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("origin:", string(res))
}