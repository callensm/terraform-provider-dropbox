package dropbox

import (
	"encoding/base64"
	"fmt"
	"strings"

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
			"content": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: convertContentToB64(),
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
		},
	}
}

func resourceDropboxPaperDocCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	content := d.Get("content").(string)
	reader := strings.NewReader(content)

	opts := &paper.PaperDocCreateArgs{
		ImportFormat:   &paper.ImportFormat{Tagged: db.Tagged{Tag: d.Get("import_format").(string)}},
		ParentFolderId: d.Get("parent_folder").(string),
	}

	results, err := client.DocsCreate(opts, reader)
	if err != nil {
		return fmt.Errorf("Doc Creation Failure: %s", err)
	}

	d.SetId(results.DocId)
	d.Set("doc_id", results.DocId)
	d.Set("title", results.Title)
	d.Set("revision", results.Revision)
	return nil
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
		return fmt.Errorf("Doc Read Failure: %s", err)
	}

	d.Set("doc_id", d.Id())
	d.Set("title", export.Title)
	d.Set("revision", export.Revision)
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

func convertContentToB64() schema.SchemaStateFunc {
	return func(data interface{}) string {
		content := data.(string)
		return base64.StdEncoding.EncodeToString([]byte(content))
	}
}
