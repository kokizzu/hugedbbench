package rqSession

import (
	"hugedbbench/2024sethgetall/testcase"

	"github.com/kokizzu/gotro/S"
)

func (s *Sessions) ToSession() (res testcase.Session) {
	res.Id = int64(s.Id)
	res.Email = s.Email
	res.Permission = map[string]bool{}
	perms := S.Split(s.Permission, ` `)
	for _, perm := range perms {
		res.Permission[perm] = true
	}
	return
}
