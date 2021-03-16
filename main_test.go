package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"notes/models"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func UserTest() *models.Account {
	return &models.Account{
		Username: "tum0xa",
		Password: "secret123",
	}
}

func UserDataTest() io.Reader {
	user := UserTest()
	return AsJSONBody(Object{
		"username": user.Username,
		"password": user.Password,
	})
}

type NotesTestSuite struct {
	suite.Suite
	ts *httptest.Server
}

func (n *NotesTestSuite) SetupSuite() {
	conn, err := gorm.Open(sqlite.Open(":memory:"))
	n.Require().Nil(err, "test db should work")

	n.ts = httptest.NewServer(GetRouter())

	models.Init(conn)
}

func (n *NotesTestSuite) SetupTest() {
	models.Migrate()
}

func (n *NotesTestSuite) TearDownTest() {
	models.Truncate()
}

func (n *NotesTestSuite) TearDownSuite() {
	n.ts.Close()
}

func (n *NotesTestSuite) TestCreateUser() {
	r := Must(http.Post(n.ts.URL+"/api/user/create", "application/json", UserDataTest()))

	n.Require().Equal(200, r.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Status, rd.Message)
}

func (n *NotesTestSuite) TestLogin() {
	user := UserTest()
	user.Password = models.HashPassword(user.Password)
	models.GetDB().Create(user)

	r := Must(http.Post(n.ts.URL+"/api/user/login", "application/json", UserDataTest()))

	n.Require().Equal(200, r.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Status, rd.Message)
}

func TestNotesTestSuite(t *testing.T) {
	suite.Run(t, &NotesTestSuite{})
}

func Must(resp *http.Response, err error) *http.Response {
	if err != nil {
		panic("when working with test server: " + err.Error())
	}

	return resp
}

func AsJSONBody(obj interface{}) io.Reader {
	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(obj)
	return b
}

type Object map[string]interface{}

type ResponseData struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
