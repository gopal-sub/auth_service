package oauth

// import (
// 	// "gopal-sub/auth_service/internal/user"
// 	// "gopal-sub/auth_service/internal/auth/oauth/google"
// )


// type Service struct {
// 	googleProvider *google.GoogleProvider
// 	userService user.UserService
// }

// func NewService(googleProvider *google.GoogleProvider, userService user.UserService) *Service{
// 	return &Service{
// 		googleProvider: googleProvider,
// 		userService: userService,
// 	}
// }


// type Oauthsvc interface {
// 	GetAuthLoginURL(provider OauthProvider) string
// 	HandleCallback(provider OauthProvider, code string) (string, error)
// }

// func (s *Service) GetAuthLoginURL(provider OauthProvider) string {
// 	if provider == "google" {
// 		return s.googleProvider.GetAuthURL()
// 	}
// 	return "ada"\}