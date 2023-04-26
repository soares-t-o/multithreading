package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ApiCep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
	Errors     error  `json:"errors"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Errors      error  `json:"errors"`
}

func getCepFromApiCep(cep string, data chan<- ApiCep) {

	req, err := http.NewRequest("GET", "https://cdn.apicep.com/file/apicep/"+cep+".json", nil)

	if err != nil {
		data <- ApiCep{Errors: err}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		data <- ApiCep{Errors: err}
	}

	if err != nil {
		data <- ApiCep{Errors: err}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		data <- ApiCep{Errors: err}
	}

	var apicep ApiCep
	json.Unmarshal(body, &apicep)
	data <- apicep

}

func getCepFromViaCep(cep string, data chan<- ViaCep) {
	req, err := http.NewRequest("GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)

	if err != nil {
		data <- ViaCep{Errors: err}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		data <- ViaCep{Errors: err}
	}

	if err != nil {
		data <- ViaCep{Errors: err}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		data <- ViaCep{Errors: err}
	}

	var viacep ViaCep
	json.Unmarshal(body, &viacep)
	data <- viacep
}

func main() {
	apicep := make(chan ApiCep)
	viacep := make(chan ViaCep)

	cep := "79400-000"

	go getCepFromApiCep(cep, apicep)
	go getCepFromViaCep(cep, viacep)

	select {

	case msg := <-apicep:
		fmt.Printf("APICEP: %+v\n", msg)
	case msg := <-viacep:
		fmt.Printf("VIACEP: %+v\n", msg)

	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

}
