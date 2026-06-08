package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := ":8080"
	dir := "/Users/myousuf/dev/obs-notes/obsdian-notes/dist/" 

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, filepath.Clean(r.URL.Path))
		
		// Get file info
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		
		// If it's a directory, return 404
		if info.IsDir() {
			http.NotFound(w, r)
			return
		}
		
		// Serve the file
		http.ServeFile(w, r, path)
	})

	println("Server running at http://localhost" + port)
	println("Directories will NOT be served (404)")
	http.ListenAndServe(port, nil)
}
