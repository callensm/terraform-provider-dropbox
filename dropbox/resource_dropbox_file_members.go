package dropbox

import (
	"fmt"
	"regexp"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDropboxFileMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxFileMembersCreate,
		Read:   resourceDropboxFileMembersRead,
		Update: resourceDropboxFileMembersUpdate,
		Delete: resourceDropboxFileMembersDelete,

		Schema: map[string]*schema.Schema{
			"file_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateWithRegExp(fileIDPattern),
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
			"access_level": {
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
		Members:     createListOfMemberSelectors(d.Get("members").([]string)),
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

	for _, member := range createListOfMemberSelectors(d.Get("members").([]string)) {
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

func createListOfMemberSelectors(members []string) []*sharing.MemberSelector {
	selectors := make([]*sharing.MemberSelector, len(members))
	emailRx := regexp.MustCompile(emailPattern)

	for _, m := range members {
		var s *sharing.MemberSelector
		if emailRx.MatchString(m) {
			s.Tag = "email"
			s.Email = m
		} else {
			s.Tag = "dropbox_id"
			s.DropboxId = m
		}
		selectors = append(selectors, s)
	}
	return selectors
}
