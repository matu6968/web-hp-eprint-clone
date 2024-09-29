package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/delete", handleDelete)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer file.Close()

	directory := r.FormValue("directory")
	if directory == "" {
		directory = "/"
	}

	tempFileName := filepath.Join(os.TempDir(), handler.Filename)
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	s3ClientPath := filepath.Join(".", "bin", "s3-client")
	configFilePath := filepath.Join(".", "bin", "s3config.toml")
	cmd := exec.Command(s3ClientPath, "-config", configFilePath, "-directory", directory, "-file", tempFile.Name())

	cmd.Args = append(cmd.Args, "-overwrite")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error uploading file: %s\n", output)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Error uploading file to S3: %s", output),
		})
		return
	}

	response := struct {
		Message string `json:"message"`
		Output  string `json:"output"`
	}{
		Message: fmt.Sprintf("File %s uploaded successfully", handler.Filename),
		Output:  string(output),
	}

	json.NewEncoder(w).Encode(response)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	filename := r.FormValue("filename")
	if filename == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Filename is required"})
		return
	}

	s3ClientPath := filepath.Join(".", "bin", "s3-client")
	configFilePath := filepath.Join(".", "bin", "s3config.toml")
	cmd := exec.Command(s3ClientPath, "-config", configFilePath, "-delete", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error deleting file: %s\n", output)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Error deleting file from S3: %s", output),
		})
		return
	}

	response := struct {
		Message string `json:"message"`
		Output  string `json:"output"`
	}{
		Message: fmt.Sprintf("File %s deleted successfully from S3", filename),
		Output:  string(output),
	}

	json.NewEncoder(w).Encode(response)
}
