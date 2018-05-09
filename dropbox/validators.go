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
