package config

type Config struct {
	ServerPort      string `json:"port" yaml:"port"`
	JWTSecret       string `json:jwt_secret" yaml:"jwt_secret"`
	UserPostgres    `json:"user_postgres" yaml:"user_postgres"`
	BankingPostgres `json:"banking_postgres" yaml:"banking_postgres"`
}

type UserPostgres struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

type BankingPostgres struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

func InitConfig() (*Config, error) {
	return &Config{}, nil
}
