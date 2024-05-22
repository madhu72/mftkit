package main

import (
	"fmt"
	"io"
	"log"
	"mftkit/mft"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./uploads/"

// uploadHandler handles file upload.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll(uploadPath, os.ModePerm)
	out, err := os.Create(filepath.Join(uploadPath, header.Filename))
	if err != nil {
		http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Unable to write the file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded successfully"))
}

// downloadHandler handles file download.
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is missing", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, file)
}

// curl -F "file=@path/to/your/file" http://localhost:8080/upload
// curl -O "http://localhost:8080/download?file=yourfilename"
func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	m := mft.NewMFT()
	fmt.Printf("Version:%v\n", m.Version())
	owner, err := m.GetFileOwner("readme.md")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("File Owner:", owner)
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
