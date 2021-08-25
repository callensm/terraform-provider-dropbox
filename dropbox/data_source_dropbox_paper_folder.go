package dropbox

import (
	"fmt"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/paper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDropboxPaperFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDropboxPaperFolderRead,

		Schema: map[string]*schema.Schema{
			"doc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folders": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of folders that contain the document reference",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDropboxPaperFolderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := paper.NewRefPaperDoc(d.Get("doc_id").(string))
	info, err := client.DocsGetFolderInfo(opts)
	if err != nil {
		return fmt.Errorf("Paper Folder Data Failure: %s", err)
	}

	outputFolders := flattenPaperFolders(info.Folders)
	d.SetId(fmt.Sprintf("folder_ds:%s", d.Get("doc_id").(string)))
	d.Set("folders", outputFolders)

	return nil
}

func flattenPaperFolders(folders []*paper.Folder) []map[string]interface{} {
	outputFolders := make([]map[string]interface{}, 0, len(folders))
	for _, f := range folders {
		newFolder := make(map[string]interface{})
		newFolder["id"] = f.Id
		newFolder["name"] = f.Name
		outputFolders = append(outputFolders, newFolder)
	}
	return outputFolders
}
