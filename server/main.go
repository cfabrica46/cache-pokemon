package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

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

		rows, err := db.Query("SELECT id,name,life,type,level FROM pokemons")

		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {

			var pokemon Pokemon

			err = rows.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Life, &pokemon.Type, &pokemon.Level)

			if err != nil {
				log.Fatal(err)
			}

			err = json.NewEncoder(w).Encode(pokemon)

			if err != nil {
				log.Fatal(err)
			}

		}

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
