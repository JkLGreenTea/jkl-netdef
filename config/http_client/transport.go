package http_client

import "time"

// Transport - HTTP транспорт клиента для проксирования запросов.
type Transport struct {
	// TLSHandshakeTimeout задает максимальное время ожидания для
	// рукопожатия TLS. Ноль означает отсутствие тайм-аута
	// (по умолчанию 15 * time.Second).
	TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout" bson:"tls_handshake_timeout" yaml:"tls_handshake_timeout" form:"tls_handshake_timeout" description:"TLSHandshakeTimeout задает максимальное время ожидания для рукопожатия TLS. Ноль означает отсутствие тайм-аута (по умолчанию 15 * time.Second)."`

	// DisableKeepAlives, если true, отключает HTTP keep-alives и
	// будет использовать соединение с сервером только для одного HTTP-запроса.
	// Это не связано с аналогичным именем TCP keep-alives
	// (по умолчанию true).
	DisableKeepAlives bool `json:"disable_keep_alives" bson:"disable_keep_alives" yaml:"disable_keep_alives" form:"disable_keep_alives" description:"DisableKeepAlives, если true, отключает HTTP keep-alives и будет использовать соединение с сервером только для одного HTTP-запроса. Это не связано с аналогичным именем TCP keep-alives (по умолчанию true)."`

	// DisableCompression, если true, запрещает транспорту
	// запрос сжатия с помощью «Accept-Encoding: gzip»,
	// если запрос не содержит заголовок.
	// Если Транспорт запрашивает gzip на свой собственный запрос и получает ответ в формате gzip,
	// это декодируется в Response.Body
	// (по умолчанию false).
	DisableCompression bool `json:"disable_compression" bson:"disable_compression" yaml:"disable_compression" form:"disable_compression" description:"DisableCompression, если true, запрещает транспорту запрос сжатия с помощью «Accept-Encoding: gzip», если запрос не содержит заголовок. Если Транспорт запрашивает gzip на свой собственный запрос и получает ответ в формате gzip, это декодируется в Response.Body (по умолчанию false)."`

	// MaxIdleConns управляет максимальным количеством простоев (keep-alive)
	// соединения на всех хостах. Ноль означает отсутствие ограничений
	// (по умолчанию 100).
	MaxIdleConns int `json:"max_idle_conns" bson:"max_idle_conns" yaml:"max_idle_conns" form:"max_idle_conns" description:"Управляет максимальным количеством простоев (keep-alive) соединения на всех хостах. Ноль означает отсутствие ограничений (по умолчанию 100). "`

	// MaxIdleConnsPerHost, если он не равен нулю, управляет максимальным простоем
	// (keep-alive) соединения для каждого хоста
	// (по умолчанию 100).
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host" bson:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host" form:"max_idle_conns_per_host" description:"Ecли он не равен нулю, управляет максимальным простоем (keep-alive) соединения для каждого хоста (по умолчанию 100)."`

	// MaxConnsPerHost ограничивает общее количество
	// подключений на хост, включая соединения в состояниях подключения,
	// активного и бездействующего. При нарушении лимита циферблаты будут заблокированы.
	//
	// Ноль означает отсутствие ограничений
	// (по умолчанию 0).
	MaxConnsPerHost int `json:"max_conns_per_host" bson:"max_conns_per_host" yaml:"max_conns_per_host" form:"max_conns_per_host" description:"Ограничивает общее количество  подключений на хост, включая соединения в состояниях подключения, активного и бездействующего. При нарушении лимита циферблаты будут заблокированы.  Ноль означает отсутствие ограничений (по умолчанию 0)."`

	// IdleConnTimeout — максимальное время бездействия
	// (keep-alive) соединение будет простаивать перед закрытием сам.
	// Ноль означает отсутствие ограничений
	// (по умолчанию 90 * time.Second).
	IdleConnTimeout time.Duration `json:"idle_conn_timeout" bson:"idle_conn_timeout" yaml:"idle_conn_timeout" form:"idle_conn_timeout" description:"Максимальное время бездействия (keep-alive) соединение будет простаивать перед закрытием сам. Ноль означает отсутствие ограничений (по умолчанию 90 * time.Second)."`

	// ResponseHeaderTimeout, если не нулевое значение, указывает количество
	// времени ожидания заголовков ответа сервера после полного
	// написание запроса (включая его текст, если таковой имеется). Это
	// время не включает время на чтение тела ответа
	// (по умолчанию 0).
	ResponseHeaderTimeout time.Duration `json:"response_header_timeout" bson:"response_header_timeout" yaml:"response_header_timeout" form:"response_header_timeout" description:"Если не нулевое значение, указывает количество ремени ожидания заголовков ответа сервера после полного написание запроса (включая его текст, если таковой имеется). Это время не включает время на чтение тела ответа (по умолчанию 0)."`

	// ExpectContinueTimeout, если не нулевое значение, указывает количество
	// времени ожидания первых заголовков ответа сервера после полной
	// записи заголовков запроса, если запрос имеет
	// заголовок "Expect: 100-continue". Ноль означает отсутствие тайм-аута и
	// приводит к немедленной отправке тела без
	// ожидание утверждения сервером.
	// Это время не включает время отправки заголовка запроса
	// (по умолчанию 1 * time.Second).
	ExpectContinueTimeout time.Duration `json:"expect_continue_timeout" bson:"expect_continue_timeout" yaml:"expect_continue_timeout" form:"expect_continue_timeout" description:"Если не нулевое значение, указывает количество времени ожидания первых заголовков ответа сервера после полной записи заголовков запроса, если запрос имеет заголовок 'Expect: 100-continue'. Ноль означает отсутствие тайм-аута и приводит к немедленной отправке тела без ожидание утверждения сервером. Это время не включает время отправки заголовка запроса (по умолчанию 1 * time.Second)."`

	// MaxResponseHeaderBytes - задает ограничение на количество
	// байтов ответа, разрешенных в ответе заголовках сервера.
	// Ноль означает использование ограничения по умолчанию
	// (по умолчанию 0).
	MaxResponseHeaderBytes int64 `json:"max_response_header_bytes" bson:"max_response_header_bytes" yaml:"max_response_header_bytes" form:"max_response_header_bytes" description:"Задает ограничение на количество байтов ответа, разрешенных в ответе заголовках сервера. Ноль означает использование ограничения по умолчанию (по умолчанию 0)."`

	// WriteBufferSize - определяет размер используемого буфера записи при записи в транспорт.
	// Если значение равно нулю, используется значение по умолчанию.
	// (по умолчанию 0).
	WriteBufferSize int `json:"write_buffer_size" bson:"write_buffer_size" yaml:"write_buffer_size" form:"write_buffer_size" description:"Определяет размер используемого буфера записи при записи в транспорт. Если значение равно нулю, используется значение по умолчанию (по умолчанию 0)."`

	// ReadBufferSize - определяет размер используемого буфера чтения
	// при считывании с транспортного средства.
	// Если значение равно нулю, используется значение по умолчанию
	// (по умолчанию 0).
	ReadBufferSize int `json:"read_buffer_size" bson:"read_buffer_size" yaml:"read_buffer_size" form:"read_buffer_size" description:"Определяет размер используемого буфера чтения при считывании с транспортного средства. Если значение равно нулю, используется значение по умолчанию (по умолчанию 0)."`

	// ForceAttemptHTTP2 - определяет, включен ли HTTP/2 при ненулевом
	// предоставляется функция Dial, DialTLS или DialContext или TLSClientConfig.
	// По умолчанию использование любых этих полей консервативно отключает HTTP/2.
	// Чтобы использовать пользовательский Dial или конфигурацию TLS и по-прежнему пытаться использовать HTTP/2
	// обновления, установите для этого значение true
	// (по умолчанию false).
	ForceAttemptHTTP2 bool `json:"force_attempt_http2" bson:"force_attempt_http2" yaml:"force_attempt_http2" form:"force_attempt_http2" description:"Определяет, включен ли HTTP/2 при ненулевом предоставляется функция Dial, DialTLS или DialContext или TLSClientConfig. По умолчанию использование любых этих полей консервативно отключает HTTP/2. Чтобы использовать пользовательский Dial или конфигурацию TLS и по-прежнему пытаться использовать HTTP/2 обновления, установите для этого значение true (по умолчанию false)."`
}
