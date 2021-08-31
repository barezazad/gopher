package middleware

import (
	"fmt"
	"gopher/entity/auth/authenum"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/model"
	"gopher/internal/response"
	"gopher/pkg/dictionary"
	"gopher/pkg/generr"
	"gopher/pkg/helper"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthGuard is used for decode the token and get public and private information
func AuthGuard(engine *core.Engine) gin.HandlerFunc {

	jwtKey := []byte(engine.Environments.JWT.SecretKey)
	fJWT := func(token *jwt.Token) (interface{}, error) { return jwtKey, nil }

	return func(c *gin.Context) {

		var token string
		// token = strings.TrimSpace(c.Query("temporary_token"))
		if token == "" {
			tokenArr, ok := c.Request.Header["Authorization"]
			if !ok || len(tokenArr[0]) == 0 {
				err := engine.ErrorLog.TickCustom(fmt.Errorf(terms.TokenIsRequired),
					"E1000060", generr.UnauthorizedErr, "", terms.TokenIsRequired)
				response.New(engine, c, authenum.Entity).Error(err).Abort().JSON()
				return
			}

			token = tokenArr[0][7:]
		}

		claims := &model.JWTClaims{}

		if tkn, err := jwt.ParseWithClaims(token, claims, fJWT); err != nil {
			checkErr(c, err, engine)
			return
		} else if !tkn.Valid {
			checkToken(c, tkn, engine)
			return
		}

		lang := c.Request.Header.Get("Accepted-Language")
		isInclodeLang, _ := helper.Includes(engine.Environments.Languages, lang)
		if lang == "" || !isInclodeLang {
			lang = string(dictionary.En)
		}

		c.Set("USERNAME", claims.Username)
		c.Set("USER_ID", claims.ID)
		c.Set("LANGUAGE", lang)
		c.Set("TOKEN", token)
		c.Next()
	}

}

func checkErr(c *gin.Context, err error, engine *core.Engine) {
	if err != nil {

		if err == jwt.ErrSignatureInvalid {
			err = engine.ErrorLog.TickCustom(fmt.Errorf(terms.TokenIsNotValid),
				"E1000060", generr.UnauthorizedErr, "", terms.TokenIsNotValid)
			response.New(engine, c, authenum.Entity).Error(err).Abort().JSON()
			return
		}

		err = engine.ErrorLog.TickCustom(fmt.Errorf(terms.TokenIsExpired),
			"E1000060", generr.UnauthorizedErr, "", terms.TokenIsExpired)
		response.New(engine, c, authenum.Entity).Error(err).Abort().JSON()
		return
	}
}

func checkToken(c *gin.Context, token *jwt.Token, engine *core.Engine) {
	if !token.Valid {
		err := engine.ErrorLog.TickCustom(fmt.Errorf(terms.TokenIsNotValid),
			"E1000060", generr.UnauthorizedErr, "", terms.TokenIsNotValid)
		response.New(engine, c, "auth").Error(err).Abort().JSON()
		return
	}
}
