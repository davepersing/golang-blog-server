package main

import (
	"./blog"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

func main() {
	db, err := sql.Open("mysql", "root:@/blogs")
	if err != nil {
		panic(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbmap.AddTableWithName(blog.Entry{}, "entry").SetKeys(true, "Id")

	m := pat.New()

	m.Get("/blog/:id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.ParseInt(req.URL.Query().Get(":id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e, err := blog.GetEntry(dbmap, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}))

	m.Post("/blog", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)

		var e blog.Entry
		err := decoder.Decode(&e)
		fmt.Printf("PostEntry : %v", e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pk, err := e.PostEntry(dbmap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]int64{"id": pk}
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	}))

	http.Handle("/", m)
	fmt.Print("Blog Server listing on port 3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
