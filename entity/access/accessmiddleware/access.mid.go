package accessmiddleware

import (
	"fmt"
	"gopher/entity/access/accessenum"
	"gopher/entity/access/accessrepo"
	"gopher/entity/auth/authenum"
	"gopher/entity/service"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accessMiddleware struct {
	engine *core.Engine
}

// NewAccessMiddleware is a simpler way for access to the struct
func NewAccessMiddleware(engine *core.Engine) accessMiddleware {
	return accessMiddleware{
		engine: engine,
	}

}

// Check will analyze if the user should have access to special resource or not
func (p *accessMiddleware) Check(resource string) gin.HandlerFunc {

	return func(c *gin.Context) {

		_, params := response.NewParam(p.engine, c, authenum.Entity)
		authCacheKey := fmt.Sprintf("%v-%v", terms.Auth, params.UserID)

		accessService := service.ProvideAccessService(accessrepo.ProvideAccessRepo(p.engine))
		accessResult := accessService.CheckAccess(c, resource, authCacheKey, params)

		if c.Query("deleted") == "true" {
			accessResult = accessService.CheckAccess(c, accessenum.ReadDeleted, authCacheKey, params)
		}

		if !accessResult {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": terms.YouDoNotHavePermission})
			return
		}

		c.Next()
	}
}
