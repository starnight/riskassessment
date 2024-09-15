package auth

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/starnight/riskassessment/backend/database"

  "go.mongodb.org/mongo-driver/bson/primitive"
)

var utils = UserUtils{
  DB_Client: database.ConnectDB(database.GetDBStr("")),
}

func TestAddAsset(t *testing.T) {
  /* Add an User */
  user := User {
    Account: "foo",
    Password: "bar",
  }

  id, err1 := utils.AddUser(&user)

  assert.Nil(t, err1)
  assert.True(t, id.Hex() != "")

  /* Get the User by ID */
  get_user2, err2 := utils.GetUserByID(id)
  assert.Nil(t, err2)
  assert.Equal(t, id, get_user2.ID)
  assert.Equal(t, user.Account, get_user2.Account)
  assert.Equal(t, user.Password, get_user2.Password)

  /* Get the user by account */
  get_user3, err3 := utils.GetUserByAccount(user.Account)
  assert.Nil(t, err3)
  assert.Equal(t, get_user2, get_user3)

  /* Get the User by account and password */
  get_user4, err4 := utils.GetUserByAccountPwd(user.Account, user.Password)
  assert.Nil(t, err4)
  assert.Equal(t, get_user2, get_user4)

  /* Update the user by adding a scope to the user */
  update_user1 := get_user2
  scope_ID := primitive.NewObjectID()
  update_user1.Scopes = append(update_user1.Scopes, scope_ID)
  err5 := utils.UpdateUser(&update_user1)
  assert.Nil(t, err5)
  update_user2, _ := utils.GetUserByID(update_user1.ID)
  assert.Equal(t, update_user1.Scopes, update_user2.Scopes)

  /* Check the user is in the scope */
  in1, err5 := utils.UserHasScopeID(update_user1.ID, scope_ID)
  assert.Nil(t, err5)
  assert.True(t, in1)

  in2, err6 := utils.UserHasScopeID(update_user1.ID, primitive.NewObjectID())
  assert.Nil(t, err6)
  assert.False(t, in2)
}

func TestHasUser(t *testing.T) {
  has, err := utils.HasUser()
  assert.Nil(t, err)
  assert.True(t, has)
}
