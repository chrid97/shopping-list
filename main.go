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
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Shopping List</title>

  <style>
    * {
      box-sizing: border-box;
    }

    html,
    body {
      height: 100%;
      margin: 0;
      font-family: sans-serif;
      background: #f5f5f5;
    }

    body {
      display: flex;
      justify-content: center;
    }

    .app {
      width: 100%;
      max-width: 600px;
      height: 100vh;
      display: flex;
      flex-direction: column;
      background: white;
    }

    header {
      padding: 1rem;
      border-bottom: 1px solid #ddd;
    }

    h1 {
      margin: 0;
      font-size: 1.5rem;
    }

    .items {
      flex: 1;
      overflow-y: auto;
      padding: 1rem;
      margin: 0;
      list-style: none;
    }

    .items li {
      padding: 0.75rem 1rem;
      margin-bottom: 0.5rem;
      border: 1px solid #ddd;
      border-radius: 4px;
    }

    form {
      display: flex;
      gap: 0.5rem;
      padding: 1rem;
      border-top: 1px solid #ddd;
      background: white;
    }

    input {
      flex: 1;
      padding: 0.75rem;
      font-size: 1rem;
    }

    button {
      padding: 0.75rem 1rem;
      font-size: 1rem;
      cursor: pointer;
    }
  </style>
</head>

<body>
  <div class="app">
    <header>
      <h1>Chris and Vivian's Shopping List</h1>
    </header>

    <ul class="items">
      {{range .}}
        {{if .}}
          <li>{{.}}</li>
        {{end}}
      {{end}}
    </ul>

    <form method="POST" action="/add">
      <input name="item" placeholder="Add item..." autocomplete="off">
      <button type="submit">Add</button>
    </form>
  </div>
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
			return
		}

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		_, err = file.WriteString(item + "\n")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
