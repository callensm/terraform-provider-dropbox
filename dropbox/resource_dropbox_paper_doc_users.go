package dropbox

import (
	"fmt"
	"regexp"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/paper"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDropboxPaperDocUsers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxPaperDocUserCreate,
		Read:   resourceDropboxPaperDocUserRead,
		Update: resourceDropboxPaperDocUserUpdate,
		Delete: resourceDropboxPaperDocUserDelete,

		Schema: map[string]*schema.Schema{
			"doc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"members": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"permissions": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "view_and_comment",
							Description:  "Value must be either `edit` or `view_and_comment`",
							ValidateFunc: validateDocUserPermissionsType(),
						},
					},
				},
			},
			"message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quiet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"shared_users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Computed list of emails of those users invited and actively sharing the document",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDropboxPaperDocUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := &paper.AddPaperDocUser{
		RefPaperDoc: *paper.NewRefPaperDoc(d.Get("doc_id").(string)),
		Members:     createListOfAddMembers(d.Get("members").([]map[string]interface{})),
		Quiet:       d.Get("quiet").(bool),
	}

	if msg, ok := d.GetOk("message"); ok {
		opts.CustomMessage = msg.(string)
	}

	statuses, err := client.DocsUsersAdd(opts)
	if err != nil {
		return fmt.Errorf("Doc Users Failure: %s", err)
	}

	for _, s := range statuses {
		tag := s.Result.Tag
		if tag != "success" && tag != "user_is_owner" && tag != "permission_already_granted" {
			var id string
			if s.Member.Email != "" {
				id = s.Member.Email
			} else {
				id = s.Member.DropboxId
			}
			return fmt.Errorf("Doc Users Failure: User %s returned status %s", id, tag)
		}
	}

	d.SetId(fmt.Sprintf("%s:%d", d.Get("doc_id").(string), len(opts.Members)))
	return resourceDropboxPaperDocUserRead(d, meta)
}

func resourceDropboxPaperDocUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := paper.NewListUsersOnPaperDocArgs(d.Get("doc_id").(string))
	response, err := client.DocsUsersList(opts)
	if err != nil {
		return fmt.Errorf("Doc Users Failure: %s", err)
	}

	invitees, shared := response.Invitees, response.Users
	emails := make([]string, 0, len(invitees)+len(shared))
	for _, s := range shared {
		emails = append(emails, s.User.Email)
	}
	for _, i := range invitees {
		emails = append(emails, i.Invitee.Email)
	}

	d.Set("shared_users", emails)

	return nil
}

func resourceDropboxPaperDocUserUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("members") {
		return resourceDropboxPaperDocUserCreate(d, meta)
	}
	d.Partial(false)

	return nil
}

func resourceDropboxPaperDocUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := &paper.RemovePaperDocUser{
		RefPaperDoc: *paper.NewRefPaperDoc(d.Get("doc_id").(string)),
	}

	membersList := d.Get("members").([]map[string]interface{})
	ids := make([]string, len(membersList))
	for i, m := range membersList {
		ids[i] = m["identity"].(string)
	}

	for _, member := range createListOfMemberSelectors(ids) {
		opts.Member = member
		err := client.DocsUsersRemove(opts)
		if err != nil {
			return fmt.Errorf("Doc Users Failure: %s", err)
		}
	}

	return nil
}

func createListOfAddMembers(m []map[string]interface{}) []*paper.AddMember {
	members := make([]*paper.AddMember, 0, len(m))
	emailRx := regexp.MustCompile(emailPattern)

	for _, i := range m {
		var selector sharing.MemberSelector
		identity := i["identity"].(string)

		if emailRx.MatchString(identity) {
			selector.Tag = "email"
			selector.Email = identity
		} else {
			selector.Tag = "dropbox_id"
			selector.DropboxId = identity
		}

		mem := &paper.AddMember{
			PermissionLevel: &paper.PaperDocPermissionLevel{Tagged: db.Tagged{Tag: i["permissions"].(string)}},
			Member:          &selector,
		}

		members = append(members, mem)
	}
	return members
}
