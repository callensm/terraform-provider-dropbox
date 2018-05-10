package dropbox

import (
	"encoding/base64"
	"fmt"
	"os"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/paper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDropboxPaperDoc() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxPaperDocCreate,
		Read:   resourceDropboxPaperDocRead,
		Update: resourceDropboxPaperDocUpdate,
		Delete: resourceDropboxPaperDocDelete,

		Schema: map[string]*schema.Schema{
			"content_file": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: convertDocContentToB64(),
			},
			"parent_folder": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"import_format": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Valid formats include: html, markdown, plain_text, other",
				ValidateFunc: validateDocImportFormat(),
			},
			"doc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"revision": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDropboxPaperDocCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	var tag paper.ImportFormat
	switch format := d.Get("import_format").(string); format {
	case "html":
	case "markdown":
	case "plain_text":
	case "other":
		tag = paper.ImportFormat{Tagged: db.Tagged{Tag: format}}
	default:
		return fmt.Errorf("Doc Creation Failure: Invalid import format given for paper document creation: %s", format)
	}

	contentFile := d.Get("content_file").(string)
	reader, err := os.Open(contentFile)
	if err != nil {
		return err
	}

	opts := &paper.PaperDocCreateArgs{
		ImportFormat:   &tag,
		ParentFolderId: d.Get("parent_folder").(string),
	}
	results, err := client.DocsCreate(opts, reader)
	if err != nil {
		return err
	}

	d.SetId(results.DocId)

	return resourceDropboxFolderRead(d, meta)
}

func resourceDropboxPaperDocRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	var tag paper.ExportFormat
	switch format := d.Get("import_format").(string); format {
	case "html":
		tag.Tag = paper.ExportFormatHtml
	case "markdown":
		tag.Tag = paper.ExportFormatMarkdown
	case "plain_text":
	case "other":
		tag.Tag = paper.ExportFormatOther
	default:
		return fmt.Errorf("Invalid import format given for paper document creation: %s", format)
	}

	opts := &paper.PaperDocExport{
		RefPaperDoc:  *paper.NewRefPaperDoc(d.Id()),
		ExportFormat: &tag,
	}
	export, _, err := client.DocsDownload(opts)
	if err != nil {
		return err
	}

	d.Set("doc_id", d.Id())
	d.Set("title", export.Title)
	d.Set("revision", export.Revision)
	d.Set("owner", export.Owner)
	return nil
}

func resourceDropboxPaperDocUpdate(d *schema.ResourceData, meta interface{}) error {
	// config := meta.(*ProviderConfig).DropboxConfig
	// client := paper.New(*config)

	// var format *paper.ImportFormat
	// format.Tag = d.Get("import_format").(string)

	// var policy *paper.PaperDocUpdatePolicy
	// policy.Tag = "overwrite_all"

	// opts := &paper.PaperDocUpdateArgs{
	// 	RefPaperDoc:     *paper.NewRefPaperDoc(d.Id()),
	// 	Revision:        d.Get("revision").(int64),
	// 	ImportFormat:    format,
	// 	DocUpdatePolicy: policy,
	// }
	// result, err := client.DocsUpdate(opts, nil)

	// TODO: Implement state usage for file contents changing

	return nil
}

func resourceDropboxPaperDocDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := paper.NewRefPaperDoc(d.Id())
	err := client.DocsArchive(opts)
	return err
}

func convertDocContentToB64() schema.SchemaStateFunc {
	return func(data interface{}) string {
		content := data.(string)
		return base64.StdEncoding.EncodeToString([]byte(content))
	}
}
