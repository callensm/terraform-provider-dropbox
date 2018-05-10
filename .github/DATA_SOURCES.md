# Data Source Documentation

## **dropbox_user_current**

Used to access information about the current accessing user account.

### Example Usage

```hcl
data "dropbox_user_current" "curr" {}

output "user_id" {
  value = "${data.dropbox_user_current.curr.account_id}"
}
```

### Argument Reference

There are no arguments for this data source

### Attributes Reference

* **account_id** - _The unique Dropbox ID for the active account_
* **account_type** - _The account type can be either `basic`, `business`, or `pro`_
* **display_name** - _Concatenation of the first and last name associated with the account_
* **email** - _The email associated with the account_

## **dropbox_space_usage**

Used to get information about the space used and allocated for the current user.

### Example Usage

```hcl
data "dropbox_space_usage" "usage" {}

output "space_used" {
  value = "${data.dropbox_space_usage.usage.used}"
}
```

### Argument Reference

There are no arguments for this data source

### Attributes Reference

* **used** - _The amount of bytes used to date for the account or team_
* **allocated** - _The amount of space allocated for the user or team in bytes_
* **is_team_allocation** - _True/False of whether the space allocation is shared with a team or not_

## **dropbox_paper_folder**

Used to access information about the folder or folders that are associated with a given Paper document.

### Example Usage

```hcl
data "dropbox_paper_folder" "target" {
  doc_id = "jGaf45sdG35aF"
}

output "num_folders" {
  value = "${len(data.dropbox_paper_folder.target.folders)}"
}
```

### Argument Reference

* **doc_id** - _The unique ID for the Paper document you are search for_

### Attribute Reference

* **folders** - _List of folders associated with Paper document. May be empty_
  * **id** - _Unique ID for the folder instance at an index_
  * **name** - _The display name of the folder instance at an index_

## **dropbox_paper_sharing_policy**

Used to get the public and team sharing policies associated with a target Paper document.

### Example Usage

```hcl
data "dropbox_paper_sharing_policy" "doc" {
  doc_id = "jGaf45sdG35aF"
}

output "doc_public" {
  value = "${data.dropbox_paper_sharing_policy.doc.public_policy}"
}
```

### Argument Reference

* **doc_id** - _The unique ID for the Paper document you are search for_

### Attribute Reference

* **public_policy** - _The policy attached to the document for public use_
* **team_policy** - _The policy attached to the document for team use_
