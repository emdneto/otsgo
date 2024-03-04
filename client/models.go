package client

import (
	"fmt"
	"strconv"
	"time"
)

type Auth struct {
	Username string
	Password string
	Enabled  bool
}

type History []string

type StatusRes struct {
	Status string `json:"status"`
}
type AuthYaml struct {
	Username string `yaml:"otsUser"`
	Password string `yaml:"otsToken"`
	//Enabled  bool   `yaml:"OTS_AUTH_ENABLED"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SecretBody struct {
	Secret     string
	Passphrase string
	Ttl        int
	Recipient  string
}

type Secrets []Secret

type Secret struct {
	Custid             string   `json:"custid"`
	MetadataKey        string   `json:"metadata_key"`
	SecretKey          string   `json:"secret_key"`
	Ttl                int      `json:"ttl"`
	MetadataTtl        int      `json:"metadata_ttl"`
	SecretTtl          int      `json:"secret_ttl"`
	Recipient          []string `json:"recipient"`
	CreatedAt          int      `json:"created"`
	UpdatedAt          int      `json:"updated"`
	PassphraseRequired bool     `json:"passphrase_required"`
	Value              string   `json:"value"`
	ReceivedAt         int      `json:"received"`
	State              string   `json:"state"`
}

type BurnSecretResponse struct {
	State Secret `json:"state"`
}

func (sm Secret) secretUrl() string {
	return fmt.Sprintf("%s/secret/%s\n", HOST, sm.SecretKey)
}

func (sm Secret) secretOtsCmd() string {
	return fmt.Sprintf("Share via ots: `ots get secret %s`", sm.SecretKey)
}

func (sm Secret) toNewSecret() []string {
	createdAt := time.Unix(int64(sm.CreatedAt), 0)
	expireAt := createdAt.Add(time.Duration(sm.Ttl/2) * time.Second)
	dateFmt := "2006-01-02 15:04:05"
	expireAtFmt := fmt.Sprint(expireAt.Format(dateFmt))
	return []string{
		sm.Custid,
		sm.secretUrl(),
		expireAtFmt,
		sm.MetadataKey,
		fmt.Sprint(sm.Recipient),
	}
}

func (sm Secret) toMetadata() []string {
	timeNow := time.Now()
	createdAt := time.Unix(int64(sm.CreatedAt), 0)
	expireAt := createdAt.Add(time.Duration(sm.Ttl/2) * time.Second)
	dateFmt := "2006-01-02 15:04:05"
	expireAtFmt := fmt.Sprint(expireAt.Format(dateFmt))
	createdAtFmt := fmt.Sprint(createdAt.Format(dateFmt))

	user := sm.Custid
	state := sm.State
	expired := strconv.FormatBool(timeNow.After(expireAt))
	metadataKey := sm.MetadataKey
	passhphrase := strconv.FormatBool(sm.PassphraseRequired)
	ttl := time.Duration(sm.SecretTtl) * time.Second

	return []string{
		user,
		state,
		expireAtFmt,
		expired,
		metadataKey,
		passhphrase,
		createdAtFmt,
		fmt.Sprint(sm.Recipient),
		ttl.String(),
	}
}
