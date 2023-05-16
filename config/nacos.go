package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"github.com/whj1990/go-core/handler"
	"gopkg.in/yaml.v3"
)

// nacos 配置客户端
var naCosConfigClient config_client.IConfigClient
var configData ConfigData

// 服务发现客户端
var naCosNamingClient naming_client.INamingClient

func NaCosInitConfigClient() {
	sc, cc := GetNaCosBaseConfig()
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
	//获取配置
	content, err := naCosConfigClient.GetConfig(vo.ConfigParam{
		DataId: viper.GetString("naCos.data-id"),
		Group:  viper.GetString("naCos.group"),
	})
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal([]byte(content), &configData)
	//配置监听
	err = naCosConfigClient.ListenConfig(vo.ConfigParam{
		DataId: viper.GetString("naCos.data-id"),
		Group:  viper.GetString("naCos.group"),
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println(data)
			yaml.Unmarshal([]byte(data), &configData)
		},
	})
	if err != nil {
		handler.HandleError(err)
	}
}
func GetNacosConfigData() *ConfigData {
	return &configData
}
func GetNaCosBaseConfig() ([]constant.ServerConfig, constant.ClientConfig) {
	if configPath == "" {
		path := flag.String("c", "conf", "config path, eg: -c conf")
		flag.Parse()
		configPath = *path
	}
	viper.AddConfigPath(configPath)
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
		constant.WithNamespaceId(viper.GetString("naCos.namespace")),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/naCos/log"),
		constant.WithCacheDir("/tmp/naCos/cache"),
		constant.WithLogLevel(viper.GetString("naCos.log-level")),
	)
	return serverConfigs, clientConfig
}

func NewNaCosNamingClient(isServer bool) {
	serverConfigs, clientConfig := GetNaCosBaseConfig()
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
	if isServer {
		RegisterServiceInstance(vo.RegisterInstanceParam{
			Ip:          configData.GrpcServer.Ip,
			Port:        uint64(configData.GrpcServer.Port),
			ServiceName: configData.GrpcServer.Name,
			GroupName:   viper.GetString("naCos.group"),
			ClusterName: viper.GetString("naCos.cluster-name"),
			Weight:      configData.GrpcServer.Weight,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    map[string]string{"idc": "shanghai", "timestamp": time.Now().Format(time.DateTime)},
		})
	}

}
func RegisterServiceInstance(param vo.RegisterInstanceParam) {

	success, err := naCosNamingClient.RegisterInstance(param)
	if !success || err != nil {
		panic("RegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}
func GetService(param vo.GetServiceParam) model.Service {
	service, err := naCosNamingClient.GetService(param)
	if err != nil {
		panic("GetService failed!" + err.Error())
	}
	return service
}
func Subscribe(param *vo.SubscribeParam) {
	naCosNamingClient.Subscribe(param)
}

func UnSubscribe(param *vo.SubscribeParam) {
	naCosNamingClient.Unsubscribe(param)
}
func SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) *model.Instance {
	instances, err := naCosNamingClient.SelectOneHealthyInstance(param)
	if err != nil {
		panic("SelectOneHealthyInstance failed!")
	}
	return instances
}
