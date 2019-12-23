package jwt

import (
	"crypto/sha256"
	"log"
	"testing"
	"time"
)

func TestNewJwt(t *testing.T) {
	header := Header{
		AlgInstance: sha256.New(),
		Alg: "sha256",
		Typ: "jwt",
	}
	var total time.Duration = 10
	var interval time.Duration = 2

	payload := Payload{
		Iss: "test",
		Iat: time.Now(),
		Exp: time.Now().Add(total*time.Second),
		Sub: "www.baidu.com",
	}

	jwt := NewJwt(header, payload, map[string]interface{}{
		"adminId": 1,
	})

	token, err := jwt.Token()
	if err != nil {
		log.Fatal(err)
	}
	jwt2 := NewJwt(Header{
		AlgInstance: sha256.New(),
		Alg: "sha256",
		Typ: "jwt",
	}, Payload{}, map[string]interface{}{
			"adminId": 1,
	})

	var count = 0
	var ticker = time.NewTicker(interval*time.Second)
	for range ticker.C {
		count++
		r, err := jwt2.VerifyToken(token)
		if err != nil && count != (int(total)/int(interval)) {
			t.Fatal(err)
		}
		if count == (int(total)/int(interval)) {
			t.Log(r)
			ticker.Stop()
			break
		} else {
			t.Log("true", count, "Time:", time.Now().Unix())
		}
	}
}