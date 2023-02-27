package validate

func User(user int64) error {
	if user == 0 {
		return ErrEmptyUser
	}
	return nil
}
