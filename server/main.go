package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var b bytes.Buffer

type Pokemon struct {
	ID, Life, Level int
	Name, Type      string
}

func main() {

	var err error

	db, err = open()

	if err != nil {
		log.Fatal(err)
	}

	go deleteCache()

	http.HandleFunc("/", getPokes)

	fmt.Println("Listening on :8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func getPokes(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		if b.Len() != 0 {
			fmt.Println("En Cache")
			fmt.Fprintf(w, "%s\n", b.Bytes())
			return
		}

		fmt.Println("En Databases")

		rows, err := db.Query("SELECT id,name,life,type,level FROM pokemons")

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		for rows.Next() {

			var pokemon Pokemon

			err = rows.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Life, &pokemon.Type, &pokemon.Level)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			dataJSON, err := json.Marshal(pokemon)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			_, err = b.Write(dataJSON)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			fmt.Fprintf(w, "%s\n", dataJSON)

		}

	}

}

func deleteCache() {

	for {
		b.Reset()
		time.Sleep(time.Second * 10)
	}

}

func open() (databases *sql.DB, err error) {

	archivo, err := os.Open("databases.db")

	if err != nil {
		if os.IsNotExist(err) {

			databases, err = migracion()

			if err != nil {

				archivo.Close()
				return
			}

			return
		}
		return
	}
	defer archivo.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	return
}

func migracion() (databases *sql.DB, err error) {

	archivoDB, err := os.Create("databases.db")

	if err != nil {
		return
	}
	archivoDB.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	archivoSQL, err := os.Open("databases.sql")

	if err != nil {
		return
	}

	defer archivoSQL.Close()

	contendio, err := io.ReadAll(archivoSQL)

	if err != nil {
		return
	}

	_, err = databases.Exec(string(contendio))

	if err != nil {
		return
	}

	return

}
