package errorController

import (
	"errors"
)

const (
	ERR_NOT_MATCH_LOGIN = "ALFA-001"

	ERR_NOT_FOUND       = "ALFA-002"
	ERR_USERID_REQUIRED = "ALFA-003"
	ERR_ID_INVALID      = "ALFA-004"

	ERR_MAIL_REQUIRED     = "ALFA-005"
	ERR_MAIL_IS_N0T_VALID = "ALFA-006"
	ERR_MAIL_IS_TAKEN     = "ALFA-007"
	ERR_MAIL_INVALID      = "ALFA-008"
	ERR_MAIL_NOT_EXIST    = "ALFA-009"
	ERR_NOT_FOUND_MODEL   = "models: resource not found"

	ERR_PSSWD_INCORRECT  = "ALFA-010"
	ERR_PSSWD_TOO_SHORT  = "ALFA-011"
	ERR_PSSWD_REQUIRED   = "ALFA-012"
	ERR_PSSWD_SAME_RESET = "ALFA-003"

	ERR_REMMEMBER_TOO_SHOT = "ALFA-014"
	ERR_REMMEMBER_REQUIRED = "ALFA-015"
)

func RetrieveError(errorCode string) error {
	return errors.New("<b>" + errorCode + ":</b> Please check the inserted data and try again")
}
