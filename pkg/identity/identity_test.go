package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"grpcapi/pkg/pb"
)

func TestIdentity(t *testing.T) {
	//-------- test 1
	id := New("123456", []pb.Role{pb.Role_ADMIN, pb.Role_OWNER}, []string{"team1", "team2"})
	s := id.String()
	t.Log(s)
	assert.Equal(t, "123456|0,1|team1,team2", s)
	id2 := &Identity{}
	ParseString(id2, s)
	assert.Equal(t, "123456", id2.UserID)
	assert.Equal(t, 2, len(id2.Roles))
	assert.Equal(t, 2, len(id2.Scopes))
	assert.Equal(t, pb.Role_ADMIN, id2.Roles[0])
	assert.Equal(t, pb.Role_OWNER, id2.Roles[1])
	assert.Equal(t, "team1", id2.Scopes[0])
	assert.Equal(t, "team2", id2.Scopes[1])

	// ------- test 2
	id = New("1234567", []pb.Role{}, []string{})
	s = id.String()
	t.Log(s)
	assert.Equal(t, "1234567", s)
	id2 = &Identity{}
	ParseString(id2, s)
	assert.Equal(t, "1234567", id2.UserID)
	assert.Equal(t, 0, len(id2.Roles))
	assert.Equal(t, 0, len(id2.Scopes))

	// ------- test 3
	id = New("12345678", []pb.Role{pb.Role_ADMIN}, []string{})
	s = id.String()
	t.Log(s)
	assert.Equal(t, "12345678|0", s)
	id2 = &Identity{}
	ParseString(id2, s)
	assert.Equal(t, "12345678", id2.UserID)
	assert.Equal(t, 1, len(id2.Roles))
	assert.Equal(t, 0, len(id2.Scopes))
	assert.Equal(t, pb.Role_ADMIN, id2.Roles[0])

	// ------- test 4
	id = New("123456789", []pb.Role{}, []string{"team1", "team2", "team3"})
	s = id.String()
	t.Log(s)
	assert.Equal(t, "123456789||team1,team2,team3", s)
	id2 = &Identity{}
	ParseString(id2, s)
	assert.Equal(t, "123456789", id2.UserID)
	assert.Equal(t, 0, len(id2.Roles))
	assert.Equal(t, 3, len(id2.Scopes))
	assert.Equal(t, "team1", id2.Scopes[0])
	assert.Equal(t, "team2", id2.Scopes[1])
	assert.Equal(t, "team3", id2.Scopes[2])
}

func newIdentity() *Identity {
	return New("123456789", []pb.Role{pb.Role_ADMIN, pb.Role_OWNER}, []string{"team1", "team2", "team3"})
}
func BenchmarkToString(b *testing.B) {
	id := newIdentity()
	for i := 0; i < b.N; i++ {
		id.String()
	}
}

func BenchmarkToJSON(b *testing.B) {
	id := newIdentity()
	for i := 0; i < b.N; i++ {
		id.JSON()
	}
}

func BenchmarkParseString(b *testing.B) {
	id := newIdentity()
	s := id.String()
	id2 := &Identity{}
	for i := 0; i < b.N; i++ {
		ParseString(id2, s)
	}
}
func BenchmarkParseJSON(b *testing.B) {
	id := newIdentity()
	j := id.JSON()
	id2 := &Identity{}
	for i := 0; i < b.N; i++ {
		ParseJSON(id2, j)
	}
}
