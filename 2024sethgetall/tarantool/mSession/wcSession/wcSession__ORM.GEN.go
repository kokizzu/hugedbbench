package wcSession

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"hugedbbench/2024sethgetall/tarantool/mSession/rqSession"

	"github.com/tarantool/go-tarantool/v2"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

// SessionsMutator DAO writer/command struct
//
//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcSession__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type wcSession__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type wcSession__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type wcSession__ORM.GEN.go
type SessionsMutator struct {
	rqSession.Sessions
	mutations *tarantool.Operations
	logs      []A.X
}

// NewSessionsMutator create new ORM writer/command object
func NewSessionsMutator(adapter *Tt.Adapter) (res *SessionsMutator) {
	res = &SessionsMutator{Sessions: rqSession.Sessions{Adapter: adapter}}
	res.mutations = tarantool.NewOperations()
	return
}

// Logs get array of logs [field, old, new]
func (s *SessionsMutator) Logs() []A.X { //nolint:dupl false positive
	return s.logs
}

// HaveMutation check whether Set* methods ever called
func (s *SessionsMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(s.logs) > 0
}

// ClearMutations clear all previously called Set* methods
func (s *SessionsMutator) ClearMutations() { //nolint:dupl false positive
	s.mutations = tarantool.NewOperations()
	s.logs = []A.X{}
}

// func (s *SessionsMutator) DoUpsert() bool { //nolint:dupl false positive
//	arr := s.ToArray()
//	_, err := s.Adapter.Upsert(s.SpaceName(), arr, A.X{
//		A.X{`=`, 0, s.Id},
//		A.X{`=`, 1, s.Email},
//		A.X{`=`, 2, s.Permission},
//		A.X{`=`, 3, s.ExpiredAt},
//		A.X{`=`, 4, s.SessionKey},
//	})
//	return !L.IsError(err, `Sessions.DoUpsert failed: `+s.SpaceName()+ `\n%#v`, arr)
// }

// DoOverwriteBySessionKey update all columns, error if not exists, not using mutations/Set*
func (s *SessionsMutator) DoOverwriteBySessionKey() bool { //nolint:dupl false positive
	_, err := s.Adapter.Connection.Do(tarantool.NewUpdateRequest(s.SpaceName()).
		Index(s.UniqueIndexSessionKey()).
		Key(A.X{s.SessionKey}).
		Operations(s.ToUpdateArray()),
	).Get()
	return !L.IsError(err, `Sessions.DoOverwriteBySessionKey failed: `+s.SpaceName())
}

// DoUpdateBySessionKey update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (s *SessionsMutator) DoUpdateBySessionKey() bool { //nolint:dupl false positive
	if !s.HaveMutation() {
		return true
	}
	_, err := s.Adapter.Connection.Do(
		tarantool.NewUpdateRequest(s.SpaceName()).
			Index(s.UniqueIndexSessionKey()).
			Key(A.X{s.SessionKey}).
			Operations(s.mutations),
	).Get()
	return !L.IsError(err, `Sessions.DoUpdateBySessionKey failed: `+s.SpaceName())
}

// DoDeletePermanentBySessionKey permanent delete
func (s *SessionsMutator) DoDeletePermanentBySessionKey() bool { //nolint:dupl false positive
	_, err := s.Adapter.Connection.Do(
		tarantool.NewDeleteRequest(s.SpaceName()).
			Index(s.UniqueIndexSessionKey()).
			Key(A.X{s.SessionKey}),
	).Get()
	return !L.IsError(err, `Sessions.DoDeletePermanentBySessionKey failed: `+s.SpaceName())
}

// DoInsert insert, error if already exists
func (s *SessionsMutator) DoInsert() bool { //nolint:dupl false positive
	arr := s.ToArray()
	_, err := s.Adapter.Connection.Do(
		tarantool.NewInsertRequest(s.SpaceName()).
			Tuple(arr),
	).Get()
	return !L.IsError(err, `Sessions.DoInsert failed: `+s.SpaceName()+`\n%#v`, arr)
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// replace = upsert, only error when there's unique secondary key
// previous name: DoReplace
func (s *SessionsMutator) DoUpsert() bool { //nolint:dupl false positive
	arr := s.ToArray()
	_, err := s.Adapter.Connection.Do(
		tarantool.NewReplaceRequest(s.SpaceName()).
			Tuple(arr),
	).Get()
	return !L.IsError(err, `Sessions.DoUpsert failed: `+s.SpaceName()+`\n%#v`, arr)
}

// SetId create mutations, should not duplicate
func (s *SessionsMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != s.Id {
		s.mutations.Assign(0, val)
		s.logs = append(s.logs, A.X{`id`, s.Id, val})
		s.Id = val
		return true
	}
	return false
}

// SetEmail create mutations, should not duplicate
func (s *SessionsMutator) SetEmail(val string) bool { //nolint:dupl false positive
	if val != s.Email {
		s.mutations.Assign(1, val)
		s.logs = append(s.logs, A.X{`email`, s.Email, val})
		s.Email = val
		return true
	}
	return false
}

// SetPermission create mutations, should not duplicate
func (s *SessionsMutator) SetPermission(val string) bool { //nolint:dupl false positive
	if val != s.Permission {
		s.mutations.Assign(2, val)
		s.logs = append(s.logs, A.X{`permission`, s.Permission, val})
		s.Permission = val
		return true
	}
	return false
}

// SetExpiredAt create mutations, should not duplicate
func (s *SessionsMutator) SetExpiredAt(val int64) bool { //nolint:dupl false positive
	if val != s.ExpiredAt {
		s.mutations.Assign(3, val)
		s.logs = append(s.logs, A.X{`expiredAt`, s.ExpiredAt, val})
		s.ExpiredAt = val
		return true
	}
	return false
}

// SetSessionKey create mutations, should not duplicate
func (s *SessionsMutator) SetSessionKey(val string) bool { //nolint:dupl false positive
	if val != s.SessionKey {
		s.mutations.Assign(4, val)
		s.logs = append(s.logs, A.X{`sessionKey`, s.SessionKey, val})
		s.SessionKey = val
		return true
	}
	return false
}

// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap
func (s *SessionsMutator) SetAll(from rqSession.Sessions, excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive
	if excludeMap == nil { // list of fields to exclude
		excludeMap = M.SB{}
	}
	if forceMap == nil { // list of fields to force overwrite
		forceMap = M.SB{}
	}
	if !excludeMap[`id`] && (forceMap[`id`] || from.Id != 0) {
		s.Id = from.Id
		changed = true
	}
	if !excludeMap[`email`] && (forceMap[`email`] || from.Email != ``) {
		s.Email = S.Trim(from.Email)
		changed = true
	}
	if !excludeMap[`permission`] && (forceMap[`permission`] || from.Permission != ``) {
		s.Permission = S.Trim(from.Permission)
		changed = true
	}
	if !excludeMap[`expiredAt`] && (forceMap[`expiredAt`] || from.ExpiredAt != 0) {
		s.ExpiredAt = from.ExpiredAt
		changed = true
	}
	if !excludeMap[`sessionKey`] && (forceMap[`sessionKey`] || from.SessionKey != ``) {
		s.SessionKey = S.Trim(from.SessionKey)
		changed = true
	}
	return
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
