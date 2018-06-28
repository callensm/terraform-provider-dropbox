package dropbox

import (
	"fmt"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/paper"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDropboxPaperDocUsers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxPaperDocUserCreate,
		Read:   resourceDropboxPaperDocUserRead,
		Update: resourceDropboxPaperDocUserUpdate,
		Delete: resourceDropboxPaperDocUserDelete,

		Schema: map[string]*schema.Schema{
			"doc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"members": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"members.account_id"},
							ValidateFunc:  validateWithRegExp(emailPattern),
						},
						"account_id": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"members.email"},
						},
						"permissions": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "view_and_comment",
							Description:  "Value must be either `edit` or `view_and_comment`",
							ValidateFunc: validateDocUserPermissionsType(),
						},
					},
				},
			},
			"message": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"quiet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"shared_users": &schema.Schema{
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
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	pid := d.Get("doc_id").(string)

	// TODO: Get members list from active doc not terraform state to get true list to remove from
	members := d.Get("members").([]map[string]string)

	d.Partial(true)
	if d.HasChange("doc_id") {
		err := removeUsersFromDocument(members, pid, &client)
		if err != nil {
			return err
		}
		d.SetPartial("doc_id")
	}

	// TODO:  Find difference between changed members and active and manage accordingly
	if d.HasChange("members") {

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

	for _, member := range createListOfMemberSelectors(d.Get("members").([]map[string]interface{})) {
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
	for _, i := range m {
		var selector sharing.MemberSelector
		if email := i["email"].(string); email != "" {
			selector.Tag = "email"
			selector.Email = email
		} else {
			selector.Tag = "dropbox_id"
			selector.DropboxId = i["account_id"].(string)
		}

		mem := &paper.AddMember{
			PermissionLevel: &paper.PaperDocPermissionLevel{Tagged: db.Tagged{Tag: i["permissions"].(string)}},
			Member:          &selector,
		}

		members = append(members, mem)
	}
	return members
}

func removeUsersFromDocument(members []map[string]string, id string, client *paper.Client) error {
	for _, m := range members {
		opts := &paper.RemovePaperDocUser{
			RefPaperDoc: *paper.NewRefPaperDoc(id),
		}

		if m["email"] != "" {
			opts.Member = &sharing.MemberSelector{
				Tagged: db.Tagged{Tag: "email"},
				Email:  m["email"],
			}
		} else {
			opts.Member = &sharing.MemberSelector{
				Tagged:    db.Tagged{Tag: "dropbox_id"},
				DropboxId: m["account_id"],
			}
		}

		err := (*client).DocsUsersRemove(opts)
		if err != nil {
			return fmt.Errorf("Doc Users Update Failure: Couldn't remove user %+v from document %s", opts.Member, id)
		}
	}

	return nil
}
