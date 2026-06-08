package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request path:", r.URL.Path)
		if r.URL.Path == "/favicon.ico"{
			return
		}

		filename := r.URL.Query().Get("file")

		if filename == "" {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Bad Request"))
			return
		}

		filePath := "/data/" + filename
		fmt.Println("filename is: ", filePath)

		info, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Unauthorized Access Prohibited")
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Unauthorized Access Prohibited"))
			return
		}

		if info.IsDir() {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Cannot Access Dir"))
			return
		}

		file_bytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("error reading file content")
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("error reading file content"))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(file_bytes)
	})

	http.ListenAndServe(":8080", nil)

}
