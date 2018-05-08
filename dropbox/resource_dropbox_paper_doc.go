package dropbox

import (
	"fmt"
	"os"

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
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_folder": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"import_format": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Valid formats include: html, markdown, plain_text, other",
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
		tag.Tag = paper.ImportFormatHtml
	case "markdown":
		tag.Tag = paper.ImportFormatMarkdown
	case "plain_text":
		tag.Tag = paper.ImportFormatPlainText
	case "other":
		tag.Tag = paper.ImportFormatOther
	default:
		return fmt.Errorf("Invalid import format given for paper document creation: %s", format)
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

	// opts := &paper.PaperDocUpdateArgs{
	// 	RefPaperDoc:  *paper.NewRefPaperDoc(d.Id()),
	// 	Revision:     d.Get("revision").(int64),
	// 	ImportFormat: d.Get("import_format").(*paper.ImportFormat),
	// }

	// TODO: Figure out updating Paper documents with different types

	return nil
}

func resourceDropboxPaperDocDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := paper.NewRefPaperDoc(d.Id())
	err := client.DocsArchive(opts)
	return err
}
