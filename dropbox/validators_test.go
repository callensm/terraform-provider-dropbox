package dropbox

import "testing"

func TestValidateRegExp(t *testing.T) {
	validIDs := map[string]string{
		"id:3kmLmQFnf1AAAAAAAAAAAw": fileIDPattern,
	}

	for key, val := range validIDs {
		_, errors := validateWithRegExp(val)(key, "test")
		if len(errors) != 0 {
			t.Fatalf("RegExp Validator Failure: %s produced an output with errors", key)
		}
	}
}

func TestDocImportFormats(t *testing.T) {
	validFormats := []string{"html", "markdown", "plain_text", "other"}

	for _, format := range validFormats {
		_, errors := validateDocImportFormat()(format, "test")
		if len(errors) != 0 {
			t.Fatalf("Import Format Validator Failure: %s produced an output with errors", format)
		}
	}
}

func TestPermissionTypes(t *testing.T) {
	validTypes := []string{"edit", "view_and_comment"}

	for _, perm := range validTypes {
		_, errors := validateDocUserPermissionsType()(perm, "test")
		if len(errors) != 0 {
			t.Fatalf("Permission Type Validator Failure: %s produced an output with errors", perm)
		}
	}
}

func TestDocPolicyType(t *testing.T) {
	validPolicies := []string{"people_with_link_can_edit", "people_with_link_can_view_and_comment", "invite_only"}

	for _, policy := range validPolicies {
		_, errors := validateDocPolicyType("team")(policy, "test")
		if len(errors) != 0 {
			t.Fatalf("Doc Policy Validator Failure: %s produced an output with errors", policy)
		}
	}

	_, errors := validateDocPolicyType("team")("disabled", "test")
	if len(errors) == 0 {
		t.Fatalf("Doc Policy Validator Failure: there was no error for passing `disabled` to `team` policy type")
	}
}

func TestFileWriteMode(t *testing.T) {
	validModes := []string{"add", "overwrite", "update"}

	for _, mode := range validModes {
		_, errors := validateFileWriteMode()(mode, "test")
		if len(errors) != 0 {
			t.Fatalf("File Write Mode Validator Failure: %s produced an output with errors", mode)
		}
	}
}

func TestAccessLevels(t *testing.T) {
	validLevels := []string{"owner", "editor", "viewer", "viewer_no_comment"}

	for _, level := range validLevels {
		_, errors := validateAccessLevel()(level, "test")
		if len(errors) != 0 {
			t.Fatalf("Access Level Validator Failure: %s produced an output with errors", level)
		}
	}
}
