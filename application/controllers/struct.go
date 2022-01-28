package controllers

//外层结构体
type cfg struct {
	Db  mysql `toml:"mysql"`
	Kfk kafka `toml:"kafka"`
	Es  es    `toml:"es"`
}

//内层结构体
type mysql struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Endpoint string `toml:"endpoint"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
	Charset  string `toml:"charset"`
}

//内层结构体
type kafka struct {
	Endpoint string `toml:"endpoint"`
	Topic    string `toml:"topic"`
	GroupID  string `toml:"group"`
}

//内层结构体
type es struct {
	Endpoint string `toml:"endpoint"`
}
