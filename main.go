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
	"regexp"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
        filelog := os.Getenv("LOG")
	if filelog == "" {
		filelog = "nolog"
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
	
	// Regular expression to check if input contains only numbers and hyphens so that user canot execute arbitary code
	printconfig := r.FormValue("printconfig")
	if printconfig == "" {
		fmt.Printf("Warning: User did not specify page range on sent request\n")
	} else {	
	    re := regexp.MustCompile(`^[0-9]+(-[0-9]+)*$`)	
	    if !re.MatchString(printconfig) {
		    json.NewEncoder(w).Encode(map[string]string{"error": "Invalid page range"})
		    fmt.Printf("debug: printconfig is: %s\n", printconfig)
		    return
	}
}
	// Regular expression to check if input contains only numbers so that user canot execute arbitary code
	printres := r.FormValue("res")
	if printres == "" {
		fmt.Printf("Warning: User did not specify quality on sent request, assuming normal quality\n")
		printres = "4"
	} else {	
	    re := regexp.MustCompile(`^[1-5]+$`)	
	    if !re.MatchString(printres) {
		    json.NewEncoder(w).Encode(map[string]string{"error": "Invalid quality range"})
		    fmt.Printf("debug: res is: %s\n", printres)
		    return
	}
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
	
	
	printingcmdPath := filepath.Join("/", "usr", "bin", "eprintcloned")
	filelog := os.Getenv("LOG")
	cmd := exec.Command(printingcmdPath, tempFile.Name(), filelog, printres, printconfig)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error printing file: %s\n", output)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Error printing file to printer: %s", output),
		})
		return
	}

	response := struct {
		Message string `json:"message"`
		Output  string `json:"output"`
	}{
		Message: fmt.Sprintf("File %s printed successfully", handler.Filename),
		Output:  string(output),
	}

	json.NewEncoder(w).Encode(response)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "To delete logs of your sent images from the instance you have sent them to, contact the owner of the instance to do it. "})
	return
}
