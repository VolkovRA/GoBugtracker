package main

// GetResponse - Получить базовый ответ сервера.
func GetResponse(result interface{}, errorCode int) interface{} {
	var m = make(map[string]interface{})
	if result != nil {
		m["results"] = result
	}
	if errorCode == 0 {
		m["success"] = true
	} else {
		m["success"] = false
		m["error"] = getErrorResponse(errorCode)
	}

	return m
}

// getErrorResponse - Получить описание ошибки в формате ответа сервера.
func getErrorResponse(errorCode int) interface{} {
	var m = make(map[string]interface{})
	m["code"] = errorCode
	m["message"] = ErrorMessages[errorCode]
	return m
}
