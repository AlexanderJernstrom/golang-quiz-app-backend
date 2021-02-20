package middlewares

import (
	"fmt"
	"log"
	database "server/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
)


func IsAuth() func(c *fiber.Ctx)  {
	return func(c *fiber.Ctx) {
		notParsedToken := c.Get("authToken")
		claims := new(database.Claims)
		token, err := jwt.ParseWithClaims(notParsedToken, claims,func (token *jwt.Token) (interface{}, error){
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
				return nil, fmt.Errorf("unecpected signing method")
			} else if ok == false{
				c.Send("Token is invalid")
			}
			return []byte(database.Key), nil
		})
		if token.Valid {
			c.Send("Token is not valid")
		}
		if err != nil{
			log.Fatal(err)
		}	
		//fmt.Println(claims.Issuer)
		c.Locals("id", claims.ID)
		c.Next()
		
	}
}