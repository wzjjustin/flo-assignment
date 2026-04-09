package service

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DBname   string `yaml:"db_name"`
	Port     int    `yaml:"port"`
}

type Config struct {
	DBConfig   DBConfig `yaml:"db"`
	NumWorkers int      `yaml:"num_workers"`
}
