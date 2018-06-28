provider "dropbox" {
  access_token = "5wEd655IJOAAAAAAAAAAigME_-IT3OR9ffdB1WY7cMM8HuRLa5vKWEO8qCtk-lGo"
}

resource "dropbox_folder" "foo" {
  path = "/Terraformed"
}

resource "dropbox_file" "doc" {
  content = "${file("../README.md")}"
  path    = "${dropbox_folder.foo.path}/README.md"
  mode    = "add"
  mute    = true
}

output "folder_id" {
  value = "${dropbox_folder.foo.folder_id}"
}

output "content_hash" {
  value = "${dropbox_file.doc.hash}"
}

output "content_size" {
  value = "${dropbox_file.doc.size}"
}
