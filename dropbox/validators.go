package dropbox

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

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
