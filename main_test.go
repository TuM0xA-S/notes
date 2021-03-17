package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func UserBodyDataTest() io.Reader {
	user := UserTest()
	return AsJSONBody(Object{
		"username": user.Username,
		"password": user.Password,
	})
}

func CreateUserTest() *models.Account {
	user := UserTest()
	user.Password = models.HashPassword(user.Password)
	models.GetDB().Create(user)

	return user
}

func NoteBodyDataTest(title, body string) io.Reader {
	return AsJSONBody(Object{
		"title": title,
		"body":  body,
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
	r := Must(http.Post(n.ts.URL+"/api/user/create", "application/json", UserBodyDataTest()))

	n.Require().Equal(200, r.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Status, rd.Message)

	user := UserTest()

	n.Require().Nil(models.GetDB().First(&models.Account{}, "username = ?", user.Username).Error)
}

func (n *NotesTestSuite) TestLogin() {
	CreateUserTest()
	r := Must(http.Post(n.ts.URL+"/api/user/login", "application/json", UserBodyDataTest()))

	n.Require().Equal(200, r.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Status, rd.Message)

	n.Require().NotEmpty(rd.AccessToken, "should send access_token on login")
}

func AuthorizeRequest(req *http.Request, user *models.Account) {
	req.Header.Add("Authorization", "Bearer "+models.GenerateToken(user.ID))
}

func (n *NotesTestSuite) TestCreateNote() {
	client := &http.Client{}

	title := "just text"
	body := "another text"

	req, _ := http.NewRequest("POST", n.ts.URL+"/api/me/notes/create", NoteBodyDataTest(title, body))
	AuthorizeRequest(req, CreateUserTest())

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))

	n.Require().True(rd.Status, rd.Message)

	n.Require().Nil(models.GetDB().First(&models.Note{}, "title = ?", title).Error)
}

func (n *NotesTestSuite) TestNotesList() {
	user := CreateUserTest()
	expectedTitles := []string{"title 1", "title 2", "another stuff"}
	for _, title := range expectedTitles {
		models.GetDB().Create(&models.Note{Title: title, UserID: user.ID})
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", n.ts.URL+"/api/me/notes", nil)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))

	actualTitles := []string{}
	for _, n := range rd.Notes {
		actualTitles = append(actualTitles, n.Title)
	}

	n.Require().ElementsMatch(actualTitles, expectedTitles)
}

func (n *NotesTestSuite) TestNoteDetail() {
	user := CreateUserTest()
	expectedNote := &models.Note{
		Title:  "wow nice title bruh",
		Body:   "text text text text text",
		UserID: user.ID,
	}
	expectedNote.Create()

	client := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf(n.ts.URL+"/api/me/notes/%d", expectedNote.ID), nil)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))

	n.Assert().Equal(expectedNote.Title, rd.Note.Title)
	n.Assert().Equal(expectedNote.Body, rd.Note.Body)
	n.Assert().Equal(expectedNote.ID, rd.Note.ID)
}

func (n *NotesTestSuite) TestNoteRemove() {
	user := CreateUserTest()
	expectedNote := &models.Note{
		Title:  "wow nice title bruh",
		Body:   "text text text text text",
		UserID: user.ID,
	}
	expectedNote.Create()

	client := &http.Client{}

	req, _ := http.NewRequest("POST", fmt.Sprintf(n.ts.URL+"/api/me/notes/%d/remove", expectedNote.ID), nil)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))

	n.Require().NotNil(models.GetDB().First(&models.Note{}, expectedNote.ID).Error)
}

func (n *NotesTestSuite) TestUserDetail() {
	user := CreateUserTest()

	client := &http.Client{}

	req, _ := http.NewRequest("GET", n.ts.URL+"/api/me", nil)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))

	n.Assert().Equal(user.Username, rd.User.Username)
	n.Assert().Equal(user.Password, rd.User.Password)
}

func (n *NotesTestSuite) TestUnauth() {
	qs := []struct {
		method, url string
	}{
		{"GET", "/api/me/notes"},
		{"POST", "/api/me/notes/create"},
		{"GET", "/api/me/notes/1"},
		{"POST", "/api/me/notes/1/remove"},
		{"GET", "/api/me"},
	}

	client := &http.Client{}
	for _, x := range qs {
		req, _ := http.NewRequest(x.method, n.ts.URL+x.url, nil)
		resp := Must(client.Do(req))
		n.Assert().Equal(403, resp.StatusCode)
	}
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
	Status      bool           `json:"status"`
	Message     string         `json:"message"`
	Notes       []models.Note  `json:"notes"`
	AccessToken string         `json:"access_token"`
	Note        models.Note    `json:"note"`
	User        models.Account `json:"user"`
}
