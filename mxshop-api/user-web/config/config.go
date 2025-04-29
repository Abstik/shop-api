package config

// 服务总配置
type ServerConfig struct {
	// 以下是此服务的配置（将此服务注册到consul时使用）
	Name string   `mapstructure:"name" json:"name"`
	Host string   `mapstructure:"host" json:"host"`
	Tags []string `mapstructure:"tags" json:"tags"`
	Port int      `mapstructure:"port" json:"port"`

	// 以下是从nacos中获取的配置
	// 用户srv服务的配置，需要用户srv服务的名称，用名称去consul中连接用户srv服务，以便在用户web服务中调用用户srv服务
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo  AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

// 用户服务
type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"` // 用户服务的名称
}

// jwt配置
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"` // jwt签名密钥
}

// 阿里云配置
type AliSmsConfig struct {
	ApiKey     string `mapstructure:"key" json:"key"`
	ApiSecrect string `mapstructure:"secrect" json:"secrect"`
}

// consul配置
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"` // consul服务器的地址
	Port int    `mapstructure:"port" json:"port"` // consul服务器的端口
}

// redis配置
type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"` // redis中的过期时间
}

// nacos配置
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
