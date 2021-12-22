package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

var SecRes SecretResponse
var BurnSecRes BurnSecretResponse
var RecentSec Secrets

func GetStatus(a Auth) bool {
	uri := fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["status"])
	_, err := AgnosticRequest(a, uri, "GET", bytes.NewBufferString(""))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func CreateSecret(a Auth, b SecretBody, g bool) bool {

	uri := fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["share"])

	if g {
		uri = fmt.Sprintf("%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["generate"])
	}

	data := url.Values{}
	data.Set("secret", b.Secret)
	data.Set("ttl", b.Ttl)
	data.Set("recipient", b.Recipient)
	data.Set("passphrase", b.Passphrase)

	body, err := AgnosticRequest(a, uri, "POST", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := json.Unmarshal(body, &SecRes); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Share this link: %s/secret/%s\n", HOST, SecRes.SecretKey)
	fmt.Printf("Metatada (DO NOT share this): %s/private/%s\n", HOST, SecRes.MetadataKey)
	fmt.Println(SecRes.SecretKey)

	return true
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

	if err := json.Unmarshal(body, &SecRes); err != nil {
		fmt.Println(err)
	}
	fmt.Println(SecRes.Value)

	return true
}

func GetMetadata(a Auth, b SecretBody) bool {
	uri := fmt.Sprintf("%s/%s/%s/%s", BASE_URI, API_VERSION, ENDPOINTS["getmetadata"], b.Secret)
	data := url.Values{}

	data.Set("METADATA_KEY", b.Secret)

	body, err := AgnosticRequest(a, uri, "POST", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := json.Unmarshal(body, &SecRes); err != nil {
		fmt.Println(err)
	}
	fmt.Println(SecRes)

	return true
}

func OutputTable([][]string) {
	//
}
func GetRecent(a Auth, b SecretBody) bool {
	uri := fmt.Sprintf("%s/%s/%s/recent", BASE_URI, API_VERSION, ENDPOINTS["getrecent"])
	//data := url.Values{}

	body, err := AgnosticRequest(a, uri, "GET", strings.NewReader(""))
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := json.Unmarshal(body, &RecentSec); err != nil {
		fmt.Println(err)
	}
	data := [][]string{}
	for _, el := range RecentSec {
		updated_unix := fmt.Sprint(time.Unix(int64(el.Updated), 0))
		created_unix := fmt.Sprint(time.Unix(int64(el.Created), 0))

		data = append(data, []string{el.Custid, el.MetadataKey, strconv.Itoa(el.MetadataTTL), strconv.Itoa(el.SecretTTL), el.State, updated_unix, created_unix, fmt.Sprint(el.Recipient)})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Custid", "Metadata Key", "Metadata TTL", "Secret TTL", "State", "Updated", "Created", "Recipient"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	fmt.Println(len(RecentSec))

	return true
}
