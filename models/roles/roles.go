package roles

const (
	Any             = "any"
	Admin           = "admin"
	Read            = "read"
	ReadOne         = "read_one"
	Write           = "write"
	WriteNewOnly    = "write_new_only"
	WriteUpdateOnly = "write_update_only"
)

func CheckRoles(requiredRoles []string, grantedRoles []string) bool {
	for _, requiredRole := range requiredRoles {

		if requiredRole == Any && len(grantedRoles) != 0 {
			return true
		}

		for _, grantedRole := range grantedRoles {
			if grantedRole == Admin {
				return true
			}
			if grantedRole == requiredRole {
				return true
			}
		}
	}
	return false
}
