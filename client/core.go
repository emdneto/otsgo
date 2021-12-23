package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

var SecRes SecretResponse
var BurnSecRes BurnSecretResponse
var RecentSec Secrets
var StsRes StatusRes

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

	data := []string{StsRes.Status}
	header := []string{"Status"}
	OutputTable(header, data)

	return true
}

func CreateSecret(a Auth, b SecretBody, g bool) bool {

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
		return false
	}
	timeNow := time.Now().Unix()
	if err := json.Unmarshal(body, &SecRes); err != nil {
		fmt.Println(err)
	}
	total := timeNow + int64(SecRes.SecretTtl)
	expiresIn := fmt.Sprint(time.Unix(int64(total), 0))
	secKeyUrl := fmt.Sprintf("%s/secret/%s\n", HOST, SecRes.SecretKey)
	privKeyUrl := fmt.Sprintf("%s/private/%s\n", HOST, SecRes.MetadataKey)
	templateData := []string{SecRes.Custid, expiresIn, secKeyUrl, privKeyUrl}
	header := []string{"User", "Expires in", "Share this link", "Private Metadata (DO NOT share this)"}
	fmt.Printf("\nShareable link copied to clipboard!\n")
	fmt.Print("\n")
	clipboard.Write(clipboard.FmtText, []byte(secKeyUrl))
	clipboard.Read(clipboard.FmtText)
	OutputTable(header, templateData)

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
	templateData := []string{SecRes.Custid,
		SecRes.MetadataKey,
		fmt.Sprint(SecRes.MetadataTtl),
		fmt.Sprint(SecRes.SecretTtl),
		SecRes.State,
		fmt.Sprint(SecRes.Updated),
		fmt.Sprint(SecRes.Created),
		fmt.Sprint(SecRes.Recipient)}
	header := []string{"Custid", "Metadata Key", "Metadata TTL", "Secret TTL", "State", "Updated", "Created", "Recipient"}
	OutputTable(header, templateData)

	return true
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
	header := []string{"Custid", "Metadata Key", "Metadata TTL", "Secret TTL", "State", "Updated", "Created", "Recipient"}
	OutputBulkTable(header, data)

	fmt.Println(len(RecentSec))

	return true
}
