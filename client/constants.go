package client

const HOST = "https://onetimesecret.com"
const BASE_URI = "https://onetimesecret.com/api"
const API_VERSION = "v1"

var ENDPOINTS = map[string]string{
	"status":      "status",
	"share":       "share",
	"generate":    "generate",
	"getsecret":   "secret",
	"getmetadata": "private",
	"burn":        "private",
	"getrecent":   "private",
}
