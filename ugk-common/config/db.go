package config

type MongoConfig struct {
	Url      string `json:"url"`
	Database string `json:"database"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int32  `json:"port"`
}
