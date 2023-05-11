package config

import (
	"flag"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"github.com/whj1990/go-core/handler"
	"strconv"
	"time"
)

// nacos 配置客户端
var naCosConfigClient config_client.IConfigClient

func NaCosInitConfigClient() {
	path := flag.String("c", "conf", "config path, eg: -c conf")
	flag.Parse()
	configPath = *path

	viper.AddConfigPath(*path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			viper.GetString("naCos.address"),
			viper.GetUint64("naCos.port"),
			constant.WithContextPath(viper.GetString("naCos.context-path")),
		),
	}
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(viper.GetString("namespace")),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/naCos/log"),
		constant.WithCacheDir("/tmp/naCos/cache"),
		constant.WithLogLevel(viper.GetString("naCos.log-level")),
	)
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	naCosConfigClient = client
}
func getConfig(params vo.ConfigParam) string {
	content, _ := naCosConfigClient.GetConfig(params)
	return content
}
func GetNaCosString(dataId, defaultValue string) string {
	content := getConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  viper.GetString("service-group"),
	})
	if content == "" {
		content = defaultValue
	}
	return content
}

func GetNaCosInt(dataId string, defaultValue int) int {
	content := getConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  viper.GetString("service-group"),
	})
	if content == "" {
		return defaultValue
	}
	res, err := strconv.Atoi(content)
	if err != nil {
		handler.HandleError(err)
	}
	return res
}

func GetNaCosBool(dataId string, defaultValue bool) bool {
	content := getConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  viper.GetString("service-group"),
	})
	if content == "" {
		return defaultValue
	}
	res, err := strconv.ParseBool(content)
	if err != nil {
		handler.HandleError(err)
	}
	return res
}

// 服务发现客户端
var naCosNamingClient naming_client.INamingClient

func NewNaCosNamingClient() {
	path := flag.String("c", "conf", "config path, eg: -c conf")
	flag.Parse()
	configPath = *path
	viper.AddConfigPath(*path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			viper.GetString("naCos.address"),
			viper.GetUint64("naCos.port"),
			constant.WithContextPath(viper.GetString("naCos.context-path")),
		),
	}
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(viper.GetString("namespace")),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/naCos/log"),
		constant.WithCacheDir("/tmp/naCos/cache"),
		constant.WithLogLevel(viper.GetString("naCos.log-level")),
	)
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	naCosNamingClient = client
	RegisterServiceInstance(vo.RegisterInstanceParam{
		Ip:          GetNaCosString("service.ip", ""),
		Port:        uint64(GetNaCosInt("server.port", 8310)),
		ServiceName: GetNaCosString("server.name", "mine.grpc"),
		GroupName:   viper.GetString("service-group"),
		ClusterName: viper.GetString("cluster-name"),
		Weight:      float64(GetNaCosInt("server.weight", 10)),
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai", "timestamp": time.Now().Format(time.DateTime)},
	})
}
func RegisterServiceInstance(param vo.RegisterInstanceParam) {

	success, err := naCosNamingClient.RegisterInstance(param)
	if !success || err != nil {
		panic("RegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}
func GetService(param vo.GetServiceParam) {
	service, err := naCosNamingClient.GetService(param)
	if err != nil {
		panic("GetService failed!" + err.Error())
	}
	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)
}
func Subscribe(param *vo.SubscribeParam) {
	naCosNamingClient.Subscribe(param)
}

func UnSubscribe(param *vo.SubscribeParam) {
	naCosNamingClient.Unsubscribe(param)
}
