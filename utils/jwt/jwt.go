package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go-blog/utils"
	"hash"
	"strings"
	"time"
)

/**
https://jwt.io/
 */
type Header struct {
	AlgInstance hash.Hash	`json:"-"`
	Alg string
	Typ string
}

type Payload struct {
	Iss string		// jwt签发者
	Iat time.Time	// 签发时间
	Exp time.Time   // 过期时间
	Nbf time.Time   // token信息生效时间.这个值可以不设置,但是设定后,一定要大于当前Unix UTC,否则token将会延迟生效.
	Aud string		// jwt接收方
	Sub string		// 面向的用户
	Jti string		// token唯一标识id
}

type Jwt struct {
	Header
	Payload
	Extra map[string]interface{}
}

type ExpireError struct {
	error
	Exp time.Time
}

func (e ExpireError) Error() string {
	return "token expired"
}

// 如token失效，则当前时间若 - 失效时间 < refreshSecond 重新刷新token
const RefreshSecond = 30 * time.Minute

func (j *Jwt) Token() (string, error) {
	var err error
	headerJson, err := json.Marshal(j.Header)
	if err != nil {
		return "", err
	}
	// 需要将Extra加入json数据
	var m = j.getPayloadMap()
	payloadJson, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s.%s", base64.RawURLEncoding.EncodeToString(headerJson), base64.RawURLEncoding.EncodeToString(payloadJson), j.Signature(headerJson, payloadJson)), nil
}

// Signature 部分是对前两部分的签名，防止数据篡改
func (j *Jwt) Signature(headerJson []byte, payloadJson []byte) string {
	c := make([]byte, 0, len(headerJson) + len(payloadJson))
	c = append(append(c, headerJson...),  payloadJson...)
	j.AlgInstance.Reset()
	j.AlgInstance.Write(c)
	return base64.RawURLEncoding.EncodeToString(j.AlgInstance.Sum(nil))
}

func (j *Jwt) VerifyToken(token string) (bool, error) {
	var err error
	arr := strings.SplitN(token, ".", 3)
	if len(arr) != 3 {
		return false, errors.New("token is not available")
	}
	var headerJson []byte
	var payloadJson []byte

	headerJson, err = base64.RawURLEncoding.DecodeString(arr[0])
	if err != nil {
		return false, err
	}

	payloadJson, err = base64.RawURLEncoding.DecodeString(arr[1])
	if err != nil {
		return false, err
	}
	if j.Signature(headerJson, payloadJson) != arr[2] {
		return false, errors.New("verify error")
	}
	var payload Payload
	err = json.Unmarshal(payloadJson, &payload)

	if err != nil {
		return false, err
	}

	var zeroTime time.Time
	// 签发时间大于当前服务器时间
	if !zeroTime.Equal(payload.Iat) && payload.Iat.After(time.Now()) {
		return false, errors.New("issue time greater than current server time")
	}

	// 过期时间小于当前服务器时间
	if !zeroTime.Equal(payload.Exp) && payload.Exp.Before(time.Now()) {
		return false, ExpireError{
			Exp: payload.Exp,
		}
	}

	if !zeroTime.Equal(payload.Nbf) && payload.Nbf.After(time.Now()) {
		return false, errors.New("token has yet to take effect")
	}
	return true, nil
}

// 追加其他信息
func (j *Jwt) setAdditional(key string, value interface{}) {
	j.Extra[key] = value
}

func (j *Jwt) getPayloadMap() map[string]interface{} {
	m := utils.StructToMap(j.Payload)
	for k, v := range j.Extra {
		m[k] = v
	}
	return m
}

func NewJwt(header Header, payload Payload, extra map[string]interface{}) *Jwt {
	return &Jwt{
		header,
		payload,
		extra,
	}
}