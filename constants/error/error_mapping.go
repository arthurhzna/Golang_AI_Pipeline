package error

func ErrMapping(err error) bool {
	var (
		GeneralErrors = GeneralErrors
	)

	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
