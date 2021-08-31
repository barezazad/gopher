package accessrepo

import (
	"gopher/entity/user/usermodel"
	"gopher/internal/core"
)

// AccessRepo for injecting engine
type AccessRepo struct {
	Engine *core.Engine
}

// ProvideAccessRepo is used in wire
func ProvideAccessRepo(engine *core.Engine) AccessRepo {
	return AccessRepo{Engine: engine}
}

// GetUserResources is used for finding all resources
func (p *AccessRepo) GetUserResources(userID uint) (result string, err error) {
	resources := struct {
		Resources string
	}{}

	err = p.Engine.DB.Table(usermodel.Table).Select("roles.resources").
		Joins("INNER JOIN roles ON users.role_id = roles.id").
		Where("users.id = ?", userID).Find(&resources).Error

	result = resources.Resources

	return
}
