package util

import (
	"github.com/whj1990/go-core/os"
	"github.com/cloudwego/kitex/client"
	dns "github.com/kitex-contrib/resolver-dns"
)

func GetClientOption(port string) client.Option {
	if os.RunningInDocker() {
		return client.WithResolver(dns.NewDNSResolver())
	}
	return client.WithHostPorts("0.0.0.0:" + port)
}
