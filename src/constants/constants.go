package constants

const (
	AllowedHeaders           = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
	ContentPrefixURL         = "/api"
	ServerInterfaceBindURL   = ":12403"
	BackupTargetTypeRsync    = "rsync"
	BackupTargetTypeTargz    = "targz"
	BackupTargetTypeGpgTargz = "gpgtargz"
)
