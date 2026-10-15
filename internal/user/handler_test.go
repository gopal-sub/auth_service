package user

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// type FakeUserService struct {}



// func (f *FakeUserService) Signup(email string, password string) (User, error){
// 	return User{
// 		Email: email,
// 	}, nil
// }
// func (f *FakeUserService) Signin(email string, password string) (token string, err error){
// 	return "", nil
// }
// func (f *FakeUserService) UserExists(email string) bool{
// 	return false
// }
// func (f *FakeUserService) IsUserAccountVerified(email string) (bool, error){
// 	return false, nil
// }
// func (f *FakeUserService) VerifyAnAccount(email string)(accountVerified bool, err error){
// 	return false, nil
// }
// func (f *FakeUserService) PasswordReset(email string, OldPassword string) error{
// 	return nil
// }
 

// //test Signup handler

// type SignupCases struct {
// 	name string
// 	email string
// 	password string
// 	responseStatus int
// 	responseMessage string

// }

// func TestSignUpHandler(t *testing.T){
// 	// arrange 
// 	fakeSvc := &FakeUserService{}
// 	fakehandler := NewHandler(fakeSvc)

// 	testCases := []SignupCases{
// 		{
// 			name: "Success creating user",
// 			email:"testuser@example.com",
// 			password: "test1243",
// 			responseStatus: 201,
// 			responseMessage: "user created",
	
// 		},
// 		{
// 			name: "failure email invalid",
// 			email:"testuser",
// 			password: "test1243",
// 			responseStatus: 400,
// 			responseMessage: "invalid email",
	
// 		},
// 		{
// 			name: "Bad body",
// 			email:"testuser@exa",
// 			password: "",
// 			responseStatus: 400,
// 			responseMessage: "invalid email",
	
// 		},
// 		{
// 			name: "Bad body",
// 			email:"",
// 			password: "nlsnfls",
// 			responseStatus: 400,
// 			responseMessage: "invalid email",
	
// 		},

// 	}

// 	for _, tt := range testCases{
// 		t.Run(tt.name, func(t *testing.T){
// 			// A Reader implements the [io.Reader]

// 			bodyStr := fmt.Sprintf(`{
// 				"email": "%s",
// 				"password": "%s"
// 			}`, tt.email, tt.password)

// 			body := strings.NewReader(bodyStr)
			
			
// 			req := httptest.NewRequest(
// 				http.MethodPost,
// 				"/signup",
// 				body,
// 			)
			
// 			rec := httptest.NewRecorder()
			
// 			// act 
			
// 			fakehandler.SignUpHandler(rec, req)
			
// 			// assert
// 			if rec.Code != tt.responseStatus {
// 				t.Fatalf("expected %d got %d", tt.responseStatus, rec.Code)
// 			} 

// 		})
// 	}
// }