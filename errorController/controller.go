package errorController

import (
	"errors"
)

func RetrieveError(errorCode string) error {
	return errors.New(errorCode)
}

func RetrieveErrorAPI(errorCode string) error {
	return errors.New("<b>" + errorCode + ":</b> Please check the inserted data and try again")
}
