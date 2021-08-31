package service

import (
	"fmt"
	"gopher/entity/access/accessenum"
	"gopher/entity/access/accessrepo"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/param"
	"strings"

	"github.com/gin-gonic/gin"
)

// AccessService defining auth service
type AccessService struct {
	Repo   accessrepo.AccessRepo
	Engine *core.Engine
}

// ProvideAccessService for auth is used in wire
func ProvideAccessService(p accessrepo.AccessRepo) AccessService {
	return AccessService{Repo: p, Engine: p.Engine}
}

// CheckAccess is used inside each method to find out if user has permission or not
func (p *AccessService) CheckAccess(c *gin.Context, resource, authCacheKey string, params param.Param) bool {

	exist, err := p.Engine.Cache.KeyExist(authCacheKey)
	if err != nil {
		return false
	}

	if !exist {
		resources, err := p.Repo.GetUserResources(params.UserID)
		if err != nil || resources == "" {
			p.Engine.ServerLog.CheckError(err, "E1000079", "can't finding the resources for user", params.Username)
			return false
		}
		splitResources := strings.Split(strings.TrimSpace(resources), ",")
		p.Engine.Cache.UpdateSet(authCacheKey, splitResources, p.Engine.Environments.Setting.PermissionTTL)
	}

	canAccess, err := p.Engine.Cache.IsMemberSet(authCacheKey, resource)
	if err != nil {
		p.Engine.ServerLog.Fatal(err)
		return false
	}

	if !canAccess {
		return false
	}

	return canAccess
}

func IsSuperAdmin(engine *core.Engine, userID uint) bool {

	authCacheKey := fmt.Sprintf("%v-%v", terms.Auth, userID)
	canAccess, err := engine.Cache.IsMemberSet(authCacheKey, accessenum.SuperAccess)
	if err != nil || !canAccess {
		return false
	}
	return true
}
