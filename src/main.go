package main

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

var smtpHost string
var smtpPort string
var smtpPswd string
var smtpEmail string

// Точка входа.
func main() {
	log.Println("Запуск сервера баг-трекера")

	// Тексты ошибок:
	InitErrorMessages()

	// Переменные окружения:
	var exists bool
	smtpHost, exists = os.LookupEnv("HOST")
	if !exists || smtpHost == "" {
		log.Fatalln("Не указана переменная окружения: HOST. (Адрес почтового SMTP сервера)")
	}

	smtpPort, exists = os.LookupEnv("PORT")
	if !exists || smtpPort == "" {
		log.Fatalln("Не указана переменная окружения: PORT. (Порт почтового SMTP сервера)")
	}

	smtpPswd, exists = os.LookupEnv("PSWD")
	if !exists || smtpPswd == "" {
		log.Fatalln("Не указана переменная окружения: PSWD. (Пароль пользователя на почтовом сервере)")
	}

	smtpEmail, exists = os.LookupEnv("EMAIL")
	if !exists || smtpEmail == "" {
		log.Fatalln("Не указана переменная окружения: EMAIL. (Почтовый ящик пользователя на почтовом сервере)")
	}

	// Запуск http:
	var ch chan error = make(chan error)
	go func() {
		var srv = &http.Server{
			Addr:           ":80",
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        http.HandlerFunc(handler),
		}
		ch <- srv.ListenAndServe()
	}()

	// Запуск https:
	go func() {
		var srv = &http.Server{
			Addr:           ":443",
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        http.HandlerFunc(handler),
		}
		ch <- srv.ListenAndServeTLS("ssl/certificate.crt", "ssl/private.key")
	}()

	// Один из серверов упал:
	log.Fatalln(<-ch)
}

// Обработчик всех входящих запросов
func handler(w http.ResponseWriter, r *http.Request) {
	var path = html.EscapeString(r.URL.Path)
	var ip = getIP(r)
	log.Println("Запрос " + r.Method + " " + path + " " + ip)

	// Ошиблись адресом:
	if r.Method != "POST" {
		send(w, GetResponse(nil, NotFound))
		return
	}

	// Парсер данных формы:
	var err = r.ParseForm()
	if err != nil {
		send(w, GetResponse(nil, BadRequest))
		log.Println("Ошибка парсера формы", err)
		return
	}
	var from = r.FormValue("from")
	var subject = r.FormValue("subject")
	var message = r.FormValue("message")

	// Валидация полученных данных:
	if from == "" {
		send(w, GetResponse(nil, FromNeed))
		return
	}
	if subject == "" {
		send(w, GetResponse(nil, SubjectNeed))
		return
	}
	if message == "" {
		send(w, GetResponse(nil, MessageNeed))
		return
	}

	// Отправка письма на почту:
	var body = []byte("To: " + smtpEmail + "\r\n" +
		"From: " + from + " <" + smtpEmail + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	var auth = smtp.PlainAuth("", smtpEmail, smtpPswd, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpEmail, []string{smtpEmail}, body)
	if err != nil {
		send(w, GetResponse(nil, InternalError))
		log.Println(err)
		return
	}

	// Сообщение отправлено:
	send(w, GetResponse(nil, 0))
}

// Получить IP клиента с учётом прокси.
func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// Отправить ответ сервера REST API.
func send(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Server", ServerName+"/"+ServerVersion)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(200)

	str, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write(str)
}
