package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		htmlContent := `
			<html>
				<title>Building Go API</title>
				<body>
					<h1>Building Go API</h1>
				</body>
			</html>
		`

		fmt.Fprint(w, htmlContent)
		w.Write([]byte(htmlContent))
	})

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
