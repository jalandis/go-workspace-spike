package acl

func HasAccess(s string) bool {
	return s != "john"
}
