package config

const (
	CLIENT_ID     = "****-****-****-****-****"
	REDIRECT_URI  = "http://localhost:9000/callback"
	CLIENT_SECRET = "****-****-****"
	CALLBACK      = "/callback"
	STATE         = "****"
	ROOT_PATH     = "/.xoneindex" // eg "" "/.xoneindex"
)

var (
	Scopes = []string{"Files.ReadWrite.All", "offline_access"}
)
