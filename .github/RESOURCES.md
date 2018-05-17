# Resource Documentation

## **dropbox_file**

Resource for the management of Dropbox files and entities (not including Paper).

### Example Usage

```hcl
# Upload the new `file.txt` content to the folder `private`
resource "dropbox_file" "file" {
  content = "${file("file.txt")}"
  path    = "/private/file.txt"
  mode    = "add"
  mute    = false
}

output "file_hash" {
  value = "${dropbox_file.file.hash}"
}
```

### Argument Reference

* **content** - (Required) _The file content that is being uploaded to the new file in the specified location. The content is converted into Base64 representation before being stored in state_
* **path** - (Required) _The destination path, including the desired file name, for the content to be uploaded as. Example: `/homework/math/Answers.txt`_
* **mode** - (Optional) _The writing mode of the file uploading. Can be one of the following: `add`, `overwrite`, or `update`. If not value is provided to the resource, it defaults to `add` mode_
* **auto_rename** - (Optional) _Indicates if the folder should automatically be renamed if there is a naming conflict on creation. Defaults to `false`_
* **mute** - (Optional) _Boolean value to determine whether users associated with the file will be notified of its update or creation. Defaults to `false` if not provided, meaning users will be notified_

### Attributes Reference

* **hash** - _A generated hash of the file's content that was uploaded_
* **size** - _The size in bytes of the content that was uploaded as the file_

## **dropbox_file_members**

Resource associated with the management of Dropbox file sharing with specific members of a team or references by account email of ID.

### Example Usage

```hcl
resource "dropbox_file" "test" {
  content       = "${file("myfile.txt")}"
  path          = "/myfile.txt"
  import_format = "plain_text"
}

resource "dropbox_file_members" "mems" {
  file_id = "${dropbox_file.test.id}"
  members = [
    {
      email = "user@example.com"
    }
  ]
  access_level = "editor"
}
```

### Argument Reference

* **file_id** - (Required) _The unique identifier for the Dropbox file you are adding the new members to_
* **members** - (Required) _List of Dropbox users/accounts that are to be give file access. `email` and `account_id` are conflicting fields so only give one per member_
  * **email** - (Optional) _Associated email address for the member being added_
  * **account_id** - (Optional) _The Dropbox account identifier for the member being added_
* **message** - (Optional) _A custom message to be emailed to the added member(s) when they are notified of their new file access_
* **quiet** - (Optional) _Boolean value to specific whether newly added members will be notified of their new file access. Defaults to `false` if no value is given, meaning they will be notified_
* **access_level** - (Optional) _The access level tier that is to be granted to all members in specified in the `members` input. Values can be either: `owner`, `editor`, `viewer`, or `viewer_no_comment`. If no value is given, the default level is `viewer`_

### Attributes Reference

There are no additional attributes that are received from the creation of this resource.

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

## **dropbox_folder_members**

Resource for managing folder sharing access levels among other accounts within a team or publicily.

### Example Usage

```hcl
resource "dropbox_folder" "foo" {
  path = "/source-code"
}

resource "dropbox_folder_members" "mems" {
  folder_id = "${dropbox_folder.foo.id}"
  members = [
    {
      email        = "user@example.com"
      access_level = "edit"
    },
    {
      account_id   = "dbid:4kjh342ljg23fds"
      access_level = "viewer"
    }
  ]
}
```

### Argument Reference

* **folder_id** - (Required) _The unique identifier for the Dropbox folder you are adding the new members to_
* **members** - (Required) _List of Dropbox users/accounts that are to be give file access. `email` and `account_id` are conflicting fields so only give one per member_
  * **email** - (Optional) _Associated email address for the member being added_
  * **account_id** - (Optional) _The Dropbox account identifier for the member being added_
  * **access_level** - (Optional) _The access level tier that is to be granted to all members in specified in the `members` input. Values can be either: `owner`, `editor`, `viewer`, or `viewer_no_comment`. If no value is given, the default level is `viewer`_
* **message** - (Optional) _A custom message to be emailed to the added member(s) when they are notified of their new file access_
* **quiet** - (Optional) _Boolean value to specific whether newly added members will be notified of their new file access. Defaults to `false` if no value is given, meaning they will be notified_

### Attributes Reference

There are no additional attributes that are received from the creation of this resource.

## **dropbox_paper_doc**

Management and file content uploading for the creation of Paper documents in Dropbox. The contents of the argued file are read and stored in state in Base64 format.

### Example Usage

```hcl
# Creates a new Paper document called 'myfile.json' in the /data folder
# containing the contents are the local 'myfile.json' file.
resource "dropbox_paper_doc" "doc" {
  content       = "${file("myfile.json")}"
  parent_folder = "/data"
  import_format = "plain_text"
}
```

### Argument Reference

* **content** - (Required) _The file content that is to be read and uploaded to Dropbox Paper as a new document. The content is converted into Base64 representation before being stored in state_
* **parent_folder** - (Optional) _The folder that the new file should be placed inside of. Defaults to the top level directory (no folder)_
* **import_format** - (Required) \_The format of that data that you are uploading into the document. Value must be of the following: `html`, `markdown`, `plain_text`, or `other`.

### Attributes Reference

* **doc_id** - _The unique identifier for the newly created Paper document_
* **revision** - _An integer value representing the current revision number of the document_
* **title** - _The name/title of the new document_

## **dropbox_paper_doc_users**

Resource for adding accounts to a Dropbox Paper document, or sharing the document with them. Able to specific the permissions level for each user you are sharing with.

### Example Usage

```hcl
resource "dropbox_paper_doc" "doc" {
  content       = "${file("myfile.json")}"
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

Set sharing policies for public and team use for a Dropbox Paper document. The deletion of this resource resets the public and team sharing policies to `invite_only`.

### Example Usage

```hcl
resource "dropbox_paper_doc" "doc" {
  content       = "${file("myfile.json")}"
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
