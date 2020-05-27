package jwt

import (
	"crypto/ecdsa"
	"github.com/spf13/cast"
	"time"
)

// token 类型： 登录用的，修改密码用的等等
const (
	Login = iota + 1 //登录签发
)

// 时间常量
const (
	Day               = 24 * time.Hour
	Month             = 720 * time.Hour
	RenewTime         = Day * 5
	DefaultExpireTime = Month
)

type TokenBuilder struct {
	// jwt payload
	claims Claims
	// 签名方法
	signMethod SigningMethod
	// 私钥
	privateKey *ecdsa.PrivateKey
}

// 创建 token builder,
// 默认参数:
// 签名算法: ES256
// 过期时间: 一个月
// 生效时间: 当前时间
// 签发者: Taurus Future
// [params]
// aud: 签发对象
// publicKey, privateKey: 公私钥
// [return]
func NewTokenBuilder(userId int64, deviceUuid string, privateKey *ecdsa.PrivateKey) *TokenBuilder {
	return new(TokenBuilder).
		Audience(cast.ToString(userId)).
		UserID(userId).
		PrivateKey(privateKey).
		SignMethod(GetSigningMethod("ES256")).
		ExpiresAt(time.Now().Add(DefaultExpireTime)).
		Issuer("Chicha").
		IssueAt(time.Now()).DeviceUuid(deviceUuid)
}

// 设置签发对象
func (builder *TokenBuilder) Audience(aud string) *TokenBuilder {
	builder.claims.Audience = aud
	return builder
}

// 设置签发机构
func (builder *TokenBuilder) Issuer(iss string) *TokenBuilder {
	builder.claims.Issuer = iss
	return builder
}

// 设置签发时间
func (builder *TokenBuilder) IssueAt(iat time.Time) *TokenBuilder {
	builder.claims.IssuedAt = iat.Unix()
	return builder
}

// 设置Token ID
func (builder *TokenBuilder) UserID(jti int64) *TokenBuilder {
	builder.claims.UserId = jti
	return builder
}

// 设置备注
func (builder *TokenBuilder) Subject(sub string) *TokenBuilder {
	builder.claims.Subject = sub
	return builder
}

// 用户设备id
func (builder *TokenBuilder) DeviceUuid(dev string) *TokenBuilder {
	builder.claims.DeviceUuid = dev
	return builder
}

// 设置过期时间
func (builder *TokenBuilder) ExpiresAt(exp time.Time) *TokenBuilder {
	builder.claims.ExpiresAt = exp.Unix()
	return builder
}

// 设置生效时间
func (builder *TokenBuilder) NotBefore(nbf time.Time) *TokenBuilder {
	builder.claims.NotBefore = nbf.Unix()
	return builder
}

// 设置私钥
func (builder *TokenBuilder) PrivateKey(privateKey *ecdsa.PrivateKey) *TokenBuilder {
	builder.privateKey = privateKey
	return builder
}

// 设置token签名算法
func (builder *TokenBuilder) SignMethod(method SigningMethod) *TokenBuilder {
	builder.signMethod = method
	return builder
}

// 构建token
func (builder *TokenBuilder) Build() (token string, err error) {
	token, err = NewWithClaims(builder.signMethod, builder.claims).SignedString(builder.privateKey)
	return
}
