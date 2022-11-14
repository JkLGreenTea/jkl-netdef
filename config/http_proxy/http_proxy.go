package http_proxy

import "time"

// HTTPProxy - конфигурация HTTP(S) прокси.
type HTTPProxy struct {
	// Host адрес на котором будет работать прокси для использования снаружи.
	// Используется в api для установления cookie.
	Host string `json:"host" bson:"host" yaml:"host" form:"host" description:"Адрес котором будет работать прокси для использования снаружи. Используется в api для установления cookie."`

	// ErrorPage путь к шаблону html странице с ошибкой
	// (по умолчанию 'system/templates/html/error.html').
	ErrorPage string `json:"error_page" bson:"error_page" yaml:"error_page" form:"error_page" description:"Путь к шаблону html странице с ошибкой (по умолчанию 'system/templates/html/error.html')."`

	// TechnicalWorksPage путь к шаблону html странице с тех. работами
	// (по умолчанию 'system/templates/html/technical_works.html').
	TechnicalWorksPage string `json:"technical_works_page" bson:"technical_works_page" yaml:"technical_works_page" form:"technical_works_page" description:"Путь к шаблону html странице с тех. работами (по умолчанию 'system/templates/html/technical_works.html')."`

	// CheckClientPage путь к шаблону html странице с проверкой клиент
	// (по умолчанию 'system/templates/html/check_client.html').
	CheckClientPage string `json:"check_client_page" bson:"check_client_page" yaml:"check_client_page" form:"check_client_page" description:"Путь к шаблону html странице с проверкой клиент (по умолчанию 'system/templates/html/check_client.html')."`

	// CheckCaptchaPage путь к шаблону html странице с проверкой клиента на каптчу
	// (по умолчанию 'system/templates/html/captcha.html').
	CheckCaptchaPage string `json:"check_captcha_page" bson:"check_captcha_page" yaml:"check_captcha_page" form:"check_captcha_page" description:"Путь к шаблону html странице с проверкой клиента на каптчу (по умолчанию 'system/templates/html/captcha.html')."`

	// ValueRequestPerMinuteForNoDataPage кол-во запросов за минуту при котором будет возвращаться
	// пустая страница с текстом (стили и js не используются)
	// (по умолчанию 300).
	ValueRequestPerMinuteForNoDataPage int `json:"value_request_per_minute_for_no_data_page" bson:"value_request_per_minute_for_no_data_page" yaml:"value_request_per_minute_for_no_data_page" form:"value_request_per_minute_for_no_data_page" description:"Кол-во запросов за минуту при котором будет возвращаться пустая страница с текстом (стили и js не используются) (по умолчанию 300)."`

	// AllowHosts список разрешенных хостов по которым будет производиться проксирование.
	AllowHosts []string `json:"allow_hosts" bson:"allow_hosts" yaml:"allow_hosts" form:"allow_hosts" description:"Список разрешенных хостов по которым будет производиться проксирование."`

	// RemoveHeaders список заголовков которые будут удаляться при проксировании
	// (по умолчанию - 'Connection'; 'Proxy-Connection'; 'Keep-Alive'; 'Proxy-Authenticate'; 'Proxy-Authorization'; 'Te'; 'Trailer'; 'Transfer-Encoding'; 'Upgrade';).
	RemoveHeaders []string `json:"remove_headers" bson:"remove_headers" yaml:"remove_headers" form:"remove_headers" description:"Список заголовков которые будут удаляться при проксировании (по умолчанию - 'Connection'; 'Proxy-Connection'; 'Keep-Alive'; 'Proxy-Authenticate'; 'Proxy-Authorization'; 'Te'; 'Trailer'; 'Transfer-Encoding'; 'Upgrade';)."`

	// TLSCertificateDir директория в которой хранятся tls сертификаты используемые для https проксирования
	// (по умолчанию 'system/certs/').
	TLSCertificateDir string `json:"tls_certificate_dir" bson:"tls_certificate_dir" yaml:"tls_certificate_dir" form:"tls_certificate_dir" description:"Директория в которой хранятся tls сертификаты используемые для https проксирования (по умолчанию 'system/certs/')."`

	// ClientLimit ограничение одновременно обрабатываемых клиентов
	// (по умолчанию 100).
	ClientLimit int `json:"client_limit" bson:"client_limit" yaml:"client_limit" form:"client_limit" description:"Ограничение одновременно обрабатываемых клиентов (по умолчанию 100)."`

	// EnableHostWhiteList если true, то хосты находящиеся в этом списке не будут блокироваться
	// (по умолчанию true).
	EnableHostWhiteList bool `json:"enable_host_white_list" bson:"enable_host_white_list" yaml:"enable_host_white_list" form:"enable_host_white_list" description:"Если true, то хосты находящиеся в этом списке не будут блокироваться (по умолчанию true)."`

	// EnableHardHostWhiteList если true, то только хосты находящиеся в этом списке будут иметь доступ к конечному серверу.
	// Обработка жеского белого списка происходит раньше чем проверка черного списка хостов, при включении обоих обработка
	// жеского белого списка произойдет раньше и на наличии в черном не будет проверяться
	// (по умолчанию false).
	EnableHardHostWhiteList bool `json:"enable_hard_host_white_list" bson:"enable_hard_host_white_list" yaml:"enable_hard_host_white_list" form:"enable_hard_host_white_list" description:"Если true, то только хосты находящиеся в этом списке будут иметь доступ к конечному серверу. Обработка жеского белого списка происходит раньше чем проверка черного списка хостов, при включении обоих обработка жеского белого списка произойдет раньше и на наличии в черном не будет проверяться (по умолчанию false)."`

	// EnableHostBlackList если true, то хосты находящиеся в этом списке не будут иметь доступ к конечному серверу.
	// (по умолчанию true).
	EnableHostBlackList bool `json:"enable_host_black_list" bson:"enable_host_black_list" yaml:"enable_host_black_list" form:"enable_host_black_list" description:"Если true, то хосты находящиеся в этом списке не будут иметь доступ к конечному серверу (по умолчанию true)."`

	// EnableLocationWhiteList если true, то только местоположения клиентов находящиеся в этом списке будут иметь
	// доступ к конечному серверу. Обработка белого списка происходит раньше чем проверка черного списка,
	// при включении обоих обработка белого списка произойдет раньше и на наличии в черном не будет проверяться
	// (по умолчанию false).
	EnableLocationWhiteList bool `json:"enable_location_white_list" bson:"enable_location_white_list" yaml:"enable_location_white_list" form:"enable_location_white_list" description:"Если true, то только местоположения клиентов находящиеся в этом списке будут иметь доступ к конечному серверу. Обработка белого списка происходит раньше чем проверка черного списка, при включении обоих обработка белого списка произойдет раньше и на наличии в черном не будет проверяться (по умолчанию false)."`

	// EnableLocationBlackList если true, то местоположения клиентов находящиеся в этом
	// списке не будут иметь доступ к конечному серверу
	// (по умолчанию false).
	EnableLocationBlackList bool `json:"enable_location_black_list" bson:"enable_location_black_list" yaml:"enable_location_black_list" form:"enable_location_black_list" description:"Если true, то местоположения клиентов находящиеся в этом списке не будут иметь доступ к конечному серверу (по умолчанию false)."`

	// Type тип прокси сервера, возможные варианты - 'Reverse';.
	Type string `json:"type" bson:"type" yaml:"type" form:"type" description:"Тип прокси сервера, возможные варианты - 'Reverse';."`

	// Protocol - протокол прок си сервера, возможные варианты - 'Http';.
	Protocol string `json:"protocol" bson:"protocol" yaml:"protocol" form:"protocol" description:"Протокол прок си сервера, возможные варианты - 'Http';."`

	// InAddr прослушиваемый адрес проки сервера.
	InAddr string `json:"in_addr" bson:"in_addr" yaml:"in_addr" form:"in_addr" description:"Прослушиваемый адрес проки сервера."`

	// OutAddr адрес перенаправления.
	OutAddr string `json:"out_addr" bson:"out_addr" yaml:"out_addr" form:"out_addr" description:"Адрес перенаправления."`

	// MethodsUnverifiedByToken список http методов которые не будут проверяться на наличии токена
	// (по умолчанию ['OPTIONS']).
	MethodsUnverifiedByToken []string `json:"methods_unverified_by_token" bson:"methods_unverified_by_token" yaml:"methods_unverified_by_token" form:"methods_unverified_by_token" description:"Список http методов которые не будут проверяться на наличии токена (по умолчанию ['OPTIONS'])."`

	// ReadTimeout - это максимальная продолжительность чтения всего
	// запрос, включая тело. Нулевое или отрицательное значение означает
	// тайм-аута не будет.
	//
	// Потому что ReadTimeout не позволяет обработчикам выполнять каждый запрос
	// решения о приемлемом крайнем сроке для каждого органа запроса или
	// скорости загрузки, большинство пользователей предпочтут использовать
	// ReadHeaderTimeout. Допустимо использовать их оба
	// (по умолчанию 30 * time.Second).
	ReadTimeout  time.Duration `json:"read_timeout" bson:"read_timeout" yaml:"read_timeout" form:"read_timeout" description:"Это максимальная продолжительность чтения всего запрос, включая тело. Нулевое или отрицательное значение означает тайм-аута не будет. Потому что ReadTimeout не позволяет обработчикам выполнять каждый запрос решения о приемлемом крайнем сроке для каждого органа запроса или скорости загрузки, большинство пользователей предпочтут использовать ReadHeaderTimeout. Допустимо использовать их оба (по умолчанию 30 * time.Second)."`

	// WriteTimeout - это максимальная продолжительность перед тайм-аутом
	// записи ответа. Он сбрасывается всякий раз, когда новый
	// заголовок запроса считывается. Как и ReadTimeout, он не
	// пусть обработчики принимают решения на основе каждого запроса.
	// Нулевое или отрицательное значение означает, что тайм-аута не будет
	// (по умолчанию 30 * time.Second).
	WriteTimeout time.Duration `json:"write_timeout" bson:"write_timeout" yaml:"write_timeout" form:"write_timeout" description:"Это максимальная продолжительность перед тайм-аутом записи ответа. Он сбрасывается всякий раз, когда новый заголовок запроса считывается. Как и ReadTimeout, он не пусть обработчики принимают решения на основе каждого запроса. Нулевое или отрицательное значение означает, что тайм-аута не будет (по умолчанию 30 * time.Second)."`

	// IdleTimeout - это максимальное время ожидания следующего запроса
	// при включенном режиме сохранения. Если время ожидания истекло
	// равно нулю, используется значение ReadTimeout. Если оба значения равны
	// нулю, тайм-аут отсутствует.
	// (по умолчанию 300 * time.Second).
	IdleTimeout time.Duration `json:"idle_timeout" bson:"idle_timeout" yaml:"idle_timeout" form:"idle_timeout" description:"Это максимальное время ожидания следующего запроса при включенном режиме сохранения. Если время ожидания истекло равно нулю, используется значение ReadTimeout. Если оба значения равны нулю, тайм-аут отсутствует (по умолчанию 300 * time.Second).."`
}
