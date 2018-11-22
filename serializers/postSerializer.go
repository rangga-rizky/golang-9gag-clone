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

func (u *PostSerializer) WithComments() *PostSerializer {
	u.PickFunc(func(t interface{}) interface{} {
		comments := []interface{}{}

		switch t := t.(type) {
		case []m.Comment:
			for _, value := range t {
				newComment := map[string]interface{}{}
				newComment["ID"] = value.ID
				newComment["Text"] = value.Text
				newComment["ImagePath"] = value.ImagePath
				newComment["User"] = value.User.Email
				newComment["UserID"] = value.UserID
				newComment["CreatedAt"] = value.CreatedAt
				comments = append(comments, newComment)
			}

		}
		return comments
	}, "Comments")
	return u
}
