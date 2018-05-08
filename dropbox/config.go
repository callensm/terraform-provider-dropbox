package dropbox

import (
	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
)

type ProviderConfig struct {
	DropboxConfig *db.Config
}
