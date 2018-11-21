package serializers

import "github.com/tuvistavie/structomap"

type UserSerializer struct {
	*structomap.Base
}

func CustomUserSerializer() *UserSerializer {
	u := &UserSerializer{structomap.New()}
	u.Pick("ID", "email")

	return u
}

/*func (u *UserSerializer) WithUser() *UserSerializer {
	u.Pick("User")
	return u
}*/
