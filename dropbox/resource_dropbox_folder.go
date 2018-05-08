package dropbox

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDropboxFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxFolderCreate,
		Read:   resourceDropboxFolderRead,
		Update: resourceDropboxFolderUpdate,
		Delete: resourceDropboxFolderDelete,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_rename": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDropboxFolderCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := files.New(*config)

	opts := &files.CreateFolderArg{
		Path:       d.Get("path").(string),
		Autorename: d.Get("auto_rename").(bool),
	}
	folder, err := client.CreateFolderV2(opts)
	if err != nil {
		return err
	}

	data := folder.Metadata
	d.SetId(data.Id)
	return resourceDropboxFolderRead(d, meta)
}

func resourceDropboxFolderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := files.New(*config)

	path := d.Get("path").(string)
	opts := files.NewGetMetadataArg(path)
	folder, err := client.GetMetadata(opts)
	if err != nil {
		return err
	}

	if folder.(*files.FolderMetadata).PathDisplay != "" {
		d.Set("path", folder.(*files.FolderMetadata).PathDisplay)
	} else {
		d.Set("path", path)
	}
	d.Set("id", folder.(*files.FolderMetadata).Id)
	d.Set("name", folder.(*files.FolderMetadata).Name)
	return nil
}

func resourceDropboxFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	// UNHIDE: config := meta.(*ProviderConfig).DropboxConfig
	// UNHIDE: client := files.New(*config)

	// TODO: Figure out how to modify path from SDK
	// d.Partial(true)
	// if d.HasChanged("path") {

	// }
	// d.Partial(false)

	return nil
}

func resourceDropboxFolderDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := files.New(*config)

	opts := &files.DeleteArg{Path: d.Get("path").(string)}
	_, err := client.DeleteV2(opts)
	if err != nil {
		return err
	}

	return nil
}
