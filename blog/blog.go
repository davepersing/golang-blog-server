package blog

import (
	_ "database/sql"
	_ "encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	_ "io"
	_ "net/http"
)

type Entry struct {
	Id     int64
	Title  string
	Body   string
	Author string
}

//p := make([]byte, req.ContentLength)
// _, err := this.Ctx.Request.Body.Read(p)
func (e *Entry) PostEntry(dbmap *gorp.DbMap) (int64, error) {
	err := dbmap.Insert(e)
	if err != nil {
		fmt.Printf("Insert error: %v", err)
		return -1, err
	}

	return e.Id, nil
}

func GetEntry(dbmap *gorp.DbMap, id int64) (*Entry, error) {
	e, err := dbmap.Get(Entry{}, id)

	if err != nil {
		fmt.Printf("Error in GetBlog %v", err)
		return &Entry{}, err
	}

	return e.(*Entry), nil
}
