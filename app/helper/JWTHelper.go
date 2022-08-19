package helper

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// DecodedToken Decode a JWT token without validation for getting payload data (claim)
func DecodedJWTToken(c *gin.Context) (map[string]interface{}, error) {

	const BearerSchema = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BearerSchema):]
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil

}
