// Author       kevin
// Time         2019-09-20 08:48
// File Desc    web 层配置类

package web

// Web层配置
type Config struct {
	Proto    string // protocol, default to https
	Address  string // "IP:PORT"
	CertPath string // 证书路径
	KeyPath  string // 秘钥路径
}
