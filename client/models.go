package client

type Auth struct {
	Username string
	Password string
	Enabled  bool
}

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

type SecretResponse struct {
	Custid             string   `json:"custid"`
	MetadataKey        string   `json:"metadata_key"`
	SecretKey          string   `json:"secret_key"`
	Ttl                int      `json:"ttl"`
	MetadataTtl        int      `json:"metadata_ttl"`
	SecretTtl          int      `json:"secret_ttl"`
	Recipient          []string `json:"recipient"`
	Created            int      `json:"created"`
	Updated            int      `json:"updated"`
	PassphraseRequired bool     `json:"passphrase_required"`
	Value              string   `json:"value"`
	State              string   `json:"state"`
}

type BurnSecretResponse struct {
	State SecretResponse `json:"state"`
}

type Secrets []struct {
	Custid      string   `json:"custid"`
	MetadataKey string   `json:"metadata_key"`
	TTL         int      `json:"ttl"`
	MetadataTTL int      `json:"metadata_ttl"`
	SecretTTL   int      `json:"secret_ttl"`
	State       string   `json:"state"`
	Updated     int      `json:"updated"`
	Created     int      `json:"created"`
	Recipient   []string `json:"recipient"`
}
