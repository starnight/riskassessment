package risk_assessment

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/starnight/riskassessment/backend/database"
)

var scope_utils = ScopeUtils{
  DB_Client: database.ConnectDB(database.GetDBStr("")),
}

func TestAddScope(t *testing.T) {
  scopes := []Scope {{Name: "test1"}, {Name: "test2"}}

  /* Add Scopes */
  id1, err1 := scope_utils.AddScope(&scopes[0])
  assert.Nil(t, err1)
  assert.True(t, id1.Hex() != "")

  id2, err2 := scope_utils.AddScope(&scopes[1])
  assert.Nil(t, err2)
  assert.True(t, id2.Hex() != "")

  /* Get Scope by ID */
  get_scope, err3 := scope_utils.GetScopeByID(id1)
  assert.Nil(t, err3)
  assert.Equal(t, scopes[0].Name, get_scope.Name)

  /* Get Scopes by IDs */
  ids := []primitive.ObjectID{id1, id2}
  get_scopes, err4 := scope_utils.GetScopeByIDs(ids)
  assert.Nil(t, err4)
  assert.Equal(t, len(ids), len(get_scopes))
  assert.Equal(t, scopes[0].Name, get_scopes[0].Name)
  assert.Equal(t, scopes[1].Name, get_scopes[1].Name)

  /* Get all Scopes */
  new_get_scopes, err5 := scope_utils.GetScopes()
  assert.Nil(t, err5)
  assert.Equal(t, get_scopes, new_get_scopes)

  /* Has Scope ID */
  has, err6 := scope_utils.HasScopeID(id1)
  assert.Nil(t, err6)
  assert.True(t, has)

  /* Update a Scope */
  scope := get_scopes[0]
  scope.Name = "new test1"
  err7 := scope_utils.UpdateScope(&scope)
  assert.Nil(t, err7)

  ids = []primitive.ObjectID{scope.ID}
  new_scopes, _ := scope_utils.GetScopeByIDs(ids)
  assert.Equal(t, scope, new_scopes[0])
}
