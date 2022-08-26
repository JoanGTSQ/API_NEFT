package controllers

func RetrieveErrorAPI(errorCode error) string {
	return "<b>" + errorCode.Error() + ":</b> Please check the inserted data and try again"
}
