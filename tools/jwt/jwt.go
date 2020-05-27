// Author       kevin
// Time         2019-09-27 13:31
// File Desc    todo

package jwt

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"yann-chat/tools/log"
)

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// service uses a different time zone than your tokens.
var TimeFunc = time.Now

// Parse methods use this callback function to supply
// the key for verification.  The function receives the parsed,
// but unverified Token.  This allows you to use properties in the
// Header of the token (such as `kid`) to identify which key to use.
type Keyfunc func(*Token) (interface{}, error)

// A JWT Token.  Different fields will be used depending on whether you're
// creating or parsing/verifying a token.
type Token struct {
	Raw       string                 // The raw token.  Populated when you Parse a token
	Method    SigningMethod          // The signing method used or to be used
	Header    map[string]interface{} // The first segment of the token
	Claims    Claims                 // The second segment of the token
	Signature string                 // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}

// 创建一个新的 JWT
// [param]
// method: JWT 签名方法
// claims: JWT 内容
// [return]
func NewWithClaims(method SigningMethod, claims Claims) *Token {
	return &Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}

// 续签token
// [param]
// renewTime: token过期续约时间，比如renewTime是5天， t的过期时间是三天后，那么就签发新的token
// [return]
// flag: true表示更新过token; false表示没有更新过token
func (t *Token) Renew(privateKey *ecdsa.PrivateKey) (flag bool) {
	// token解析器
	parser := new(Parser)
	parser.ValidMethods = []string{t.Method.Alg()}
	jwtToken, err := parser.Parse(t.Raw, &privateKey.PublicKey)
	//fmt.Println("过期时间:", jwtToken.Claims.ExpiresAt)
	if err != nil {
		log.Error(err.Error())
	}
	if jwtToken == nil {
		log.Error("jwt is nil")
	} else {
		if jwtToken.Claims.NotBefore < time.Now().Unix() && jwtToken.Claims.ExpiresAt-time.Now().Unix() < int64(
			RenewTime/1e9) {
			t.Claims.IssuedAt = time.Now().Unix()
			t.Claims.ExpiresAt = time.Now().Add(DefaultExpireTime).Unix()
			t.Claims.NotBefore = t.Claims.IssuedAt
			return true
		}
	}
	return false
}

// 为JWT签名
func (t *Token) SignedString(key interface{}) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = t.SigningString(); err != nil {
		return "", err
	}
	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", err
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

// Generate the signing string.  This is the
// most expensive part of the whole deal.  Unless you
// need this for something special, just go straight for
// the SignedString.
func (t *Token) SigningString() (string, error) {
	var err error
	parts := make([]string, 2)
	for i := range parts {
		var jsonValue []byte
		if i == 0 {
			if jsonValue, err = json.Marshal(t.Header); err != nil {
				return "", err
			}
		} else {
			if jsonValue, err = json.Marshal(t.Claims); err != nil {
				return "", err
			}
		}

		parts[i] = EncodeSegment(jsonValue)
	}
	return strings.Join(parts, "."), nil
}

// 打印过期时间
func (t *Token) PrintExpireTime() {
	timeString := fmt.Sprintf("%d", t.Claims.ExpiresAt)
	i, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		panic(err)
	}
	expTime := time.Unix(i, 0)
	fmt.Printf("exp: %s\n", expTime)
}

// Parse, validate, and return a token
func Parse(tokenString string, publicKey *ecdsa.PublicKey) (*Token, error) {
	return new(Parser).Parse(tokenString, publicKey)
}

// Encode JWT specific base64url encoding with padding stripped
func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

// Decode JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(seg)
}

// 校验 token 并更新
// [return]
// flag: true, 更新过token, false, 没有更新过token
// err: 如果不为空, 表示校验token失败
// 调用方应该先判断err是否为空, 如果不为空, 再判断flag是否为true, 如果为true, 表示token更新过.
func VerifyAndRenewToken(tokenStr string, privateKey *ecdsa.PrivateKey) (tokenObj *Token, err error, flag bool) {
	parser := new(Parser)
	parser.ValidMethods = []string{"ES256"}
	tokenObj, err = parser.Parse(tokenStr, &privateKey.PublicKey)
	if err != nil {
		return tokenObj, err, false
	}
	return tokenObj, nil, flag
}
