package config

type Config struct {
	Filename  string
	Address   string
	Route     string
	Port      string
	Threshold int
}

var ServerConfig = Config{
	Filename:  "timestamps.log",
	Address:   "localhost",
	Route:     "/",
	Port:      "8000",
	Threshold: 60,
}
