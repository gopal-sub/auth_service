package user


// sign up types
type SignUpRequest struct {
    Email    string 
    Password string 
}

type SignUpResponse struct {
	id string
    Email    string 
}

//signin types

type SigninRequest struct {
    Email    string 
    Password string 
}

type SigninResponse struct {
	id string
    Email    string 
}





