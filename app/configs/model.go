package configs

type Configs struct {
	App AppConfig
	DB  DbConfig
}

type AppConfig struct {
	Name string `json:"name"`
	Env  string `json:"env"`
	Port string `json:"port"`
}

type DbConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	DbName      string `json:"db_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DbIsMigrate bool   `json:"db_is_migrate"`
	DebugMode   bool   `json:"debug_mode"`
}
