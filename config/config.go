package config

type Environment struct {
	Server   `json:"server"`
	Database `json:"database"`
	API      `json:"api"`
	ErrorLog `json:"error_log"`
	Dict     `json:"dictionary"`
	Setting  `json:"settings"`
	JWT      `json:"jwt"`
	Activity `json:"activity"`
	Redis    `json:"redis"`
	RSA      `json:"rsa"`
	Email    `json:"email"`
	Document `json:"document"`
	GinMode  string `json:"gin_mode" env:"GIN_MOD" envDefault:"debug"`
}

type Server struct {
	Port           string `json:"port" env:"SERVER_PORT,required,notEmpty" `
	Addr           string `json:"addr" env:"SERVER_ADDR,required,notEmpty"`
	LogFormat      string `json:"log_format" env:"SERVER_LOG_FORMAT,required" `
	LogOutput      string `json:"log_output" env:"SERVER_LOG_OUTPUT,required"`
	LogLevel       string `json:"log_level" env:"SERVER_LOG_LEVEL,required" `
	LogIndentation bool   `json:"log_indentation" env:"SERVER_LOG_INDENTATION,required"`
	SendLibError   bool   `json:"send_lib_error" env:"SERVER_SEND_LIB_ERROR"`
	TimeZone       string `json:"time_zone" env:"SERVER_TIME_ZONE,required,notEmpty" `
}

type Database struct {
	DataDSN        string `json:"data_dsn" env:"DATABASE_DATA_DSN,required,notEmpty"`
	DataDBType     string `json:"data_db_type" env:"DATABASE_DATA_TYPE,required"`
	DataDBLog      string `json:"data_db_log" env:"DATABASE_DATA_LOG,required"`
	ActivityDSN    string `json:"activity_dsn" env:"DATABASE_ACTIVITY_DSN,required,notEmpty"`
	ActivityDBType string `json:"activity_db_type" env:"DATABASE_ACTIVITY_TYPE,required"`
	ActivityDBLog  string `json:"activity_db_log" env:"DATABASE_ACTIVITY_LOG,required"`
	AutoMigrate    bool   `json:"auto_migrate" env:"DATABASE_AUTO_MIGRATE,required" envtype:"bool"`
	ShowQueryLogs  bool   `json:"show_query_logs" env:"DATABASE_SHOW_QUERY_LOGS,required"`
}

type API struct {
	LogFormat      string `json:"log_format" env:"API_LOG_FORMAT,required"`
	LogOutput      string `json:"log_output" env:"API_LOG_OUTPUT,required"`
	LogLevel       string `json:"log_level" env:"API_LOG_LEVEL,required"`
	LogIndentation bool   `json:"log_indentation" env:"API_LOG_INDENTATION,required"`
}

type ErrorLog struct {
	LogFormat      string `json:"log_format" env:"API_LOG_FORMAT,required"`
	LogOutput      string `json:"log_output" env:"API_LOG_OUTPUT,required"`
	LogLevel       string `json:"log_level" env:"API_LOG_LEVEL,required"`
	LogIndentation bool   `json:"log_indentation" env:"API_LOG_INDENTATION,required"`
}

type Redis struct {
	Address     string `json:"address" env:"REDIS_ADDR,required"`
	Password    string `json:"password" env:"REDIS_PASSWORD,required"`
	DB          int    `json:"db" env:"REDIS_DB,required"`
	CacheApiTTL int    `json:"cache_api_ttl" env:"REDIS_CACHE_API_TTL,required"`
}

type Dict struct {
	TermsPath          string   `json:"terms_path" env:"DICT_TERMS_PATH,required"`
	DefaultLanguage    string   `json:"default_language" env:"DICT_DEFAULT_LANGUAGE,required"`
	TranslateInBackend bool     `json:"translate_in_backend" env:"DICT_TRANSLATE_IN_BACKEND,required"`
	Languages          []string `json:"languages" env:"DICT_LANGUAGES,required"`
}

type Setting struct {
	ExcelMaxRows     int    `json:"excel_max_rows" env:"SETTING_EXCEL_MAX_ROWS"`
	AdminUsername    string `json:"admin_username" env:"SETTING_ADMIN_USERNAME,required,notEmpty"`
	AdminPassword    string `json:"admin_password" env:"SETTING_ADMIN_PASSWORD,required,notEmpty"`
	PermissionTTL    int    `json:"permission_ttl" env:"SETTING_PERMISSIONS_TTL,required"`
	ResetPasswordUrl string `json:"reset_password_url" env:"SETTING_RESET_PASSWORD_URL,required"`
}

type JWT struct {
	PasswordSalt string `json:"password_salt" env:"JWT_PASSWORD_SALT,required,notEmpty"`
	SecretKey    string `json:"secret_key" env:"JWT_SECRET_KEY,required,notEmpty"`
	Expiration   uint64 `json:"expiration" env:"JWT_EXPIRATION,required,notEmpty"`
}

type Activity struct {
	Read      bool   `json:"read" env:"ACTIVITY_READ,required"`
	Write     bool   `json:"write" env:"ACTIVITY_WRITE,required"`
	BatchSize uint64 `json:"batch_size" env:"ACTIVITY_BATCH_SIZE,required"`
	SaveAfter uint64 `json:"save_after" env:"ACTIVITY_SAVE_AFTER,required"`
}

type RSA struct {
	CipherPrivateKey string `json:"cipher_private_key" env:"RSA_CIPHER_PRIVATE_KEY,required"`
	CipherPublicKey  string `json:"cipher_public_key" env:"RSA_CIPHER_PUBLIC_KEY,required"`
}

type Email struct {
	Username string `json:"username" env:"EMAIL_USERNAME,required"`
	Password string `json:"password" env:"EMAIL_PASSWORD,required"`
	Host     string `json:"host" env:"EMAIL_HOST,required"`
	Port     int    `json:"port" env:"EMAIL_PORT,required"`
}

type Document struct {
	CitiesDir string `json:"cities_dir" env:"DOCUMENTS_CITIES_DIR,required"`
	GiftsDir  string `json:"gifts_dir" env:"DOCUMENTS_GIFTS_DIR,required"`
}
