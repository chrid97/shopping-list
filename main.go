package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func main() {
	page := template.Must(template.New("page").Parse(`
<!doctype html>
<html>

<head>
  <meta charset="utf-8">
  <title>Shopping List</title>
</head>

<body>
  <h1>Chris and Vivian's Shopping List</h1>

	<ul>
		{{range .}}
			<li>{{.}}</li>
		{{end}}
	</ul>

	<form method="POST" action="/add">
		<input name="item">
		<button type="submit">Add</button>
	</form>
</body>

</html>
	`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("shopping-list.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		items := strings.Split(string(data), "\n")

		err = page.Execute(w, items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// ------------------------------- Add Item -----------------------
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		item := strings.TrimSpace(r.FormValue("item"))
		// (todo) Display empty input is invalid error message on the frontend
		if item == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		file, err := os.OpenFile(
			"shopping-list.txt",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0644,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := file.Close(); err != nil {
			fmt.Println(err)
		}

		_, err = file.WriteString(item + "\n")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// ------------------------------- Delete Item -----------------------

	// ------------------------------- Start Server -----------------------
	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
