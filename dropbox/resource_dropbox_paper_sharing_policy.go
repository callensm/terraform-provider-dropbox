package dropbox

import (
	"fmt"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/paper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDropboxPaperSharingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxPaperSharingPolicyCreate,
		Read:   resourceDropboxPaperSharingPolicyRead,
		Update: resourceDropboxPaperSharingPolicyUpdate,
		Delete: resourceDropboxPaperSharingPolicyDelete,

		Schema: map[string]*schema.Schema{
			"doc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public_policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDocPolicyType("public"),
			},
			"team_policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDocPolicyType("team"),
			},
		},
	}
}

func resourceDropboxPaperSharingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := &paper.PaperDocSharingPolicy{
		RefPaperDoc: *paper.NewRefPaperDoc(d.Get("doc_id").(string)),
		SharingPolicy: &paper.SharingPolicy{
			PublicSharingPolicy: &paper.SharingPublicPolicyType{Tagged: db.Tagged{Tag: d.Get("public_policy").(string)}},
			TeamSharingPolicy:   &paper.SharingTeamPolicyType{Tagged: db.Tagged{Tag: d.Get("team_policy").(string)}},
		},
	}

	err := client.DocsSharingPolicySet(opts)
	if err != nil {
		return fmt.Errorf("Sharing Policy Creation Failure: %s", err)
	}

	d.SetId(fmt.Sprintf("pset:%s", d.Get("doc_id").(string)))

	return nil
}

func resourceDropboxPaperSharingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDropboxPaperSharingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDropboxPaperSharingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
