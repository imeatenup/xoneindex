package config

const (
	CLIENT_ID     = "****-****-****-****-****"
	REDIRECT_URI  = "http://localhost:9000/callback"
	CLIENT_SECRET = "****-****-****"
	CALLBACK      = "/callback"
	STATE         = "****"
	ROOT_PATH     = "/.xoneindex" // eg "" "/.xoneindex"
	// scf
	XINDEX_TENCENTCLOUD_SECRETID  = "******"
	XINDEX_TENCENTCLOUD_SECRETKEY = "******"
)

var (
	Scopes = []string{"Files.ReadWrite.All", "offline_access"}
)
