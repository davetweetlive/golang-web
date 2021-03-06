package checker

import (
	"html/template"
	"net/http"
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type Word struct {
	Name string
}

var tpl *template.Template

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/api/check", wordCheck)

	// serve public resources
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	// parse templates
	tpl = template.Must(template.ParseGlob("*.html"))
}

func index(res http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {

		var w Word
		w.Name = req.FormValue("new-word")

		ctx := appengine.NewContext(req)
		log.Infof(ctx, "WORD SUBMITTED: %v", w.Name)

		key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
		_, err := datastore.Put(ctx, key, &w)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tpl.ExecuteTemplate(res, "index.html", nil)
}

func wordCheck(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)

	// acquire the incoming word
	var w Word
	json.NewDecoder(req.Body).Decode(&w)
	log.Infof(ctx, "ENTERED wordCheck - w.Name: %v", w.Name)

	// check the incoming word against the datastore
	key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
	err := datastore.Get(ctx, key, &w)
	if err != nil {
		json.NewEncoder(res).Encode("false")
		return
	}
	json.NewEncoder(res).Encode("true")
}