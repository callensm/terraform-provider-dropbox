package dropbox

import (
	"fmt"
	"regexp"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDropboxFolderMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropboxFolderMembersCreate,
		Read:   resourceDropboxFolderMembersRead,
		Update: resourceDropboxFolderMembersUpdate,
		Delete: resourceDropboxFolderMembersDelete,

		Schema: map[string]*schema.Schema{
			"folder_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateWithRegExp(folderIDPattern),
			},
			"members": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "viewer",
							ValidateFunc: validateAccessLevel(),
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
	output := make([]map[string]string, len(m))
	for _, i := range m {
		member := map[string]string{
			"access_level": i.AccessType.Tag,
		}

		if i.User.Email != "" {
			member["identity"] = i.User.Email
		} else if i.User.AccountId != "" {
			member["identity"] = i.User.AccountId
		}

		output = append(output, member)
	}
	return output
}

func createListOfFolderMembers(m []map[string]interface{}) []*sharing.AddMember {
	members := make([]*sharing.AddMember, 0, len(m))
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

		mem := &sharing.AddMember{
			AccessLevel: &sharing.AccessLevel{Tagged: db.Tagged{Tag: i["access_level"].(string)}},
			Member:      &selector,
		}

		members = append(members, mem)
	}
	return members
}

func removeFolderShareMembers(arg *sharing.RemoveFolderMemberArg, client *sharing.Client, members []map[string]string) error {
	emailRx := regexp.MustCompile(emailPattern)
	for _, mem := range members {
		if emailRx.MatchString(mem["identity"]) {
			arg.Member = &sharing.MemberSelector{
				Tagged: db.Tagged{Tag: "email"},
				Email:  mem["identity"],
			}
		} else {
			arg.Member = &sharing.MemberSelector{
				Tagged:    db.Tagged{Tag: "dropbox_id"},
				DropboxId: mem["identity"],
			}
		}
		_, err := (*client).RemoveFolderMember(arg)
		if err != nil {
			return err
		}
	}
	return nil
}
