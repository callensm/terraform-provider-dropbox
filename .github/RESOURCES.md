# Resource Documentation

## **dropbox_folder**

Management of folders within your Dropbox account (not including Dropbox Paper).

### Example Usage

```hcl
# Create a folder called 'code' at the top level directory
resource "dropbox_folder" "folder_a" {
  path        = "/code"
  auto_rename = false
}

# Create a nested folder within the previously make folder
resource "dropbox_folder" "folder_b" {
  path        = "${dropbox_folder.folder_a.path}/source"
  auto_rename = false
}

# Creates folder system /code/source
```

### Argument Reference

* **path** - (Required) _The directory path to create and name the new folder. An empty string value indicates the top level directory_
* **auto_rename** - (Optional) _Indicates if the folder should automatically be renamed if there is a naming conflict on creation. Defaults to `false`_

### Attributes Reference

* **folder_id** - _The unique identifier for the newly created folder_
* **name** - _The display name of the new folder_
* **property_group_templates** - _A list of template IDs for the property groups that are assigned to the new folder. May be empty_

## **dropbox_paper_doc**

Management and file content uploading for the creation of Paper documents in Dropbox. The contents of the argued file are read and stored in state in Base64 format.

### Example Usage

```hcl
# Creates a new Paper document called 'myfile.json' in the /data folder
# containing the contents are the local 'myfile.json' file.
resource "dropbox_paper_doc" "doc" {
  content_file  = "${file("myfile.json")}"
  parent_folder = "/data"
  import_format = "plain_text"
}
```

### Argument Reference

* **content_file** - (Required) _The file that is to be read and uploaded to Dropbox Paper as a new document_
* **parent_folder** - (Optional) _The folder that the new file should be placed inside of. Defaults to the top level directory (no folder)_
* **import_format** - (Required) \_The format of that data that you are uploading into the document. Value must be of the following: `html`, `markdown`, `plain_text`, or `other`.

### Attributes Reference

* **doc_id** - _The unique identifier for the newly created Paper document_
* **revision** - _An integer value representing the current revision number of the document_
* **title** - _The name/title of the new document_
* **owner** - _The name or email address of the owner account for the new document_

## **dropbox_paper_doc_users**

Resource for adding accounts to a Dropbox Paper document, or sharing the document with them. Able to specific the permissions level for each user you are sharing with.

### Example Usage

```hcl
resource "dropbox_paper_doc" "doc" {
  content_file  = "${file("myfile.json")}"
  import_format = "plain_text"
}

resource "dropbox_paper_doc_users" "doc_users" {
  doc_id  = "${dropbox_paper_doc.doc.doc_id}"
  members = [
    {
      email       = "user1@example.com"
      permissions = "view_and_comment"
    },
    {
      account_id  = "dbid:adfijrogqhofqer"
      permissions = "edit"
    }
  ]
  message = "You have been invited to share my new document!"
}
```

### Argument Reference

* **doc_id** - (Required) _The unique identifier for the Paper document you are adding the share users to_
* **members** - (Required) _A list of member selectors for the accounts you wish to share the document with and assign permissions. The `email` and `account_id` are conflicting fields so only enter one of the two in order to specific the user to be added._
  * **email** - (Optional) _The email address of the user's account to add_
  * **account_id** - (Optional) _The unique account identifier for the account to add_
  * **permissions** - (Optional) \_The permission level of the user for their access to the document. Must be `view_and_comment` or `edit`. Defaults to `view_and_comment`.
* **message** - (Optional) _A custom message to be emailed to the newly added users in their invitation to share the document_
* **quiet** - (Optional) _Boolean to specific if the users should receive an email about their invitation to share or not. Defaults to false (they will be enabled)_

### Attributes Reference

* **shared_users** - _List of email addresses of the users who are currently sharing or have been invited to share the document_

## **dropbox_paper_sharing_policy**

Set sharing policies for public and team use for a Dropbox Paper document.

### Example Usage

```hcl
resource "dropbox_paper_doc" "doc" {
  content_file  = "${file("myfile.json")}"
  import_format = "plain_text"
}

# Only allows invite access to the public
# and only view/comment permissions to team members with document link
resource "dropbox_paper_sharing_policy" "policies" {
  doc_id        = "${dropbox_paper_doc.doc.doc_id}"
  public_policy = "invite_only"
  team_policy   = "people_with_link_can_view_and_comment"
}
```

### Argument Reference

* **doc_id** - (Required) _The unique identifier for the Paper document you are assigning the new sharing policies to_
* **public_policy** - (Optional) _The policy to assign the document for public sharing use. Must be value of the following: `people_with_link_can_edit`, `people_with_link_can_view_and_comment`, `invite_only`, or `disabled`_
* **team_policy** - (Optional) _The policy to assign the document for team sharing use. Must be value of the following: `people_with_link_can_edit`, `people_with_link_can_view_and_comment`, or `invite_only`_

### Attributes Reference

No additional attributes are created or returned from the creation of this resource.
