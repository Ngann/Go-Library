package main

import (
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

//all imports package used in the main.go file
// "fmt" similar to C's printf/Scanf
// "net/http" provides HTTP client a dn server implementations.. GET, Head, post etc..
// "html/template" data driven templates for gernatinf HTML output safe against code injection, uses the same interfaces as the text/template package
//"database/sql" provide interfaces for SQL databases
//"github.com/mattn/go-sqlite3" sqlite3 driver conforming to the built-in database/sql interface
//"encoding/json" implements encoding/decoding of Json, Marshal and Unmarshal functions
//"net/url" parse URLs and implements query escaping
//"io/ioutil" allow to read a file from response body or some form of request output
//"encoding/xml"

import (
	"fmt"
	"net/http"
	"html/template"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"encoding/json"
	"net/url"
	"io/ioutil"
	"encoding/xml"
)

//Page struct with Name field is used in the route handler to create a new instance of Page with the name Gopher, then we will pass the page object as the third parameter in the Execute Template
//DBStatus
type Page struct {
	Name string
	DBStatus bool
}

type SearchResult struct {
	Title string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year string `xml:"hyr,attr"`
	ID string `xml:"owi,attr"`
}

func main() {
	// initialize the template variable with the call to method template.Parsefiles
	// ParseFiles will build the template object from a list of file names and return an error
	//Must will absorb the error from the ParseFiles and halt execution of the program if it cannot parse the template
	// the template is executed in the route handler
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "dev.db")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//p variable is used in template.ExecuteTemplate
		p := Page{Name: "Gopher"}
		if name := r.FormValue("name"); name != "" {
			p.Name = name
		}
		p.DBStatus = db.Ping() == nil

	// here we call execute template on a template's object with the reponse writer object
	// ExecuteTemplate returns ad error object, if the error not nil, alert user of internal server error 500
		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var results []SearchResult
		var err error

		if results, err = search(r.FormValue("search")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/books/add", func (w http.ResponseWriter, r *http.Request) {
		var book ClassifyBookResponse
		var err error

		if book, err = find(r.FormValue("id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err = db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = db.Exec("insert into books (pk, title, author, id, classification) values (?, ?, ?, ?, ?)",
			nil, book.BookData.Title, book.BookData.Author, book.BookData.ID, book.Classification.MostPopular)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}

type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

type ClassifyBookResponse struct {
	BookData struct {
		Title string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

func find(id string) (ClassifyBookResponse, error) {
	var c ClassifyBookResponse
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?summary=true&owi=" + url.QueryEscape(id))

	if err != nil {
		return ClassifyBookResponse{}, err
	}

	err = xml.Unmarshal(body, &c)
	return c, err
}

func search(query string) ([]SearchResult, error) {
	var c ClassifySearchResponse
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?summary=true&title=" + url.QueryEscape(query))

	if err != nil {
		return []SearchResult{}, err
	}

	err = xml.Unmarshal(body, &c)
	return c.Results, err
}

func classifyAPI(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get(url); err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}