package dropbox

import (
	"fmt"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDropboxFolderMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxFolderMembersCreate,
		Read:   resourceDropboxFolderMembersRead,
		Update: resourceDropboxFolderMembersUpdate,
		Delete: resourceDropboxFolderMembersDelete,

		Schema: map[string]*schema.Schema{
			"folder_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateWithRegExp(folderIDPattern),
			},
			"members": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
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
						"access_level": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "viewer",
							ValidateFunc: validateAccessLevel(),
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
		},
	}
}

func resourceDropboxFolderMembersCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := sharing.New(*config)

	opts := &sharing.AddFolderMemberArg{
		SharedFolderId: d.Get("folder_id").(string),
		Members:        createListOfFolderMembers(d.Get("members").([]map[string]interface{})),
		Quiet:          d.Get("quiet").(bool),
	}

	if msg, ok := d.GetOk("message"); ok {
		opts.CustomMessage = msg.(string)
	}

	err := client.AddFolderMember(opts)
	if err != nil {
		return fmt.Errorf("Folder Member Creation Failure: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%d", d.Get("folder_id").(string), len(opts.Members)))
	return nil
}

func resourceDropboxFolderMembersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := sharing.New(*config)

	opts := &sharing.ListFolderMembersArgs{
		ListFolderMembersCursorArg: *sharing.NewListFolderMembersCursorArg(),
		SharedFolderId:             d.Get("folder_id").(string),
	}
	members, err := client.ListFolderMembers(opts)
	if err != nil {
		return fmt.Errorf("Folder Member Read Failure: %s", err)
	}

	d.Set("members", createListOfTerraformMembers(members.Users))
	return nil
}

func resourceDropboxFolderMembersUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceDropboxFolderMembersCreate(d, meta)
}

func resourceDropboxFolderMembersDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := sharing.New(*config)

	opts := &sharing.RemoveFolderMemberArg{
		SharedFolderId: d.Get("folder_id").(string),
		LeaveACopy:     false,
	}
	err := removeFolderShareMembers(opts, &client, d.Get("members").([]map[string]string))
	if err != nil {
		return fmt.Errorf("Folder Member Deletion Failure: %s", err)
	}

	return nil
}

func createListOfTerraformMembers(m []*sharing.UserMembershipInfo) []map[string]string {
	output := make([]map[string]string, 0, len(m))
	for _, i := range m {
		member := make(map[string]string)
		member["email"] = i.User.Email
		member["account_id"] = i.User.AccountId
		member["access_level"] = i.AccessType.Tag
		output = append(output, member)
	}
	return output
}

func createListOfFolderMembers(m []map[string]interface{}) []*sharing.AddMember {
	members := make([]*sharing.AddMember, 0, len(m))
	for _, i := range m {
		var selector sharing.MemberSelector
		if email := i["email"]; email != "" {
			selector.Tag = "email"
			selector.Email = email.(string)
		} else {
			selector.Tag = "dropbox_id"
			selector.DropboxId = i["account_id"].(string)
		}

		mem := &sharing.AddMember{
			AccessLevel: &sharing.AccessLevel{Tagged: db.Tagged{Tag: i["access_level"].(string)}},
			Member:      &selector,
		}

		members = append(members, mem)
	}
	return members
}

func removeFolderShareMembers(arg *sharing.RemoveFolderMemberArg, client *sharing.Client, members []map[string]string) error {
	for _, mem := range members {
		if mem["email"] != "" {
			arg.Member = &sharing.MemberSelector{
				Tagged: db.Tagged{Tag: "email"},
				Email:  mem["email"],
			}
		} else {
			arg.Member = &sharing.MemberSelector{
				Tagged:    db.Tagged{Tag: "dropbox_id"},
				DropboxId: mem["account_id"],
			}
		}
		_, err := (*client).RemoveFolderMember(arg)
		if err != nil {
			return err
		}
	}
	return nil
}
