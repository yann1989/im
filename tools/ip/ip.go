// Author       kevin
// Time         2019-08-08 20:16
// File Desc    ip 相关的工具类, 比如获取当前主机的网卡ip等

package ip

import (
	"net"
	"strings"
)

// ExternalIP 获取系统绑定的所有外网IP
func ExternalIP() (res []string) {

	// 获取所有网卡
	inters, err := net.Interfaces()
	if err != nil {
		return
	}

	// 遍历网卡, 查找所有的外网IP, 如果没有返回""
	for _, inter := range inters {

		// 排除loopback interface
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ipnet.IP.IsLoopback() || ipnet.IP.IsLinkLocalMulticast() || ipnet.IP.IsLinkLocalUnicast() {
						continue
					}
					// 剔除一些特殊的ip, 比如, 10开头的A类IP
					if ip4 := ipnet.IP.To4(); ip4 != nil {
						switch true {
						case ip4[0] == 10:
							continue
						case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
							continue
						case ip4[0] == 192 && ip4[1] == 168:
							continue
						default:
							res = append(res, ipnet.IP.String())
						}
					}
				}
			}
		}
	}
	return
}

// 获取系统绑定的内网IPv4地址
func InternalIP() string {

	// 获取所有网卡
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}

	// 遍历网卡, 查找IPv4地址, 如果没有返回""
	for _, inter := range inters {
		// 排除loopback interface
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}
