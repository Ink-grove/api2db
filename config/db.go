package config

type DBConfig struct {
	DbType      string `json:"db_type"`
	LogMode     string `json:"log_mode"`
	DSN         string `json:"dsn"`          // data source name.
	Active      int    `json:"active"`       // pool
	Idle        int    `json:"idle"`         // pool
	IdleTimeout string `json:"idle_timeout"` // connect max life time.
}
