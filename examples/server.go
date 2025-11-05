package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/h1", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for key, values := range r.Form {
			fmt.Println("Form", key, values)
		}
		for key, values := range r.PostForm {
			fmt.Println("PostForm", key, values)
		}

		w.Write([]byte("Form: "))
		w.Write([]byte(r.Form.Encode()))
		w.Write([]byte(" "))
		w.Write([]byte("PostForm: "))
		w.Write([]byte(r.PostForm.Encode()))
		w.Write([]byte("\n"))
	})
	http.ListenAndServe(":8080", nil)
}
