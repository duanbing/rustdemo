module github.com/xuperdata/xuperdid

go 1.14

replace golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d => github.com/golang/oauth2 v0.0.0-20200107190931-bf48bf16ab8d

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-session/session v3.1.2+incompatible
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/oauth2.v3 v3.12.0
)
