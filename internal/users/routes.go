package users

import "net/http"

func RegisterRoutes(mux *http.ServeMux){
	mux.HandleFunc("POST /socialLogin", SocialLogin)
}