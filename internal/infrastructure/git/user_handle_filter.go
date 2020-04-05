package git

func userIn(user string, allowedUsersHandles []string) bool {
	if allowedUsersHandles == nil {
		return true
	}

	for _, allowedUsersHandle := range allowedUsersHandles {
		if user == allowedUsersHandle {
			return true
		}
	}

	return false
}
