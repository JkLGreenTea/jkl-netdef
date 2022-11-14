package loggers

// Loggers - конфигурация логгеров.
type Loggers struct {
	// Api конфигурация api логгера.
	Api *Api `json:"api" bson:"api" yaml:"api" form:"api" description:"Конфигурация api логгера."`

	// Global конфигурация глобального логгера.
	Global *Global `json:"global" bson:"global" yaml:"global" form:"global" description:"Конфигурация глобального логгера."`

	// HttpProxy конфигурация http proxy логгера.
	HttpProxy *HttpProxy `json:"http_proxy" bson:"http_proxy" yaml:"http_proxy" form:"http_proxy" description:"Конфигурация http proxy логгера."`
}

// Api - конфигурация api логгера.
type Api struct {
	// Title название логгера
	// (по умолчанию 'Api-Log').
	Title string `json:"title" bson:"title" yaml:"title" form:"title" description:"Название логгера (по умолчанию 'Api-Log')."`

	// LogLevel уровень логгирования
	//(по умолчанию DEBUG').
	LogLevel string `json:"log_level" bson:"log_level" yaml:"log_level" form:"log_level" description:"Уровень логгирования (по умолчанию DEBUG')."`

	// TimeFormat формат времени для вывода лога
	// (по умолчанию 'Monday, 02 Jan 2006 15:04:05').
	TimeFormat string `json:"time_format" bson:"time_format" yaml:"time_format" form:"time_format" description:"Формат времени для вывода лога (по умолчанию 'Monday, 02 Jan 2006 15:04:05')."`

	// LogFilePath путь к директории сохранения лога
	// (по умолчанию 'system/logs/api/').
	LogFilePath string `json:"log_file_path" bson:"log_file_path" yaml:"log_file_path" form:"log_file_path" description:"Путь к директории сохранения лога (по умолчанию 'system/logs/api/')."`

	// EnableCallerPath конфигурация включения отображения путя вызова.
	EnableCallerPath *EnableCallerPath `json:"enable_caller_path" bson:"enable_caller_path" yaml:"enable_caller_path" form:"enable_caller_path" description:"Конфигурация включения отображения путя вызова."`

	// LogEnableOutputFile конфигурация включения записи в файл.
	EnableOutputFile *LogEnableOutputFile `json:"enable_output_file" bson:"enable_output_file" yaml:"enable_output_file" form:"enable_output_file" description:"Конфигурация включения записи в файл."`

	// EnableOutputInTerminal если true, тогда логгирование будет отображаться в терминале (грузит syslog)
	// (по умолчанию true).
	EnableOutputInTerminal bool `json:"enable_output_in_terminal" bson:"enable_output_in_terminal" yaml:"enable_output_in_terminal" form:"enable_output_in_terminal" description:"EnableOutputTerminal если true, тогда логгирование будет отображаться в терминале (грузит syslog) (по умолчанию true)."`

	// LogColors цвета заголовков лог вывода.
	Colors *LogColors `json:"colors" bson:"colors" yaml:"colors" form:"colors" description:"Цвета заголовков лог вывода."`
}

// Global - конфигурация глобального логгера.
type Global struct {
	// Title название логгера
	// (по умолчанию 'Global-Log').
	Title string `json:"title" bson:"title" yaml:"title" form:"title" description:"Название логгера (по умолчанию 'Api-Log')."`

	// LogLevel уровень логгирования
	//(по умолчанию DEBUG').
	LogLevel string `json:"log_level" bson:"log_level" yaml:"log_level" form:"log_level" description:"Уровень логгирования (по умолчанию DEBUG')."`

	// TimeFormat формат времени для вывода лога
	// (по умолчанию 'Monday, 02 Jan 2006 15:04:05').
	TimeFormat string `json:"time_format" bson:"time_format" yaml:"time_format" form:"time_format" description:"Формат времени для вывода лога (по умолчанию 'Monday, 02 Jan 2006 15:04:05')."`

	// LogFilePath путь к директории сохранения лога
	// (по умолчанию 'system/logs/api/').
	LogFilePath string `json:"log_file_path" bson:"log_file_path" yaml:"log_file_path" form:"log_file_path" description:"Путь к директории сохранения лога (по умолчанию 'system/logs/api/')."`

	// EnableCallerPath включение отображения путя вызова.
	EnableCallerPath *EnableCallerPath `json:"enable_caller_path" bson:"enable_caller_path" yaml:"enable_caller_path" form:"enable_caller_path" description:"Включение отображения путя вызова."`

	// LogEnableOutputFile включение записи в файл.
	EnableOutputFile *LogEnableOutputFile `json:"enable_output_file" bson:"enable_output_file" yaml:"enable_output_file" form:"enable_output_file" description:"Включение записи в файл."`

	// EnableOutputInTerminal если true, тогда логгирование будет отображаться в терминале (грузит syslog)
	// (по умолчанию true).
	EnableOutputInTerminal bool `json:"enable_output_in_terminal" bson:"enable_output_in_terminal" yaml:"enable_output_in_terminal" form:"enable_output_in_terminal" description:"EnableOutputTerminal если true, тогда логгирование будет отображаться в терминале (грузит syslog) (по умолчанию true)."`

	// LogColors цвета заголовков лог вывода.
	Colors *LogColors `json:"colors" bson:"colors" yaml:"colors" form:"colors" description:"Цвета заголовков лог вывода."`
}

// HttpProxy - конфигурация http proxy логгера.
type HttpProxy struct {
	// Title название логгера
	// (по умолчанию 'HttpProxy-Log').
	Title string `json:"title" bson:"title" yaml:"title" form:"title" description:"Название логгера (по умолчанию 'Api-Log')."`

	// LogLevel уровень логгирования
	//(по умолчанию DEBUG').
	LogLevel string `json:"log_level" bson:"log_level" yaml:"log_level" form:"log_level" description:"Уровень логгирования (по умолчанию DEBUG')."`

	// TimeFormat формат времени для вывода лога
	// (по умолчанию 'Monday, 02 Jan 2006 15:04:05').
	TimeFormat string `json:"time_format" bson:"time_format" yaml:"time_format" form:"time_format" description:"Формат времени для вывода лога (по умолчанию 'Monday, 02 Jan 2006 15:04:05')."`

	// LogFilePath путь к директории сохранения лога
	// (по умолчанию 'system/logs/api/').
	LogFilePath string `json:"log_file_path" bson:"log_file_path" yaml:"log_file_path" form:"log_file_path" description:"Путь к директории сохранения лога (по умолчанию 'system/logs/api/')."`

	// EnableCallerPath включение отображения путя вызова.
	EnableCallerPath *EnableCallerPath `json:"enable_caller_path" bson:"enable_caller_path" yaml:"enable_caller_path" form:"enable_caller_path" description:"Включение отображения путя вызова."`

	// LogEnableOutputFile включение записи в файл.
	EnableOutputFile *LogEnableOutputFile `json:"enable_output_file" bson:"enable_output_file" yaml:"enable_output_file" form:"enable_output_file" description:"Включение записи в файл."`

	// EnableOutputInTerminal если true, тогда логгирование будет отображаться в терминале (грузит syslog)
	// (по умолчанию true).
	EnableOutputInTerminal bool `json:"enable_output_in_terminal" bson:"enable_output_in_terminal" yaml:"enable_output_in_terminal" form:"enable_output_in_terminal" description:"EnableOutputTerminal если true, тогда логгирование будет отображаться в терминале (грузит syslog) (по умолчанию true)."`

	// LogColors цвета заголовков лог вывода.
	Colors *LogColors `json:"colors" bson:"colors" yaml:"colors" form:"colors" description:"Цвета заголовков лог вывода."`
}

// LogEnableOutputFile - включение записи в файл.
type LogEnableOutputFile struct {
	// DEBUG если true, то запись 'DEBUG' лога будет сохраняться в журнал
	// (по умолчанию true).
	DEBUG bool `json:"debug" bson:"debug" yaml:"debug" form:"debug" description:"Если true, то запись 'DEBUG' лога будет сохраняться в журнал (по умолчанию true)."`

	// INFO если true, то запись 'INFO' лога будет сохраняться в журнал
	// (по умолчанию true).
	INFO bool `json:"info" bson:"info" yaml:"info" form:"info" description:"Если true, то запись 'INFO' лога будет сохраняться в журнал (по умолчанию true)."`

	// WARN если true, то запись 'WARN' лога будет сохраняться в журнал
	// (по умолчанию true).
	WARN bool `json:"warn" bson:"warn" yaml:"warn" form:"warn" description:"Если true, то запись 'WARN' лога будет сохраняться в журнал (по умолчанию true)."`

	// ERROR если true, то запись 'ERROR' лога будет сохраняться в журнал
	// (по умолчанию true).
	ERROR bool `json:"error" bson:"error" yaml:"error" form:"error" description:"Если true, то запись 'ERROR' лога будет сохраняться в журнал (по умолчанию true)."`

	// FATAL если true, то запись 'FATAL' лога будет сохраняться в журнал
	// (по умолчанию true).
	FATAL bool `json:"fatal" bson:"fatal" yaml:"fatal" form:"fatal" description:"Если true, то запись 'FATAL' лога будет сохраняться в журнал (по умолчанию true)."`
}

// LogColors - цвета заголовков лог вывода.
type LogColors struct {
	// DEBUG цвет заголовка для отображения в терминале 'DEBUG' лога
	// (по умолчанию 'Green').
	DEBUG string `json:"debug" bson:"debug" yaml:"debug" form:"debug" description:"Цвет заголовка для отображения в терминале 'DEBUG' лога (по умолчанию 'Green')."`

	// INFO цвет заголовка для отображения в терминале 'INFO' лога
	// (по умолчанию 'Cyan').
	INFO string `json:"info" bson:"info" yaml:"info" form:"info" description:"Цвет заголовка для отображения в терминале 'INFO' лога (по умолчанию 'Cyan')."`

	// WARN цвет заголовка для отображения в терминале 'WARN' лога
	// (по умолчанию 'Yellow').
	WARN string `json:"warn" bson:"warn" yaml:"warn" form:"warn" description:"Цвет заголовка для отображения в терминале 'WARN' лога (по умолчанию 'Yellow')."`

	// ERROR цвет заголовка для отображения в терминале 'ERROR' лога
	// (по умолчанию 'Red').
	ERROR string `json:"error" bson:"error" yaml:"error" form:"error" description:"Цвет заголовка для отображения в терминале 'ERROR' лога (по умолчанию 'Red')."`

	// FATAL цвет заголовка для отображения в терминале 'FATAL' лога
	// (по умолчанию 'Purple').
	FATAL string `json:"fatal" bson:"fatal" yaml:"fatal" form:"fatal" description:"Цвет заголовка для отображения в терминале 'FATAL' лога (по умолчанию 'Purple')."`

	// ALogColor цвет заголовка для отображения в терминале api лога
	// (по умолчанию 'Blue').
	ALogColor string `json:"a_log_color" bson:"a_log_color" yaml:"a_log_color" form:"a_log_color" description:"Цвет заголовка для отображения в терминале api лога (по умолчанию 'Blue')."`

	// PLogColor цвет заголовка для отображения в терминале proxy лога
	// (по умолчанию 'Cyan').
	PLogColor string `json:"p_log_color" bson:"p_log_color" yaml:"p_log_color" form:"p_log_color" description:"Цвет заголовка для отображения в терминале proxy лога (по умолчанию 'Cyan')."`

	// GLogColor цвет заголовка для отображения в терминале глобального лога
	// (по умолчанию 'Purple').
	GLogColor string `json:"g_log_color" bson:"g_log_color" yaml:"g_log_color" form:"g_log_color" description:"Цвет заголовка для отображения в терминале глобального лога (по умолчанию 'Purple')."`

	// GET цвет заголовка для отображения в терминале лога 'GET' запросов
	// (по умолчанию 'Green').
	GET string `json:"get" bson:"get" yaml:"get" form:"get" description:"Цвет заголовка для отображения в терминале лога 'GET' запросов (по умолчанию 'Green')"`

	// POST цвет заголовка для отображения в терминале лога 'POST' запросов
	// (по умолчанию 'Cyan').
	POST string `json:"post" bson:"post" yaml:"post" form:"post" description:"Цвет заголовка для отображения в терминале лога 'POST' запросов (по умолчанию 'Cyan')"`

	// PUT цвет заголовка для отображения в терминале лога 'PUT' запросов
	// (по умолчанию 'Purple').
	PUT string `json:"put" bson:"put" yaml:"put" form:"put" description:"Цвет заголовка для отображения в терминале лога 'PUT' запросов (по умолчанию 'Purple')"`

	// DELETE цвет заголовка для отображения в терминале лога 'DELETE' запросов
	// (по умолчанию 'Red').
	DELETE string `json:"delete" bson:"delete" yaml:"delete" form:"delete" description:"Цвет заголовка для отображения в терминале лога 'DELETE' запросов (по умолчанию 'Red')"`

	// PATCH цвет заголовка для отображения в терминале лога 'PATCH' запросов
	// (по умолчанию 'Yellow').
	PATCH string `json:"patch" bson:"patch" yaml:"patch" form:"patch" description:"Цвет заголовка для отображения в терминале лога 'PATCH' запросов (по умолчанию 'Yellow')"`

	// OPTIONS цвет заголовка для отображения в терминале лога 'OPTIONS' запросов
	// (по умолчанию 'Blue').
	OPTIONS string `json:"options" bson:"options" yaml:"options" form:"options" description:"Цвет заголовка для отображения в терминале лога 'OPTIONS' запросов (по умолчанию 'Blue')"`

	// HEAD цвет заголовка для отображения в терминале лога 'HEAD' запросов
	// (по умолчанию 'Gray').
	HEAD string `json:"head" bson:"head" yaml:"head" form:"head" description:"Цвет заголовка для отображения в терминале лога 'HEAD' запросов (по умолчанию 'Gray')"`

	// HTTPCode100 цвет заголовка для отображения в терминале лога http запросов со статусами 100
	// (по умолчанию 'Green').
	HTTPCode100 string `json:"http_code_100" bson:"http_code_100" yaml:"http_code_100" form:"http_code_100" description:"Цвет заголовка для отображения в терминале лога http запросов со статусами 100 (по умолчанию 'Green')."`

	// HTTPCode200 цвет заголовка для отображения в терминале лога http запросов со статусами 200
	// (по умолчанию 'Cyan').
	HTTPCode200 string `json:"http_code_200" bson:"http_code_200" yaml:"http_code_200" form:"http_code_200" description:"Цвет заголовка для отображения в терминале лога http запросов со статусами 200 (по умолчанию 'Cyan')."`

	// HTTPCode300 цвет заголовка для отображения в терминале лога http запросов со статусами 300
	// (по умолчанию 'Purple').
	HTTPCode300 string `json:"http_code_300" bson:"http_code_300" yaml:"http_code_300" form:"http_code_300" description:"Цвет заголовка для отображения в терминале лога http запросов со статусами 300 (по умолчанию 'Purple')."`

	// HTTPCode400 цвет заголовка для отображения в терминале лога http запросов со статусами 400
	// (по умолчанию 'Yellow').
	HTTPCode400 string `json:"http_code_400" bson:"http_code_400" yaml:"http_code_400" form:"http_code_400" description:"Цвет заголовка для отображения в терминале лога http запросов со статусами 400 (по умолчанию 'Yellow')."`

	// HTTPCode500 цвет заголовка для отображения в терминале лога http запросов со статусами 500
	// (по умолчанию 'Red').
	HTTPCode500 string `json:"http_code_500" bson:"http_code_500" yaml:"http_code_500" form:"http_code_500" description:"Цвет заголовка для отображения в терминале лога http запросов со статусами 500 (по умолчанию 'Red')."`
}

// EnableCallerPath - включение отображения путя вызова.
type EnableCallerPath struct {
	// DEBUG если true, то в отображении 'DEBUG' лога будет отображаться место вызова
	// (по умолчанию true).
	DEBUG bool `json:"debug" bson:"debug" yaml:"debug" form:"debug" description:"Если true, то в отображении 'DEBUG' лога будет отображаться место вызова (по умолчанию true)."`

	// INFO если true, то в отображении 'INFO' лога будет отображаться место вызова
	// (по умолчанию false).
	INFO bool `json:"info" bson:"info" yaml:"info" form:"info" description:"Если true, то в отображении 'INFO' лога будет отображаться место вызова (по умолчанию false)."`

	// WARN если true, то в отображении 'WARN' лога будет отображаться место вызова
	// (по умолчанию false).
	WARN bool `json:"warn" bson:"warn" yaml:"warn" form:"warn" description:"Если true, то в отображении 'WARN' лога будет отображаться место вызова (по умолчанию false)."`

	// ERROR если true, то в отображении 'ERROR' лога будет отображаться место вызова
	// (по умолчанию true).
	ERROR bool `json:"error" bson:"error" yaml:"error" form:"error" description:"Если true, то в отображении 'ERROR' лога будет отображаться место вызова (по умолчанию true)."`

	// FATAL если true, то в отображении 'FATAL' лога будет отображаться место вызова
	// (по умолчанию true).
	FATAL bool `json:"fatal" bson:"fatal" yaml:"fatal" form:"fatal" description:"Если true, то в отображении 'FATAL' лога будет отображаться место вызова (по умолчанию true)."`
}
