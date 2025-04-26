package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
)

func serveHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./static/" + filename)
		if err != nil {
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	}
}

func submitEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	subject := r.FormValue("subject")
	message := r.FormValue("message")
	fmt.Printf("Name: %s\nEmail: %s\nSubject: %s\nMessage: %s\n", name, email, subject, message)
	w.WriteHeader(http.StatusOK)
	successHTML := `
<div class="flex items-center gap-2 px-4 py-2 rounded-lg bg-ctp-green/10 text-ctp-green border border-ctp-green/40 shadow-sm">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-ctp-green" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
    </svg>
    <span class="text-sm font-medium">Message sent successfully! We’ll get back to you soon.</span>
</div>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(successHTML))
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", "rexsy.bimq12@gmail.com")
	m.SetHeader("Cc", "rexsy.bimq12@gmail.com", "rexsy.bimq12@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text", message)
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("email_name"), os.Getenv("app_mail_password"))
	fmt.Println(os.Getenv("app_mail_password"))
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}
	// https://github.com/go-gomail/gomail
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	fmt.Println(os.Getenv("app_mail_password"))
	port := flag.String("port", "8080", "Port number")
	fs := http.FileServer(http.Dir("./static/imgs"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// CaskaydiaCove Nerd Font
	flag.Parse()
	// Directory you want to serve

	// Serve the files at root path "/"
	http.Handle("/blogposts", serveHTML("wip.html"))
	http.Handle("/dotfiles", serveHTML("dotfiles.html"))
	http.Handle("/projects", serveHTML("projects.html"))
	http.HandleFunc("/submit", submitEmailHandler)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			serveHTML("index.html").ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	// Start the server
	log.Printf("Starting server on :%s...", *port)
	addressPort := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(addressPort, nil))
}
