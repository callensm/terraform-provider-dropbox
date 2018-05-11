package dropbox

import (
	"fmt"
	"strings"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDropboxFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxFileCreate,
		Read:   resourceDropboxFileRead,
		Update: resourceDropboxFileUpdate,
		Delete: resourceDropboxFileDelete,

		Schema: map[string]*schema.Schema{
			"content": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: convertContentToB64(),
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "add",
				Description:  "Valid formats modes: add, overwrite, and update",
				ValidateFunc: validateFileWriteMode(),
			},
			"auto_rename": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mute": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"hash": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A generated hash of the uploaded content",
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the uploaded data in bytes",
			},
		},
	}
}

func resourceDropboxFileCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := files.New(*config)

	content := d.Get("content").(string)
	reader := strings.NewReader(content)

	opts := &files.CommitInfo{
		Path:       d.Get("path").(string),
		Mode:       &files.WriteMode{Tagged: db.Tagged{Tag: d.Get("mode").(string)}},
		Autorename: d.Get("auto_rename").(bool),
		Mute:       d.Get("mute").(bool),
	}

	metadata, err := client.Upload(opts, reader)
	if err != nil {
		return fmt.Errorf("File Creation Failure: %s", err)
	}

	d.SetId(metadata.Id)
	d.Set("hash", metadata.ContentHash)
	d.Set("size", metadata.Size)
	return nil
}

func resourceDropboxFileRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDropboxFileUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDropboxFileDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := files.New(*config)

	opts := &files.DeleteArg{Path: d.Get("path").(string)}
	_, err := client.DeleteV2(opts)
	if err != nil {
		return fmt.Errorf("File Deletion Failure: %s", err)
	}

	return nil
}
