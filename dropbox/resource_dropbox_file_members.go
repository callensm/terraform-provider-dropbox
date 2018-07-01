package dropbox

import (
	"fmt"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDropboxFileMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxFileMembersCreate,
		Read:   resourceDropboxFileMembersRead,
		Update: resourceDropboxFileMembersUpdate,
		Delete: resourceDropboxFileMembersDelete,

		Schema: map[string]*schema.Schema{
			"file_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateWithRegExp(fileIDPattern),
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
			"access_level": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "viewer",
				ValidateFunc: validateAccessLevel(),
			},
		},
	}
}

func resourceDropboxFileMembersCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := sharing.New(*config)

	opts := &sharing.AddFileMemberArgs{
		File:        d.Get("file_id").(string),
		Members:     createListOfMemberSelectors(d.Get("members").([]map[string]interface{})),
		Quiet:       d.Get("quiet").(bool),
		AccessLevel: &sharing.AccessLevel{Tagged: db.Tagged{Tag: d.Get("access_level").(string)}},
	}

	if msg, ok := d.GetOk("message"); ok {
		opts.CustomMessage = msg.(string)
	}

	res, err := client.AddFileMember(opts)
	if err != nil {
		return fmt.Errorf("File Member Creation Failure: %s", err)
	}

	for _, r := range res {
		if r.Result.MemberError != nil {
			return fmt.Errorf("File Member Creation Failure: %+v", r.Result.MemberError)
		}
	}

	d.SetId(fmt.Sprintf("%s:%d", d.Get("file_id").(string), len(opts.Members)))
	return nil
}

func resourceDropboxFileMembersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := sharing.New(*config)

	opts := sharing.NewListFileMembersArg(d.Get("file_id").(string))
	fileMembers, err := client.ListFileMembers(opts)
	if err != nil {
		return fmt.Errorf("File Member Read Failure: %s", err)
	}

	var members []string
	for _, m := range fileMembers.Users {
		if email := m.User.Email; email != "" {
			members = append(members, email)
		} else {
			members = append(members, m.User.AccountId)
		}
	}

	d.Set("members", members)
	return nil
}

func resourceDropboxFileMembersUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("memebrs") {
		return resourceDropboxFileMembersCreate(d, meta)
	}
	d.Partial(false)

	return nil
}

func resourceDropboxFileMembersDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := sharing.New(*config)

	opts := &sharing.RemoveFileMemberArg{
		File: d.Get("file_id").(string),
	}

	for _, member := range createListOfMemberSelectors(d.Get("members").([]map[string]interface{})) {
		opts.Member = member
		res, err := client.RemoveFileMember2(opts)
		if err != nil {
			return fmt.Errorf("File Member Deletion Failure: %s", err)
		}
		if res.MemberError != nil {
			return fmt.Errorf("File Member Deletion Failure: %+v", res.MemberError)
		}
	}

	return nil
}

func createListOfMemberSelectors(m []map[string]interface{}) []*sharing.MemberSelector {
	members := make([]*sharing.MemberSelector, 0, len(m))
	for _, i := range m {
		var selector sharing.MemberSelector
		if email := i["email"].(string); email != "" {
			selector.Tag = "email"
			selector.Email = email
		} else {
			selector.Tag = "dropbox_id"
			selector.DropboxId = i["account_id"].(string)
		}

		members = append(members, &selector)
	}
	return members
}
