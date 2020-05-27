// Author       kevin
// Time         2019-09-20 18:15
// File Desc    jwt test

package jwt

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var (
	privateKey *ecdsa.PrivateKey
)

// 初始化，读取JWT签名的公私钥
func init() {
	var err error
	// private key
	privateKeyBytes, err := ioutil.ReadFile("/home/kevin/go/key.pem")
	privateKey, err = ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}
}

// 签发token
func ExampleIssueToken() {
	// 生成token
	builder := NewTokenBuilder(100, "kevin", privateKey)
	builder.IssueAt(time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local))
	builder.NotBefore(time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local))
	builder.ExpiresAt(time.Date(2100, 2, 1, 0, 0, 0, 0, time.Local))
	token, err := builder.Build()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(token)
	// Output:
}

// 解析token
func ExampleParseToken() {

	var err error

	// 生成token
	builder := NewTokenBuilder(100, "kevin", privateKey)
	builder.IssueAt(time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local))
	builder.NotBefore(time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local))
	builder.ExpiresAt(time.Date(2100, 2, 1, 0, 0, 0, 0, time.Local))
	token, err := builder.Build()
	if err != nil {
		log.Println(err)
	}

	parser := new(Parser)
	parser.ValidMethods = []string{"ES256"}
	jwtToken, err := parser.Parse(token, &privateKey.PublicKey)
	if err != nil {
		log.Println(err)
	}
	jwtToken.PrintExpireTime()
	// Output:
	// exp: 2100-02-01 00:00:00 +0800 CST
}

// 续签token
func ExampleRenewToken() {

	// 生成一个1小时后过期的token
	// 生成token
	builder := NewTokenBuilder(100, "kevin", privateKey)
	builder.IssueAt(time.Now())
	builder.NotBefore(time.Now())
	builder.ExpiresAt(time.Now().Add(time.Hour))
	token, err := builder.Build()
	if err != nil {
		log.Println(err)
	}
	// 解析token
	parser := new(Parser)
	parser.ValidMethods = []string{"ES256"}
	jwtToken, err := parser.Parse(token, &privateKey.PublicKey)

	if err != nil {
		log.Println(err)
	}
	jwtToken.PrintExpireTime()
	jwtToken.Renew(privateKey)
	jwtToken.PrintExpireTime()
	// Output:
	//
}
