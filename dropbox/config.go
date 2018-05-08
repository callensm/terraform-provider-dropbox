package dropbox

import (
	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
)

// ProviderConfig is passed in as meta data to all data
// source and resource management functions, containing
// the configuration struct for the Dropbox API
type ProviderConfig struct {
	DropboxConfig *db.Config
}
