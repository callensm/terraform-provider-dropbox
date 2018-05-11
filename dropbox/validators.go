package dropbox

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

const fileIDPattern = "((/|id:).*|nspath:[0-9]+:.*)|ns:[0-9]+(/.*)?"

const emailPattern = "^['&A-Za-z0-9._%+-]+@[A-Za-z0-9-][A-Za-z0-9.-]*.[A-Za-z]{2,15}$"

const folderPathPattern = "(/(.|[\r\n])*)|(ns:[0-9]+(/.*)?)"

const uploadPathPattern = "(/(.|[\r\n])*)|(ns:[0-9]+(/.*)?)|(id:.*)"

func validateWithRegExp(pattern string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		ok, err := regexp.MatchString(pattern, value)
		if err != nil {
			errors = append(errors, fmt.Errorf("RegExp Validation Failure: %s", err))
		}
		if !ok {
			errors = append(errors, fmt.Errorf("RegExp Validation Failure: %s does not match the pattern %s", value, pattern))
		}
		return
	}
}

func validateDocImportFormat() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		format := v.(string)
		validFormats := []string{"html", "markdown", "plain_text", "other"}

		for _, f := range validFormats {
			if f == format {
				return
			}
		}

		errors = append(errors, fmt.Errorf("Import Format Validation Failure: %s is not a valid import format", format))
		return
	}
}

func validateDocUserPermissionsType() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		level := v.(string)

		if level == "" {
			ws = append(ws, fmt.Sprintf("User Permission Validation Warning: No permission level value was given, resorting to default."))
			return
		}

		if level == "edit" || level == "view_and_comment" {
			return
		}

		errors = append(errors, fmt.Errorf("User Permission Validation Error: %s is not a valid permission level for a Paper doc", level))
		return
	}
}

func validateDocPolicyType(env string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		policy := v.(string)
		approved := []string{"people_with_link_can_edit", "people_with_link_can_view_and_comment", "invite_only"}
		approvedPublic := append(approved, "disabled")

		if env == "public" {
			for _, p := range approvedPublic {
				if policy == p {
					return
				}
			}
			errors = append(errors, fmt.Errorf("Share Policy Validation Failure: %s is not an approved public policy", policy))
		} else {
			for _, p := range approved {
				if policy == p {
					return
				}
			}
			errors = append(errors, fmt.Errorf("Share Policy Validation Failure: %s is not an approved team policy", policy))
		}
		return
	}
}

func validateFileWriteMode() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		mode := v.(string)
		valid := []string{"add", "overwrite", "update"}

		for _, v := range valid {
			if mode == v {
				return
			}
		}

		errors = append(errors, fmt.Errorf("Write Mode Validation Failure: %s is not a valid write mode", mode))
		return
	}
}

func validateFileAccessLevel() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		level := v.(string)
		valid := []string{"owner", "editor", "viewer", "viewer_no_comment"}

		for _, v := range valid {
			if level == v {
				return
			}
		}

		errors = append(errors, fmt.Errorf("Access Level Validation Failure: %s is not a valid file access level", level))
		return
	}
}
