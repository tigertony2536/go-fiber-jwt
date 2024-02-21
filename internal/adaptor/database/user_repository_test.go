package database_test

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

func init() {
	err := godotenv.Load("D:\\dev\\go\\src\\03-side-projects\\go-login\\.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestCreateUser(t *testing.T) {
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	ur := database.NewUserRepositoryImpl(gorm)
	defer db.Close()
	t.Run("Create user successfully", func(t *testing.T) {
		expected := domain.UserLogin{Email: "user1@user.com", Password: "user1"}
		_, err := ur.CreateUser(&expected)
		if assert.NoError(t, err) {
			result, _ := ur.GetUserByEmail("user1@user.com")
			r := *result
			assert.Equal(t, expected.Email, r.Email)
			assert.Equal(t, expected.Password, r.Password)
		}
	})
	t.Run("Create user fail, Existed Email", func(t *testing.T) {
		expected := domain.UserLogin{Email: "user1@user.com", Password: "user1"}
		_, err := ur.CreateUser(&expected)
		assert.Error(t, err, "Email already exist")
	})
}

// Get User by ID
func TestGetUserByID(t *testing.T) {
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ur := database.NewUserRepositoryImpl(gorm)
	t.Run("Get user successfully", func(t *testing.T) {
		expected := domain.UserLogin{Email: "user1@user.com", Password: "user1"}
		ur.CreateUser(&expected)
		result, _ := ur.GetUserByEmail("user1@user.com")
		r := *result
		res, err := ur.GetUserByID(r.ID)
		assert.NoError(t, err)
		assert.Equal(t, res.Email, res.Email)
		assert.Equal(t, res.Password, res.Password)
	})
	t.Run("Not found this ID in database", func(t *testing.T) {

		result, err := ur.GetUserByID(999)
		assert.Error(t, err, "Not found this ID in database")
		var nilpt *domain.UserLogin
		assert.Equal(t, nilpt, result)
	})
}

// GetUserByEmail
func TestGetUserByEmail(t *testing.T) {
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ur := database.NewUserRepositoryImpl(gorm)
	t.Run("Get user successfully", func(t *testing.T) {
		expected := domain.UserLogin{Email: "user1@user.com", Password: "password1"}
		_, err := ur.CreateUser(&expected)
		if err != nil {
			t.Fatal(err)
		}
		result, err := ur.GetUserByEmail("user1@user.com")
		assert.NoError(t, err)
		r := *result
		assert.Equal(t, expected.Email, r.Email)
		assert.Equal(t, expected.Password, r.Password)
	})
	t.Run("Not found this email in database", func(t *testing.T) {

		result, err := ur.GetUserByEmail("user2@user.com")
		assert.Error(t, err, "Not found this email in database")
		var nilpt *domain.UserLogin
		assert.Equal(t, nilpt, result)
	})
}

// // GetUsers
func TestGetUsers(t *testing.T) {
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	gorm.AutoMigrate(&domain.Session{})
	defer db.Close()
	ur := database.NewUserRepositoryImpl(gorm)
	t.Run("No users in database", func(t *testing.T) {
		result, err := ur.GetUsers()
		assert.NoError(t, err)
		var nilpt *[]domain.UserLogin
		assert.Nil(t, nilpt, result)
	})

	t.Run("Get users successfully", func(t *testing.T) {
		expected := []domain.UserLogin{
			{Email: "user1@user.com", Password: "user1"},
			{Email: "user2@user.com", Password: "user2"},
		}
		_, _ = ur.CreateUser(&expected[0])
		_, _ = ur.CreateUser(&expected[1])
		result, err := ur.GetUsers()
		assert.NoError(t, err)
		r := *result
		assert.Equal(t, expected[0].Email, r[0].Email)
		assert.Equal(t, expected[0].Password, r[0].Password)
		assert.Equal(t, expected[1].Email, r[1].Email)
		assert.Equal(t, expected[1].Password, r[1].Password)
	})
}

// // UpdateUser
func TestUpdateUser(t *testing.T) {
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	gorm.AutoMigrate(&domain.Session{})
	defer db.Close()
	ur := database.NewUserRepositoryImpl(gorm)
	t.Run("Update user successfully", func(t *testing.T) {
		oldUser := domain.UserLogin{Email: "old@user.com", Password: "oldpass"}
		expected := domain.UserLogin{Email: "new@user.com", Password: "newpass"}

		_, err := ur.CreateUser(&oldUser)
		if err != nil {
			t.Fatal("Creating user failed: ", err)
		}
		updateErr := ur.UpdateUser("old@user.com", &expected)
		result, GetErr := ur.GetUserByEmail("new@user.com")
		if GetErr != nil {
			t.Fatal("Updated user not found: " + GetErr.Error())
		}

		r := *result //Convert pointer to concrete type

		assert.NoError(t, updateErr)
		assert.Equal(t, expected.Email, r.Email)
		assert.Equal(t, expected.Password, r.Password)
	})
}

// DeleteUser
func TestDeleteUser(t *testing.T) {
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ur := database.NewUserRepositoryImpl(gorm)
	t.Run("Delete user successfully", func(t *testing.T) {
		user := domain.UserLogin{Email: "user1@user.com", Password: "user1"}

		_, err := ur.CreateUser(&user)
		if err != nil {
			t.Fatal("Creating user failed: ", err)
		}
		err = ur.DeleteUser(user.Email)
		assert.NoError(t, err)

		result, err := ur.GetUserByEmail(user.Email)
		var nilpointer *domain.UserLogin
		assert.Equal(t, nilpointer, result)
		assert.Error(t, err)
	})

}
