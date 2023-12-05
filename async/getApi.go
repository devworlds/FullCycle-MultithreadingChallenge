package async

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const cep = "54774200"

func GetApi() {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()

	ch1 := make(chan string)
	ch2 := make(chan string)

	go getFirstApi(ch1)
	go getSecondApi(ch2)

	select {
	case viacep := <-ch1:
		println("Viacep: ", viacep)
	case apicep := <-ch2:
		println("Apicep: ", apicep)
	}

}

func getFirstApi(ch1 chan string) chan string {
	apiURL := "http://viacep.com.br/ws/" + cep + "/json/"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("A solicitação excedeu o tempo limite de 1 segundo.")
		} else {
			fmt.Println("Erro ao fazer a solicitação:", err)
		}
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ch1 <- string(respData)
	return ch1
}

func getSecondApi(ch2 chan string) chan string {
	apiURL := "https://cdn.apicep.com/file/apicep/" + cep + ".json"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("A solicitação excedeu o tempo limite de 1 segundo.")
		} else {
			fmt.Println("Erro ao fazer a solicitação:", err)
		}
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ch2 <- string(respData)
	return ch2
}
