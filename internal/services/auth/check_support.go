package auth

import "net/http"

func checkOK() CheckOutput {
	return CheckOutput{
		StatusCode: http.StatusOK,
		Allowed:    true,
	}
}

func checkForbidden() CheckOutput {
	return CheckOutput{
		StatusCode: http.StatusForbidden,
		Allowed:    false,
	}
}

func checkUnauthorized(message string) CheckOutput {
	return CheckOutput{
		StatusCode: http.StatusUnauthorized,
		Allowed:    false,
		Message:    message,
	}
}

func checkInternalServerError(err error) CheckOutput {
	return CheckOutput{
		StatusCode: http.StatusInternalServerError,
		Allowed:    false,
		Message:    err.Error(),
	}
}

func checkNotFound() CheckOutput {
	return CheckOutput{
		StatusCode: http.StatusNotFound,
		Allowed:    false,
	}
}

func checkMethodNotAllowed() CheckOutput {
	return CheckOutput{
		StatusCode: http.StatusMethodNotAllowed,
		Allowed:    false,
	}
}
