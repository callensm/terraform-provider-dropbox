package dropbox

import (
	"fmt"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/paper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDropboxPaperSharingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxPaperSharingPolicyCreate,
		Read:   resourceDropboxPaperSharingPolicyRead,
		Update: resourceDropboxPaperSharingPolicyUpdate,
		Delete: resourceDropboxPaperSharingPolicyDelete,

		Schema: map[string]*schema.Schema{
			"doc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDocPolicyType("public"),
			},
			"team_policy": {
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
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := paper.NewRefPaperDoc(d.Get("doc_id").(string))
	policy, err := client.DocsSharingPolicyGet(opts)
	if err != nil {
		return fmt.Errorf("Sharing Policy Read Failure: %s", err)
	}

	d.Set("public_policy", policy.PublicSharingPolicy.Tag)
	d.Set("team_policy", policy.TeamSharingPolicy.Tag)
	return nil
}

func resourceDropboxPaperSharingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	d.Partial(true)
	if d.HasChange("doc_id") || d.HasChange("public_policy") || d.HasChange("team_policy") {
		opts := &paper.PaperDocSharingPolicy{
			RefPaperDoc: *paper.NewRefPaperDoc(d.Get("doc_id").(string)),
			SharingPolicy: &paper.SharingPolicy{
				PublicSharingPolicy: &paper.SharingPublicPolicyType{Tagged: db.Tagged{Tag: d.Get("public_policy").(string)}},
				TeamSharingPolicy:   &paper.SharingTeamPolicyType{Tagged: db.Tagged{Tag: d.Get("team_policy").(string)}},
			},
		}

		err := client.DocsSharingPolicySet(opts)
		if err != nil {
			return fmt.Errorf("Sharing Policy Update Failure: %s", err)
		}
	}
	d.Partial(false)

	return nil
}

func resourceDropboxPaperSharingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := &paper.PaperDocSharingPolicy{
		RefPaperDoc: *paper.NewRefPaperDoc(d.Get("doc_id").(string)),
		SharingPolicy: &paper.SharingPolicy{
			PublicSharingPolicy: &paper.SharingPublicPolicyType{Tagged: db.Tagged{Tag: "invite_only"}},
			TeamSharingPolicy:   &paper.SharingTeamPolicyType{Tagged: db.Tagged{Tag: "invite_only"}},
		},
	}

	err := client.DocsSharingPolicySet(opts)
	if err != nil {
		return fmt.Errorf("Sharing Policy Deletion Failure: %s", err)
	}

	return nil
}
