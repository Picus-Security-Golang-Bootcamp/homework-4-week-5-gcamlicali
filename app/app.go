package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gcamlicali/RESTHW/app/model"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"

	postgres "github.com/gcamlicali/RESTHW/app/DB"
	. "github.com/gcamlicali/RESTHW/app/handler"
	Repo "github.com/gcamlicali/RESTHW/app/repo"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// App has router and db instances
type App struct {
	Router     *mux.Router
	DB         *gorm.DB
	BookRepo   *Repo.BookRepository
	AuthorRepo *Repo.AuthorRepository
}

// Initializer App inits with predefined ENV file
func (a *App) Initializer() error {

	//Reading ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Creating new DB
	db, err := postgres.NewPsqlDB()
	if err != nil {
		return fmt.Errorf("cannot open database: %v", err)
	}

	fmt.Println("Tables Migrated")
	a.DB = db

	a.AuthorRepo = Repo.NewAuthorRepository(db)
	a.BookRepo = Repo.NewBookRepository(db)

	fmt.Println("New Repo OK")

	// Migrating Author Tables
	err = model.AuthorDBMigrate(db)
	if err != nil {
		return fmt.Errorf("author database cannot migrated: %v", err)
	}
	// Migrating Book Tables
	err = model.BookDBMigrate(db)
	if err != nil {
		return fmt.Errorf("book database cannot migrated: %v", err)
	}

	a.AuthorRepo.FillAuthorData()
	a.BookRepo.FillBookData()
	fmt.Println("Fill Book OK")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	fmt.Println("Route Initialize OK")

	return nil
}

//Run Starts App
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

//getBooks Get all books from DB via own methods
func (a *App) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := a.BookRepo.FindAllBooks()

	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, books)
}

//getAuthors Get all Authors from DB via own methods
func (a *App) getAuthors(w http.ResponseWriter, r *http.Request) {
	author, err := a.AuthorRepo.GetAllAuthors()

	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, author)
}

//getBookByID Get given id book from DB via own methods
func (a *App) getBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	p, err := a.BookRepo.GetBookByID(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondError(w, http.StatusNotFound, "Book not found")
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusOK, p)
}

//getAuthorByID Get given id Author from DB via own methods
func (a *App) getAuthorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	author, err := a.AuthorRepo.GetAuthorByIdWithBooks(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondError(w, http.StatusNotFound, "Author not found")
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusOK, author)
}

//buyBookByID Updates given id book on DB via own methods
func (a *App) buyBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	type Quant struct {
		Quantity int `json:"quantity"`
	}
	q := Quant{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&q)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if q.Quantity == 0 {

		RespondError(w, http.StatusBadRequest, "Invalid Argument")
		return
	}

	quantity := uint(q.Quantity)

	p, err := a.BookRepo.BuyBookByID(id, quantity)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondError(w, http.StatusNotFound, "Book not found")
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusOK, p)
}

//searchBookByName Get given name book from DB via own methods
func (a *App) searchBookByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]

	p, err := a.BookRepo.SearchBooksByName(name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondError(w, http.StatusNotFound, "Book not found")
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusOK, p)
}

//searchAuthorByName Get given name Author from DB via own methods
func (a *App) searchAuthorByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]

	p, err := a.AuthorRepo.SearchAuthorByName(name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondError(w, http.StatusNotFound, "Author not found")
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusOK, p)
}

//deleteBookByID Soft deletes given id book in DB via own methods
func (a *App) deleteBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err = a.BookRepo.DeleteBookById(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondError(w, http.StatusNotFound, "Book not found")
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusOK, nil)
}

//deleteAuthorByID Soft deletes given id author in DB via own methods
func (a *App) deleteAuthorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err = a.AuthorRepo.DeleteAuthorById(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondError(w, http.StatusNotFound, "Author not found")
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusOK, nil)
}

// Initialize Routes
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/book", a.getBooks).Methods("GET")
	a.Router.HandleFunc("/book/name/{name}", a.searchBookByName).Methods("GET")
	a.Router.HandleFunc("/book/id/{id:[0-9]+}", a.getBookByID).Methods("GET")
	a.Router.HandleFunc("/book/id/{id:[0-9]+}", a.buyBookByID).Methods("PATCH")
	a.Router.HandleFunc("/book/id/{id:[0-9]+}", a.deleteBookByID).Methods("DELETE")
	a.Router.HandleFunc("/author", a.getAuthors).Methods("GET")
	a.Router.HandleFunc("/author/name/{name}", a.searchAuthorByName).Methods("GET")
	a.Router.HandleFunc("/author/id/{id:[0-9]+}", a.getAuthorByID).Methods("GET")
	a.Router.HandleFunc("/author/id/{id:[0-9]+}", a.deleteAuthorByID).Methods("DELETE")
}
