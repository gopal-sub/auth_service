package user

import (
	"database/sql"
	"errors"
	"testing"
)

type FakeRepository struct {
	user User
	findErr error
	createErr error
}




func (f *FakeRepository) FindUserByEmail(email string) (User, error){
	return f.user, f.findErr
}

func (f *FakeRepository) Create(user User) error {
	return  f.createErr
}


func TestSigninUserNotFound(t *testing.T) {
	//arrange
	fakeRepo := &FakeRepository {
		findErr: sql.ErrNoRows,
	}
	service := NewService(fakeRepo)
	//act
	token, err := service.Signin("alice@example.com", "hi there")
	//assert
	if !errors.Is(err, ErrInvalidCredentials){
		t.Fatal("exprected ErrInvalidCredentials")
	}
	if token != ""{
		t.Fatal("exprected empty token")
	}

}


func TestDBConnectionErr(t *testing.T){
	// arrange
	fakeRepo := &FakeRepository{
		findErr: dberr,
	}
	service := NewService(fakeRepo)
	

	// act
	token, err := service.Signin("testuser@gmail.com", "hi there")

	// assert
	if !errors.Is(err, dberr){
		t.Fatal("expected database conenction error")
	}
	if token != ""{
		t.Fatal("expected empty token")
	}
}

func TestWrongPassword(t *testing.T) {
	// arrange
	password := "hit ehre"
	hashPassword, err := CreateHash(password)

	testUser := User{Email: "testuser123@gmail.com", PasswordHash: hashPassword}
	fakeRepo := &FakeRepository{
		user: testUser,
		findErr: nil,
	}
	service := NewService(fakeRepo)


	// act
	token, err := service.Signin("testuser123@gmail.com", "password")

	// assert
	if !errors.Is(err, ErrInvalidCredentials){
		t.Fatal("expected ErrInvalidCredentials")
	}
	if token != ""{
		t.Fatal("expected empty token")
	}
}

func TestRightUserAndPassword(t *testing.T) {
	// arrange
	password := "hit ehre"
	hashPassword, _ := CreateHash(password)

	testRepo := &FakeRepository{
		user: User{Email: "test123@gmail.com", PasswordHash: hashPassword},
		findErr: nil,
	}

	service := NewService(testRepo)

	// act

	token, err := service.Signin("test123@gmail.com", password)

	
	// assert

	if err != nil{
		t.Error("expected no error for a succcessfull sigjin")
	}
	if token == "" {
		t.Fatal("a valid jwt token was expected")
	}
}


//signup test
// -> test for successful signup
// -> database error find user
// -> database error creating the user
// -> user already exists


func Test_SignupSuccesss(t *testing.T){
	
	// arrange
	repo := &FakeRepository{
		findErr: sql.ErrNoRows,
		createErr: nil,
	}
	service := NewService(repo)
	// act

	_ , err := service.Signup("different@gmail.com", "hidahida")
	// assert
	if err != nil {
		t.Fatalf("expected no error for a success sign up error: %v", err)
	}
	

}
func Test_DatabaseErrFindUser(t *testing.T){
	// arrange
	repo := &FakeRepository{
		findErr: dberr,
		createErr: nil,
	}
	service := NewService(repo)


	// act

	_, err := service.Signup("testuser@gmail.com", "hi there")
	// assert
	if !errors.Is(err, dberr){
		t.Fatalf("expected error dberr got : %v", err)
	}
}
func Test_DatabaseErrInsertUser(t *testing.T){
	// arrange
	// act
	// assert
}
func Test_UserExistsConflict(t *testing.T){
	// arrange
	// act
	// assert
}
