package models

import (
	"regexp"
)

func UserCreateValidation(email string, password string) bool {
	return UserEmailIsValid(email) && UserPasswordIsValid(password)
}

func UserEmailIsValid(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return re.MatchString(email)
}

func UserPasswordIsValid(password string) bool {
	// Password must be at least 8 characters long
	re := regexp.MustCompile(`^.{8,}$`)

	// Password must be at most 64 characters long
	if len(password) > 64 {
		return false
	}

	// Password must contain at least one uppercase letter and
	// Password must contain at least one lowercase letter and
	// Password must contain at least one number and
	// Password must contain at least one special character
	reUpper := regexp.MustCompile(`[A-Z]`)
	reLower := regexp.MustCompile(`[a-z]`)
	reNumber := regexp.MustCompile(`[0-9]`)
	reSpecial := regexp.MustCompile(`[!@#~$%^&*()\-_=+[\]{}|\\;:'",.<>/?]+`)

	return re.MatchString(password) &&
		reUpper.MatchString(password) &&
		reLower.MatchString(password) &&
		reNumber.MatchString(password) &&
		reSpecial.MatchString(password)
}
