package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/losinggeneration/pipeln"
	helper "github.com/losinggeneration/pipeln/html/template"
)

func testTemplate(w http.ResponseWriter, req *http.Request) {
	t := template.New("example")
	helper.FuncMap(t)

	t, err := t.ParseFiles("./templates/example.html")
	if err != nil {
		log.Println("Unable to parse template: %v", err)
		return
	}

	t.ExecuteTemplate(w, "example.html", nil)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/template", testTemplate)

	n := negroni.Classic()

	n.Use(pipeln.NewAssets())
	n.UseHandler(mux)

	n.Run(":8080")
}
