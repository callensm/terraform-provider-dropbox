package dropbox

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/file_properties"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/hashicorp/terraform/helper/schema"
)

var folderPathPattern = "(/(.|[\r\n])*)|(ns:[0-9]+(/.*)?)"

func resourceDropboxFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxFolderCreate,
		Read:   resourceDropboxFolderRead,
		Delete: resourceDropboxFolderDelete,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateWithRegExp(folderPathPattern),
			},
			"auto_rename": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"folder_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"property_group_templates": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of template IDs corresponding to the associated folder property groups",
				Elem:        &schema.Schema{Type: schema.TypeString},
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

	if data.PropertyGroups != nil {
		d.Set("property_group_templates", flattenPropertyGroupIds(data.PropertyGroups))
	}

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
	d.Set("folder_id", folder.(*files.FolderMetadata).Id)
	d.Set("name", folder.(*files.FolderMetadata).Name)
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

func flattenPropertyGroupIds(groups []*file_properties.PropertyGroup) []string {
	groupIds := make([]string, 0, len(groups))
	for _, g := range groups {
		groupIds = append(groupIds, g.TemplateId)
	}
	return groupIds
}
