package config

type ConfigData struct {
	ClientServer `yaml:"clientServer"` //客户端配置
	GrpcServer   `yaml:"grpcServer"`   //微服务服务端配置
	Jaeger
	Db
	Redis
	Influxdb
	Kafka
	Mongo
	OSS
}
type ClientServer struct {
	Name            string
	GrpcServerName  string `yaml:"grpcServerName"`
	GrpcGroupName   string `yaml:"grpcGroupName"`
	GrpcClusterName string `yaml:"grpcClusterName"` //以英文逗号分割
	Port            string
}
type GrpcServer struct {
	Name    string
	Ip      string
	Port    int
	Network string
	Address string
	Weight  float64
}
type Jaeger struct {
	ServiceName string `yaml:"serviceName"`
	HostPort    string `yaml:"hostPort"`
}
type Db struct {
	Name                string
	Address             string
	UserName            string `yaml:"userName"`
	Password            string
	MaxIdleConnects     int `yaml:"maxIdleConnects"`
	MaxOpenConnects     int `yaml:"maxOpenConnects"`
	ConnMaxLifetimeHour int `yaml:"connMaxLifetimeHour"`
	Write               DbWrite
	Read                DbRead
}
type DbWrite struct {
	Address  string
	UserName string `yaml:"userName"`
	Password string
}
type DbRead struct {
	Address  string
	UserName string `yaml:"userName"`
	Password string
}
type Redis struct {
	Address  string
	Password string
	Database int
}
type Influxdb struct {
	Url      string
	Username string
	Password string
}
type Kafka struct {
	Addrs   string
	Topics  string //多个主题以英文逗号分割
	GroupId string `yaml:"groupId"`
	Reset   string
}
type Mongo struct {
	Address  string
	Username string //多个主题以英文逗号分割
	Password string
}

type OSS struct {
	EndPoint        string `yaml:"endPoint"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Bucket          string
	Host            string
	RegionId        string `yaml:"regionId"`
	RoleArn         string `yaml:"roleArn"`
}
