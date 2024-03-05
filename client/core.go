package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var secret Secret
var secrets Secrets
var BurnSecRes BurnSecretResponse
var StsRes StatusRes

var tblheader = []string{"user", "state", "expires", "expired", "metadata", "passhphrase", "created", "sent", "ttl"}

func Login(a Auth) bool {
	uri := fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["status"])
	body, err := AgnosticRequest(a, uri, "GET", bytes.NewBufferString(""))

	if err != nil {
		fmt.Println(err)
		return false
	}
	if err := json.Unmarshal(body, &StsRes); err != nil {
		fmt.Println(err)
	}

	return true
}
func GetStatus(a Auth) bool {
	uri := fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["status"])
	body, err := AgnosticRequest(a, uri, "GET", bytes.NewBufferString(""))
	if err != nil {
		fmt.Println(err)
		return false
	}
	if err := json.Unmarshal(body, &StsRes); err != nil {
		fmt.Println(err)
	}

	// table output
	tblheader := []string{"Status", "Timestamp"}
	tbldata := [][]string{{StsRes.Status, fmt.Sprintf("%d", time.Now().Unix())}}
	fmtTableOutput(tblheader, tbldata)

	return true
}

func CreateSecret(a Auth, b SecretBody, g bool) {

	uri := fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["share"])

	if g {
		uri = fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["generate"])
	}

	data := url.Values{}
	data.Set("secret", b.Secret)
	data.Set("ttl", strconv.Itoa(b.Ttl))
	data.Set("recipient", b.Recipient)
	data.Set("passphrase", b.Passphrase)

	body, err := AgnosticRequest(a, uri, "POST", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(body, &secret); err != nil {
		fmt.Println(err)
	}

	// table output
	tblheader = []string{"user", "secret", "expires", "metadata", "sent"}
	tbldata := [][]string{}
	tbldata = append(tbldata, secret.toNewSecret())
	fmtTableOutput(tblheader, tbldata)

	if !a.Enabled {
		if err := writeHistory(secret.MetadataKey); err != nil {
			fmt.Printf("Error writing entry to file: %s\n", err)
		}
	}

	fmt.Println(secret.secretOtsCmd())
}

func BurnSecret(a Auth, b SecretBody) bool {
	uri := fmt.Sprintf("%s/%s/%s/%s/burn", BASE_URI, API_VERSION, ENDPOINTS["getmetadata"], b.Secret)
	data := url.Values{}
	fmt.Println(uri)
	body, err := AgnosticRequest(a, uri, "POST", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := json.Unmarshal(body, &BurnSecRes); err != nil {
		fmt.Println(err)
	}
	// table output
	fmt.Println(BurnSecRes.State.State)

	return true
}
func GetSecret(a Auth, b SecretBody) bool {
	uri := fmt.Sprintf("%s/%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["getsecret"], b.Secret)
	data := url.Values{}

	data.Set("SECRET_KEY", b.Secret)
	data.Set("passphrase", b.Passphrase)

	body, err := AgnosticRequest(a, uri, "POST", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := json.Unmarshal(body, &secret); err != nil {
		fmt.Println(err)
	}

	// table output
	fmt.Println(secret.Value)

	return true
}

func GetMetadata(a Auth, b SecretBody) {
	uri := fmt.Sprintf("%s/%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["getmetadata"], b.Secret)
	data := url.Values{}

	data.Set("METADATA_KEY", b.Secret)

	body, err := AgnosticRequest(a, uri, "POST", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	if err := json.Unmarshal(body, &secret); err != nil {
		fmt.Println(err)
	}

	tbldata := [][]string{}
	tbldata = append(tbldata, secret.toMetadata())
	fmtTableOutput(tblheader, tbldata)
}

func GetRecent(a Auth, b SecretBody) {
	uri := fmt.Sprintf("%s/%s/%s/recent", BASE_URI, API_VERSION, ENDPOINTS["getrecent"])
	tbldata := [][]string{}

	if !a.Enabled {
		fmt.Println("WARNING: Unable to locate credentials. You can configure credentials by running ots login.")
	}

	History, err := loadHistory(50)
	if err != nil {
		fmt.Printf("failed loadHistory: %v", err)
	}

	ch := make(chan Secret, len(History))
	errorCh := make(chan error)
	var wg sync.WaitGroup

	get_meta_uri := fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["getmetadata"])

	for _, item := range History {
		wg.Add(1)
		uri := fmt.Sprintf("%s/%s", get_meta_uri, item)
		data := url.Values{}
		data.Set("METADATA_KEY", item)

		go func(uri string, i string) {
			defer wg.Done()
			body, err := AgnosticRequest(a, uri, "POST", strings.NewReader(data.Encode()))
			if err != nil {
				if errors.Is(err, errors.New("404 Unknown secret")) || err.Error() == "404 Unknown secret" {
					return
				} else {
					errorCh <- err
					return
				}
			}

			if err := json.Unmarshal(body, &secret); err != nil {
				errorCh <- fmt.Errorf("failed unmarshalling JSON: %v\nResponse: %s", err, string(body))
				return
			}

			ch <- secret

		}(uri, item)
	}
	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	for range History {
		select {
		case secret, ok := <-ch:
			if !ok {
				break
			}
			tbldata = append(tbldata, secret.toMetadata())
		case err := <-errorCh:
			fmt.Printf("%v\n", err)

		}
	}

	if a.Enabled {
		body, err := AgnosticRequest(a, uri, "GET", strings.NewReader(""))
		if err != nil {
			fmt.Println(err)
		}

		if err := json.Unmarshal(body, &secrets); err != nil {
			fmt.Println(err)
		}
		for _, secret := range secrets {
			tbldata = append(tbldata, secret.toMetadata())
		}
	}

	if len(tbldata) == 0 {
		fmt.Printf("No secrets found.\n")
		return
	}

	fmtTableOutput(tblheader, tbldata)
}
