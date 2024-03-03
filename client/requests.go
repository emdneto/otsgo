package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var ErrResponse ErrorResponse

var methods = map[string]string{
	"GET":  http.MethodGet,
	"POST": http.MethodPost,
}

func AgnosticRequest(a Auth, uri string, m string, b io.Reader) ([]byte, error) {
	client := http.Client{Timeout: 5 * time.Second}

	method := methods[m]

	body := b
	if method == "GET" {
		body = http.NoBody
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		log.Fatal(err)
	}

	if a.Enabled {
		req.SetBasicAuth(a.Username, a.Password)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		if err := json.Unmarshal(resBody, &ErrResponse); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		return nil, fmt.Errorf(string(fmt.Sprintf("%d %s", res.StatusCode, ErrResponse.Message)))
	}

	return resBody, nil
}
