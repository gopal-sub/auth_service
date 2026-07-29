package httputils

import (
	"encoding/json"
	"errors"
)

var JsonSyntaxError = errors.New("Syntax error in JSON")
var DisallowedFieldsError = errors.New("JSON has disallowed fields")





var SyntaxError *json.SyntaxError
