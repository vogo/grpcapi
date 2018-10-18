package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	//-------- test 1
	id := New("123456", []string{"admin", "manager"}, []string{"team1", "team2"})
	s := id.String()
	t.Log(s)
	assert.Equal(t, "123456|admin,manager|team1,team2", s)
	id2 := &Identity{}
	ParseString(id2, s)
	assert.Equal(t, "123456", id2.UserID)
	assert.Equal(t, 2, len(id2.Roles))
	assert.Equal(t, 2, len(id2.Scopes))
	assert.Equal(t, "admin", id2.Roles[0])
	assert.Equal(t, "manager", id2.Roles[1])
	assert.Equal(t, "team1", id2.Scopes[0])
	assert.Equal(t, "team2", id2.Scopes[1])

	// ------- test 2
	id = New("1234567", []string{}, []string{})
	s = id.String()
	t.Log(s)
	assert.Equal(t, "1234567", s)
	id2 = &Identity{}
	ParseString(id2, s)
	assert.Equal(t, "1234567", id2.UserID)
	assert.Equal(t, 0, len(id2.Roles))
	assert.Equal(t, 0, len(id2.Scopes))

	// ------- test 3
	id = New("12345678", []string{"admin"}, []string{})
	s = id.String()
	t.Log(s)
	assert.Equal(t, "12345678|admin", s)
	id2 = &Identity{}
	ParseString(id2, s)
	assert.Equal(t, "12345678", id2.UserID)
	assert.Equal(t, 1, len(id2.Roles))
	assert.Equal(t, 0, len(id2.Scopes))
	assert.Equal(t, "admin", id2.Roles[0])

	// ------- test 4
	id = New("123456789", []string{}, []string{"team1", "team2", "team3"})
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

func BenchmarkToString(b *testing.B) {
	id := New("123456", []string{"admin", "manager"}, []string{"team1", "team2"})
	for i := 0; i < b.N; i++ {
		id.String()
	}
}

func BenchmarkToJSON(b *testing.B) {
	id := New("123456", []string{"admin", "manager"}, []string{"team1", "team2"})
	for i := 0; i < b.N; i++ {
		id.JSON()
	}
}

func BenchmarkParseString(b *testing.B) {
	id := New("123456", []string{"admin", "manager"}, []string{"team1", "team2"})
	s := id.String()
	id2 := &Identity{}
	for i := 0; i < b.N; i++ {
		ParseString(id2, s)
	}
}
func BenchmarkParseJSON(b *testing.B) {
	id := New("123456", []string{"admin", "manager"}, []string{"team1", "team2"})
	j := id.JSON()
	id2 := &Identity{}
	for i := 0; i < b.N; i++ {
		ParseJSON(id2, j)
	}
}
