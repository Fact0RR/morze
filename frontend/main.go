package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed index.html script.js style.css env.js send_mes.png
var staticFiles embed.FS

func main() {
	// Создаем подфайловую систему, которая будет показывать файлы из корня
	// Это нужно чтобы файлы были доступны по корневому пути
	staticFS, err := fs.Sub(staticFiles, ".")
	if err != nil {
		log.Fatal(err)
	}

	// Создаем файловый сервер
	fileServer := http.FileServer(http.FS(staticFS))

	// Обрабатываем все запросы
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если запрос к корню, отдаем index.html
		if r.URL.Path == "/" {
			http.ServeFileFS(w, r, staticFiles, "index.html")
			return
		}
		// Иначе используем файловый сервер
		fileServer.ServeHTTP(w, r)
	})

	// Запускаем сервер на порту 8080
	log.Println("Сервер запущен на http://localhost:8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}