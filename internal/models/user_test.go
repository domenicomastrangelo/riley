package models

import (
	"testing"

	"riley/internal/config"
	"riley/internal/sql"
)

func TestUserCreate(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "testemailsignup@example.com"
	password := "password123%A%"

	user, err := UserCreate(email, password, db)
	if err != nil {
		t.Error("Testing signup: Wanted nil, got", err)
	}

	if !user.Exists(db) {
		t.Error("Testing check user exists: Wanted true, got false")
	}

	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Error("Delete user during signup test failed: Wanted nil, got", err)
		}
	}()

	// Test signup with existing email
	user2, err := UserCreate(email, password, db)
	if err == nil {
		t.Error("Testing signup with existing email: Wanted error, got nil")

		err = user2.Delete(false, db)
		if err != nil {
			t.Error("Delete user during signup test failed: Wanted nil, got", err)
		}
	}
}

func TestSignupWithAllSpacesPassword(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "email@exampleallspaces.com"
	password := "        p3%A"

	user, err := UserCreate(email, password, db)
	if err == nil {
		t.Error("Testing signup with all spaces password: Wanted error, got nil")

		if user.Exists(db) {
			t.Error("Testing check user exists: Wanted false, got true")
		}

		err = user.Delete(false, db)
		if err != nil {
			t.Error("Delete user during signup test failed: Wanted nil, got", err)
		}
	}
}

func TestUserCreateWithWrongEmailNoDomain(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "testemailsignupwithwrongemail"
	password := "password"

	user, err := UserCreate(email, password, db)
	if err == nil {
		t.Error("Testing signup with wrong email: Wanted error, got nil")

		if user.Exists(db) {
			t.Error("Testing check user exists: Wanted false, got true")
		}

		err = user.Delete(false, db)
		if err != nil {
			t.Error("Delete user during signup test failed: Wanted nil, got", err)
		}
	}
}

func TestUserCreateWithWrongEmailNoAt(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "testemailsignupwithwrongemail.com"
	password := "password"

	user, err := UserCreate(email, password, db)
	if err == nil {
		t.Error("Testing signup with wrong email: Wanted error, got nil")

		if user.Exists(db) {
			t.Error("Testing check user exists: Wanted false, got true")
		}

		err = user.Delete(false, db)
		if err != nil {
			t.Error("Delete user during signup test failed: Wanted nil, got", err)
		}
	}
}

func TestUserCreateWithWrongPassword(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "testuserCreatewrongpass@example.com"
	passwords := []string{
		"pass",
		"passqo",
		"pass123",
		"pas$123",
	}

	for _, password := range passwords {
		user, err := UserCreate(email, password, db)
		if err == nil {
			t.Error("Testing signup with wrong password: Wanted error, got nil")

			if user.Exists(db) {
				t.Error("Testing check user exists: Wanted false, got true")
			}

			err = user.Delete(false, db)
			if err != nil {
				t.Error("Delete user during signup test failed: Wanted nil, got", err)
			}
		}
	}
}

func TestUserCreateWithCorrectPassword(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "testwithcorrectpass@example.com"

	passwords := []string{
		"password123$A",
		"pa$$1234AA",
		"pa.'-s00A",
		"pa$$woArd1234567890$üäö:_^°!\"§$%&/()=?`´*+~#'-.,;<>|",
	}

	for _, password := range passwords {
		user, err := UserCreate(email, password, db)
		if err != nil {
			t.Error("Testing signup with correct password: Wanted nil, got", err)
		}

		if !user.Exists(db) {
			t.Error("Testing check user exists: Wanted true, got false")
		}

		err = user.Delete(false, db)
		if err != nil {
			t.Error("Delete user during signup test failed: Wanted nil, got", err)
		}
	}
}

func TestUserDeleteSoft(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "testuserdeletesoft@example.com"
	password := "password123!§$%AA"

	user, err := UserCreate(email, password, db)
	if err != nil {
		t.Fatal("Signup during delete soft test failed: Wanted nil, got", err)
	}

	err = user.Delete(true, db)
	if err != nil {
		t.Error("Testing delete soft: Wanted nil, got", err)
	}

	if user.Exists(db) {
		t.Error("Testing check user exists: Wanted false, got true")
	}

	if !user.IsSoftDeleted(db) {
		t.Error("Testing check soft deleted: Wanted true, got false")
	}

	err = user.Delete(false, db)
	if err != nil {
		t.Error("Delete user during delete soft test failed: Wanted nil, got", err)
	}
}

func TestUserDelete(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "testemaildelete@example.com"
	password := "password123$$AA"

	user, err := UserCreate(email, password, db)
	if err != nil {
		t.Error("Signup during delete test failed: Wanted nil, got", err)
	}

	err = user.Delete(false, db)
	if err != nil {
		t.Error("Testing delete: Wanted nil, got", err)
	}
}

func TestUserCheckLogin(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	email := "wrong@example.com"
	password := "password123$$AA"

	userID, err := UserCheckLogin(email, password, db)
	if err != nil {
		t.Error("Testing check login with wrong email and password: Wanted nil, got", err)
	}

	if userID != 0 {
		t.Error("Testing check login with wrong email and password: Wanted 0, got", userID)
	}

	userID, err = UserCheckLogin("", password, db)
	if err != nil {
		t.Error("Testing check login with empty email: Wanted nil, got", err)
	}

	if userID != 0 {
		t.Error("Testing check login with empty email: Wanted 0, got", userID)
	}

	userID, err = UserCheckLogin(email, "", db)
	if err != nil {
		t.Error("Testing check login with empty password: Wanted nil, got", err)
	}

	if userID != 0 {
		t.Error("Testing check login with empty password: Wanted 0, got", userID)
	}

	userID, err = UserCheckLogin("", "", db)
	if err != nil {
		t.Error("Testing check login with empty email and password: Wanted nil, got", err)
	}

	if userID != 0 {
		t.Error("Testing check login with empty email and password: Wanted 0, got", userID)
	}

	email = "correct@example.com"
	password = "correctpasswordpassword123$$AA"

	user, err := UserCreate(email, password, db)
	if err != nil {
		t.Error("Signup during login test failed: Wanted nil, got", err)
	}

	if userID, err := UserCheckLogin(email, password, db); userID == 0 && err != nil {
		t.Error("Testing check login with correct email and password: Wanted true, got false")
	}

	err = user.Delete(false, db)
	if err != nil {
		t.Error("Delete user during login test failed: Wanted nil, got", err)
	}
}
