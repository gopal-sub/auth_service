package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeUserService struct {}



func (f *FakeUserService) Signup(email string, password string) (User, error){
	return User{
		Email: email,
	}, nil
}
func (f *FakeUserService) Signin(email string, password string) (token string, err error){
	return "", nil
}
func (f *FakeUserService) UserExists(email string) bool{
	return false
}
func (f *FakeUserService) IsUserAccountVerified(email string) (bool, error){
	return false, nil
}
func (f *FakeUserService) VerifyAnAccount(email string)(accountVerified bool, err error){
	return false, nil
}
func (f *FakeUserService) PasswordReset(email string, OldPassword string) error{
	return nil
}
 

//test Signup handler

type SignupCases struct {
	email string
	password string
	method string
	responseStatus int
	responseMessage string

}

func TestSignUpHandler(t *testing.T){
	// arrange 
	fakeSvc := &FakeUserService{}
	fakehandler := NewHandler(fakeSvc)


	// A Reader implements the [io.Reader]
	body := strings.NewReader(`{
		"email": "test@gmail.com",
		"password": "hi there"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/signup",
		body,
	)

	rec := httptest.NewRecorder()

	// act 

	fakehandler.SignUpHandler(rec, req)

	// assert

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, rec.Code)
	}


}