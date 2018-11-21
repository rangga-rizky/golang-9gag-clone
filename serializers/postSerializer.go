package serializers

import (
	m "../models"
	"github.com/tuvistavie/structomap"
)

type PostSerializer struct {
	*structomap.Base
}

func CustomPostSerializer() *PostSerializer {
	u := &PostSerializer{structomap.New()}
	u.Pick("ID", "Title", "ImagePath", "CreatedAt", "UpdatedAt")
	u.PickFunc(func(t interface{}) interface{} {
		user := map[string]interface{}{}
		user["ID"] = t.(m.Account).ID
		user["email"] = t.(m.Account).Email
		return user
	}, "User")
	u.PickFunc(func(t interface{}) interface{} {
		section := map[string]interface{}{}
		section["ID"] = t.(m.Section).ID
		section["name"] = t.(m.Section).Name
		return section
	}, "Section")
	return u
}

/*func (u *PostSerializer) WithUser() *PostSerializer {
	u.Pick("User")
	return u
}*/
