package http_api_server

// Server - конфигурация Http api сервера.
type Server struct {
	// Domain используемый домен для Api (опционально).
	// Если домен не используется, укажите IP-адрес прокси сервера.
	Domain string `json:"domain" bson:"domain" yaml:"domain" form:"domain" description:"Domain используемый домен для Api (опционально). Если домен не используется, укажите IP-адрес прокси сервера."`

	// Наименование Api, используется при логгировании и url (по умолчанию 'api').
	Name string `json:"name" bson:"name" yaml:"name" form:"name" description:"Наименование Api, используется при логгировании и url (по умолчанию 'api')."`

	// Версия Api, используется в url (по умолчанию указывается разработчиком).
	Version string `json:"version" bson:"version" yaml:"version" form:"version" description:"Версия Api, используется в url (по умолчанию указывается разработчиком)."`

	// Gin mode:
	// gin.DebugMode указывает, что режим gin является отладочным;
	// gin.ReleaseMode указывает на то, что режим джина отключен;
	// gin.TestMode указывает, что режим джина является тестовым;
	GinMode string `json:"gin_mode" bson:"gin_mode" yaml:"gin_mode" form:"gin_mode" description:"Gin mode: gin.DebugMode указывает, что режим gin является отладочным; gin.ReleaseMode указывает на то, что режим джина отключен; gin.TestMode указывает, что режим джина является тестовым;"`

	// Адрес http api сервера (по умолчанию 0.0.0.0:9000).
	HttpAddr string `json:"http_addr" bson:"http_addr" yaml:"http_addr" form:"http_addr" description:"Адрес http api сервера (по умолчанию 0.0.0.0:9000)."`

	// Адрес https api сервера (по умолчанию 0.0.0.0:9001).
	HttpsAddr string `json:"https_addr" bson:"https_addr" yaml:"https_addr" form:"https_addr" description:"Адрес https api сервера (по умолчанию 0.0.0.0:9001)."`

	// Путь к директории файлов шаблонов (по умолчанию 'engine/api/http/templates/*/*').
	TemplateFiles []string `json:"template_files" bson:"template_files" yaml:"template_files" form:"template_files" description:"Путь к директории файлов шаблонов (по умолчанию 'engine/api/http/templates/*/*')."`

	// Путь к директории статичных файлов (по умолчанию path: 'engine/api/http/static/' root: '/static/').
	StaticFiles []StaticFiles `json:"static_files" bson:"static_files" yaml:"static_files" form:"static_files" description:"Путь к директории статичных файлов (по умолчанию path: 'engine/api/http/static/' root: '/static/')."`

	// Директория pprof файлов (по умолчанию 'system/pprof/api/').
	PprofFileDirectory string `json:"pprof_file_directory" bson:"pprof_file_directory" yaml:"pprof_file_directory" form:"pprof_file_directory" description:"Директория pprof файлов (по умолчанию 'system/pprof/api/')."`

	// Директория postman коллекций (по умолчанию 'system/postman/api/').
	PostmanCollectionsDirectory string `json:"postman_collections_directory" bson:"postman_collections_directory" yaml:"postman_collections_directory" form:"postman_collections_directory" description:"Директория postman коллекций (по умолчанию 'system/postman/api/')."`

	// Вкл. генерации postman коллекций (по умолчанию true).
	EnableGenPostmanCollections bool `json:"enable_gen_postman_collections" bson:"enable_gen_postman_collections" yaml:"enable_gen_postman_collections" form:"enable_gen_postman_collections" description:"Вкл. генерации postman коллекций (по умолчанию true)."`

	// Вкл. дописывания путей (по умолчанию false).
	EnablePathRewrite bool `json:"enable_path_rewrite" bson:"enable_path_rewrite" yaml:"enable_path_rewrite" form:"enable_path_rewrite" description:"Вкл. дописывания путей (по умолчанию false)."`

	// Вкл. pprof (по умолчанию true).
	EnablePprof bool `json:"enable_pprof" bson:"enable_pprof" yaml:"enable_pprof" form:"enable_pprof" description:"Вкл. pprof (по умолчанию true)."`

	// Вкл. кросс-доменные запросы (по умолчанию true).
	EnableCors bool `json:"enable_cors" bson:"enable_cors" yaml:"enable_cors" form:"enable_cors" description:"Вкл. кросс-доменные запросы (по умолчанию true)."`

	// Конфигурация cистемы доступа.
	SystemAccess *SystemAccess `json:"system_access" bson:"system_access" yaml:"system_access" form:"system_access" description:"Конфигурация cистемы доступа."`

	// Cors конфигурация кросс-доменных запросов.
	Cors *Cors `json:"cors" bson:"cors" yaml:"cors" form:"cors" description:"Конфигурация кросс-доменных запросов."`

	// TLSCertificateDir директория в которой хранятся tls сертификаты используемые для https проксирования
	// (по умолчанию 'system/certs/').
	TLSCertificateDir string `json:"tls_certificate_dir" bson:"tls_certificate_dir" yaml:"tls_certificate_dir" form:"tls_certificate_dir" description:"Директория в которой хранятся tls сертификаты используемые для https проксирования (по умолчанию 'system/certs/')."`

	// CaptchaImgDir путь к директории с хранилищем изображений используемых в капчи при проверки клиента
	// (по умолчанию 'engine/api/http/static/captcha/img/').
	CaptchaImgDir string `json:"captcha_img_dir" bson:"captcha_img_dir" yaml:"captcha_img_dir" form:"captcha_img_dir" description:"Путь к директории с хранилищем изображений используемых в капчи при проверки клиента (по умолчанию 'engine/api/http/static/captcha/img/')."`
}

// StaticFiles - статичные файлы.
type StaticFiles struct {
	// Путь к директории статичных файлов.
	Path string `json:"path" bson:"path" yaml:"path" form:"path" description:"Путь к директории статичных файлов."`

	// URL к директории статичных файлов.
	Root string `json:"root" bson:"root" yaml:"root" form:"root" description:"URL к директории статичных файлов."`
}

// SystemAccess - конфигурация системы доступа.
type SystemAccess struct {
	// Наименование системы доступа, используется при логгировании (по умолчанию 'SystemAccess').
	Title string `json:"title" bson:"title" yaml:"title" form:"title" description:"Наименование системы доступа, используется при логгировании (по умолчанию 'SystemAccess')."`

	// Конфигурация токена.
	Token *Token `json:"token" bson:"token" yaml:"token" form:"token" description:"Конфигурация токена."`
}

// Token - конфигурация токена системы доступа.
type Token struct {
	// Соль, используется для генерации ключа.
	HasSalt string `json:"has_salt" bson:"has_salt" yaml:"has_salt" form:"has_salt" description:"Соль, используется для генерации ключа."`

	// Время жизни токена.
	LifeTime float64 `json:"life_time" bson:"life_time" yaml:"life_time" form:"life_time" description:"Время жизни токена."`

	// Ключ подписи.
	SignedKey string `json:"signed_key" bson:"signed_key" yaml:"signed_key" form:"signed_key" description:"Ключ подписи."`
}

// Cors - конфигурация кросс-доменных запросов.
type Cors struct {
	AllowOrigins           []string `toml:"allow_origins" json:"allow_origins" yaml:"allow_origins"`
	AllowMethods           []string `toml:"allow_methods" json:"allow_methods" yaml:"allow_methods"`
	AllowHeaders           []string `toml:"allow_headers" json:"allow_headers" yaml:"allow_headers"`
	AllowCredentials       bool     `toml:"allow_credentials" json:"allow_credentials" yaml:"allow_credentials"`
	ExposeHeaders          []string `toml:"expose_headers" json:"expose_headers" yaml:"expose_headers"`
	MaxAge                 float64  `toml:"max_age" json:"max_age" yaml:"max_age"`
	AllowAllOrigins        bool     `toml:"allow_all_origins" json:"allow_all_origins" yaml:"allow_all_origins"`
	AllowFiles             bool     `toml:"allow_files" json:"allow_files" yaml:"allow_files"`
	AllowBrowserExtensions bool     `toml:"allow_browser_extensions" json:"allow_browser_extensions" yaml:"allow_browser_extensions"`
	AllowWildcard          bool     `toml:"allow_wildcard" json:"allow_wildcard" yaml:"allow_wildcard"`
	AllowWebSockets        bool     `toml:"allow_web_sockets" json:"allow_web_sockets" yaml:"allow_web_sockets"`
}
