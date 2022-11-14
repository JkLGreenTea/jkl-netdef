package status

const (
	// Success - Все прошло хорошо, и (обычно) некоторые данные были возвращены.
	Success Status = "success"
	// Failed - Возникла проблема с представленными данными, или какое-то предварительное условие вызова API не было выполнено.
	Failed Status = "failed"
	// 	Error - При обработке запроса произошла ошибка, т.е. Было выдано исключение.
	Error Status = "error"
	// Stopped - Запрос не может быть обработан т.к сервер выключен.
	Stopped Status = "stopped"
	// Restarting - Запрос не может быть обработан т.к сервер перезагружается.
	Restarting Status = "restarting"
)

// Status - статус запроса.
type Status string
