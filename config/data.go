package config

type ConfigData struct {
	HttpServer `yaml:"HttpServer"`
	GrpcServer `yaml:"grpcServer"`
	Jaeger
	Db
	Redis
	Influxdb
	Kafka
	Mongo
	OSS
}
type HttpServer struct {
	Port string
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