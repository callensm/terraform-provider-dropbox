package dropbox

import (
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
				ValidateFunc: validateFileAccessLevel(),
			},
		},
	}
}

func resourceDropboxFileMembersCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDropboxFileMembersRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDropboxFileMembersUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDropboxFileMembersDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
