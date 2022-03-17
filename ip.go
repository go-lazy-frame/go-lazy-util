//             ,%%%%%%%%,
//           ,%%/\%%%%/\%%
//          ,%%%\c "" J/%%%
// %.       %%%%/ o  o \%%%
// `%%.     %%%%    _  |%%%
//  `%%     `%%%%(__Y__)%%'
//  //       ;%%%%`\-/%%%'
// ((       /  `%%%%%%%'
//  \\    .'          |
//   \\  /       \  | |
//    \\/攻城狮保佑) | |
//     \         /_ | |__
//     (___________)))))))                   `\/'
/*
 * 修订记录:
 * long.qian 2021-10-04 18:32 创建
 */

/**
 * @author long.qian
 */

package go_lazy_util

import (
	"fmt"
	"github.com/go-ping/ping"
	"net"
	"runtime"
	"strings"
	"time"
)

var (
	IpUtil = new(ipUtil)
)

type ipUtil struct {
}

// GetLocalIp 获取本地IP
func (me *ipUtil) GetLocalIp() string {
	ip := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				if strings.HasPrefix(ip, "192.168") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.") {
					break
				}
			}

		}
	}
	if ip == "" {
		return "127.0.0.1"
	}
	return ip
}

func (me *ipUtil) IsPing(host string, timeout time.Duration) (bool, error) {
	if timeout <= 0 {
		return false, fmt.Errorf("%s", "ping timeout 必须大于 0")
	}
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return false, err
	}
	pinger.Count = 1
	pinger.Timeout = timeout
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}
	// Linux 上，注意查看：https://github.com/go-ping/ping/blob/master/README.md#supported-operating-systems
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		if runtime.GOOS == "linux" {
			return false, fmt.Errorf("%s %v", `Linux 平台，您可能需要执行以下命令：sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"`, err)
		}
		return false, err
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv >= 1 {
		return true, nil
	}
	return false, nil
}
