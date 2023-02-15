package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Message struct {
	ID  int64
	Msg string
}

type ViaCEP struct {
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
}

type CDNapiCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	c1 := make(chan []byte)
	c2 := make(chan []byte)

	go func() {
		cep := "09130-110"
		req, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		c1 <- res
	}()

	go func() {
		cep := "09130110"
		req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		c2 <- res
	}()
	select {
	case resp := <-c1:
		var cep CDNapiCEP
		err := json.Unmarshal(resp, &cep)
		if err != nil {
			panic(err)
		}
		fmt.Printf(`Resultado obtido primeiro pela API "https://cdn.apicep.com/file/apicep/" + cep + ".json": %v`, cep)
	case resp := <-c2:
		var cep ViaCEP
		err := json.Unmarshal(resp, &cep)
		if err != nil {
			panic(err)
		}
		fmt.Printf(`Resultado obtido primeiro pela API "http://viacep.com.br/ws/" + cep + "/json/": %v`, cep)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}
