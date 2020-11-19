package main

// ErrorMessages - Текстовое описание ошибок сервера API.
var ErrorMessages map[int]string

// Коды ошибок API.
const (

	// Некорректный запрос.
	BadRequest int = 400

	// Метод не найден.
	NotFound int = 404

	// Внутренняя ошибка сервера.
	InternalError int = 500

	// Требуется указать параметр: "from" для отправки сообщения. (От какого приложения это сообщение)
	FromNeed int = 1000

	// Требуется указать параметр: "subject" для отправки сообщения. (Тема сообщения, тип или класс ошибки)
	SubjectNeed int = 1001

	// Требуется указать параметр: "message" для отправки сообщения. (Текст сообщения)
	MessageNeed int = 1002
)

// InitErrorMessages - Инициализировать текстовое описание ошибок сервера API.
func InitErrorMessages() {
	ErrorMessages = map[int]string{
		BadRequest:    "Некорректный запрос",
		NotFound:      "Метод не найден",
		InternalError: "Внутренняя ошибка сервера",
		FromNeed:      "Требуется параметр: \"from\" для отправки сообщения. (От какого приложения это сообщение)",
		SubjectNeed:   "Требуется параметр: \"subject\" для отправки сообщения. (Тема сообщения, тип или класс ошибки)",
		MessageNeed:   "Требуется параметр: \"message\" для отправки сообщения. (Текст сообщения)",
	}
}
