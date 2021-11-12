package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type PersonDTO struct {
	Name string `json:"name"`
}
type Person struct {
	Id int `json:"id"`
	PersonDTO
}

type Database struct {
	persons []*Person
}
func (database *Database) init() {
	database.persons = []*Person{}
}
func (database *Database) findAll() *[]*Person {
	return &database.persons
}
func (database *Database) save(personDTO *PersonDTO) *Person {
	var id int
	persons := database.persons
	if len(persons)>0 {
		id = persons[len(persons)-1].Id + 1
	} else {
		id = 0
	}
	person := &Person{id, *personDTO}
	database.persons = append(persons, person)
	return person
}
func (database *Database) delete(id int) *Person {
	persons := database.persons
	found := -1
	for i, record := range persons {
		if record.Id == id {
			found = i
			break
		}
	}
	if found < 0 {
		return nil

	} else {
		person := persons[found]
		database.persons = append(persons[:found], persons[found+1:]...)
		return person
	}
}

func parseQueryId(r *http.Request) (int, error) {
	r.ParseForm()
	idRaw := r.Form["id"]
	if len(idRaw) == 0 {
		return -1, errors.New("Missing id")
	} else {
		id, idError := strconv.Atoi(idRaw[0])
		if idError==nil && id>=0 {
			return id, nil
		} else {
			return -1, errors.New("Invalid id")
		}
	}
}

func main() {
	var port int
	flag.IntVar(&port, "p", 8080, "Server port")
	flag.Parse()

	database := Database{}
	database.init()

	http.HandleFunc("/person", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.Header().Set("Access-Control-Allow-Methods", "DELETE")

		switch r.Method {
		case "GET":
			resJSON, _ := json.Marshal(database.findAll())

			rw.WriteHeader(http.StatusOK)
			rw.Write(resJSON)
			return

		case "POST":
			bodyRaw, bodyRawError := io.ReadAll(r.Body)
			r.Body.Close()
			var personDTO PersonDTO
			jsonError := json.Unmarshal(bodyRaw, &personDTO)

			if bodyRawError == nil && jsonError == nil &&
			personDTO.Name != "" {
				person := database.save(&personDTO)
				resJSON, _ := json.Marshal(person)

				rw.WriteHeader(http.StatusCreated)
				rw.Write(resJSON)
				return

			} else {
				rw.WriteHeader(http.StatusUnprocessableEntity)
				rw.Write([]byte("Could not parse body"))
				return
			}

		case "DELETE":
			id, idError := parseQueryId(r)
			if idError==nil {
				person := database.delete(id)
				if person == nil {
					rw.WriteHeader(http.StatusNotFound)
					rw.Write([]byte("Could not find person"))
					return

				} else {
					resJSON, _ := json.Marshal(person)
					rw.WriteHeader(http.StatusOK)
					rw.Write(resJSON)
					return
				}

			} else {
				rw.WriteHeader(http.StatusUnprocessableEntity)
				rw.Write([]byte(idError.Error()))
				return
			}
		}
	})

	fmt.Printf("Server listening on port %v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}