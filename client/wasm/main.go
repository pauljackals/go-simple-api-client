package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"syscall/js"
)

func document() js.Value {
	return js.Global().Get("document")
}
func preventDefault(event js.Value) {
	event.Call("preventDefault")
}
func getElementById(id string) js.Value {
	return document().Call("getElementById", id)
}
func addEventListener(element js.Value, eventType string, listener func(js.Value, []js.Value)interface{}) {
	element.Call("addEventListener", eventType, js.FuncOf(listener))
}
func createElement(elementType string) js.Value {
	return document().Call("createElement", elementType)
}
func appendChild(element js.Value, child js.Value) {
	element.Call("appendChild", child)
}

type Person struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

func createPerson(person Person, personsList js.Value, apiAddress string) {
	listElement := createElement("li")
	listElement.Set("id", fmt.Sprintf("person-%v", person.Id))
	listElement.Set("innerHTML", fmt.Sprintf("<span>%v</span>", person.Name))
	button := createElement("button")
	button.Set("textContent", "x")
	appendChild(listElement, button)
	appendChild(personsList, listElement)

	addEventListener(button, "click", func(this js.Value, args []js.Value) interface{} {
		go func() {
			client := http.Client{}
			request, _ := http.NewRequest("DELETE", fmt.Sprintf("%v?id=%v", apiAddress, person.Id), nil)
			_, responseError := client.Do(request)
			if responseError==nil {
				listElement.Call("remove")
			}
		}()
		return nil
	})
}

var ApiUrl string = "http://localhost:8080"

func main() {
	apiAddress := fmt.Sprintf("%v/person", ApiUrl)

	personsList := getElementById("persons")
	
	go func() {
		response, _ := http.Get(apiAddress)
		data, _ := io.ReadAll(response.Body)
		response.Body.Close()

		var persons []Person
		json.Unmarshal(data, &persons)
		
		for _, person := range persons {
			createPerson(person, personsList, apiAddress)
		}
	}()

	personForm := getElementById("person-form")
	addEventListener(personForm, "submit", func (this js.Value, args []js.Value) interface{} {
		event := args[0]
		preventDefault(event)

		target := event.Get("target")
		name := fmt.Sprint(target.Get("name").Get("value"))

		if len(name) == 0 {
			return nil
		}

		reqBody, _ := json.Marshal(map[string]string{"name": name})
	
		go func () {
			response, _ := http.Post(apiAddress, "application/json", bytes.NewBuffer(reqBody))
			data, _ := io.ReadAll(response.Body)
			response.Body.Close()

			var person Person
			json.Unmarshal(data, &person)
			createPerson(person, personsList, apiAddress)
		}()
		target.Call("reset")
	
		return nil
	})

	<-make(chan bool)
}
