package handlers

import (
	"html/template"
	"net/http"
)

var registeredUsers = []map[string]string{}

const templatesDir = "templates/"

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := templatesDir + tmpl
	tmplParsed, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmplParsed.Execute(w, data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		name := r.FormValue("name")
		email := r.FormValue("email")

		if name == "" || email == "" {
			renderTemplate(w, "register.html", "Заполните все поля!")
			return
		}

		newUser := map[string]string{"Name": name, "Email": email}
		registeredUsers = append(registeredUsers, newUser)

		data := map[string]string{"Message": "Регистрация прошла успешно!"}
		renderTemplate(w, "register.html", data)
		return
	}
	renderTemplate(w, "register.html", nil)
}

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	books := []string{"Book 1", "Book 2", "Book 3"}
	renderTemplate(w, "books.html", books)
}

func CartHandler(w http.ResponseWriter, r *http.Request) {
	cart := []string{"Selected Book 1", "Selected Book 2"}
	renderTemplate(w, "cart.html", cart)
}
