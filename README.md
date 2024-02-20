# OTS Go client

otsgo is a simple CLI and API client for [One-Time Secret](https://onetimesecret.com/) written in Go. 

## Features

* Get system status
* Authentication
* Share, generate, get and burn Secrets as well as retrieve it's metadata

# Quickstart

## Getting otsgo from Webinstall

```
curl -sS https://webi.sh/ots | sh
```

## Getting otsgo from go install

```
$ go install github.com/emdneto/otsgo@latest
```

## Getting otsgo from binary
You may download otsgo binary from the [latest releases on Github](https://github.com/emdneto/otsgo/releases/latest).

# Using otsgo CLI

## Overview 
```
$ ots -h
A simple CLI and API client for One-Time Secret

Usage:
  ots [flags]
  ots [command]

Available Commands:
  burn        Burn a secret that has not been read yet
  completion  generate the autocompletion script for the specified shell
  get         Get secret, metadata or recent
  help        Help about any command
  login       Perform basic http auth and store credentials
  share       Share or generate a random secret
  status      Current status of the system

Flags:
      --config string   config file (default is $HOME/.otsgo.yaml)
  -h, --help            help for ots
  -t, --toggle          Help message for toggle
  -v, --version         Displays current version

Use "ots [command] --help" for more information about a command.
```

### Show OTS status

```
$ ots status

| STATUS  |
|---------|
| nominal |
```

## Share Secrets 

```
$ ots share -h
Share or generate a random secret

Usage:
  ots share [flags]

Flags:
  -f, --from-stdin          Read from stdin
  -g, --generate            Generate a short, unique secret. This is useful for temporary passwords, one-time pads, salts, etc.
  -h, --help                help for share
  -p, --passphrase string   a string that the recipient must know to view the secret. This value is also used to encrypt the secret and is bcrypted before being stored so we only have this value in transit.
  -r, --recipient string    an email address. We will send a friendly email containing the secret link (NOT the secret itself).
  -s, --secret string       the secret value which is encrypted before being stored. There is a maximum length based on your plan that is enforced (1k-10k)
  -t, --ttl int             the maximum amount of time, in seconds, that the secret should survive (i.e. time-to-live). Once this time expires, the secret will be deleted and not recoverable. (default 604800)
```

### Share a generated secret

```
$ ots share -g
```

### Share custom secret with ttl and passphrase

```
$ ots share -s hellosecret -t 300 -p hello
```

### Share secret from file
```
$ cat <<EOF | ots share -f -
secret: hello
seret: secret
EOF

$ echo "hellosecret" | ots share -f
```

### Burn secrets
```
$ ots burn METADATA_KEY
```

## Get secrets, metadata and recent
```
$ ots get -h
Get secret, metadata or recent

Usage:
  ots get [flags]
  ots get [command]

Available Commands:
  meta        Retrieve secret associated metadata
  recent      Retreive a list of recent metadata.
  secret      Retrieve a Secret

```

### Get secret value
```
$ ots get secret SECRET_KEY
```
### Get secret metadata
```
$ ots get meta METADATA_KEY
```
### Get recent secrets (requires auth)
```
$ ots get recent
```

## Authentication

### Auth with Environment Variables

otsgo will try to locate the credentials present in the environment variables. If found, every request will be made with HTTP Basic Authentication. If you get `404 Not authorized` in any command, probably your credentials are wrong.

```
$ export OTS_USER=demo; export OTS_TOKEN=xyz
$ ots get recent
$ ots share -g -r demo@demo.com
```
### Store auth credentials
Your password will be stored unencrypted in ~/.otsgo.yaml
```
$ ots login -u USERNAME -p API_TOKEN
$ ots get recent
$ ots share -g -r demo@demo.com
```

# Common aliases
```
alias oss="ots share secret"
alias osgs="ots share secret -g"
alias ogs="ots get secret"
alias obs="ots burn secret"
```