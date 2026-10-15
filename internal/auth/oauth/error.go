package oauth

import (
	"errors"
)


var StateDoesNotExist = errors.New("The state was not present in the store")