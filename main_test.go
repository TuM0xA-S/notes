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

	. "notes/config"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func UserTest() *models.User {
	return &models.User{
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

func CreateUserTest() *models.User {
	user := UserTest()
	user.Create()
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
	n.ts = httptest.NewServer(GetRouter())
}

func (n *NotesTestSuite) SetupTest() {
	conn, err := gorm.Open(sqlite.Open(":memory:"))
	n.Require().Nil(err, "test db should work")

	models.Init(conn)
	models.Migrate()
}

func (n *NotesTestSuite) TearDownTest() {
	models.Truncate()
}

func (n *NotesTestSuite) TearDownSuite() {
	n.ts.Close()
}

func (n *NotesTestSuite) TestCreateUser() {
	r := Must(http.Post(n.ts.URL+"/api/users", "application/json", UserBodyDataTest()))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	user := UserTest()

	n.Require().Nil(models.GetDB().First(&models.User{}, "username = ?", user.Username).Error)
}

func (n *NotesTestSuite) TestLogin() {
	CreateUserTest()
	r := Must(http.Post(n.ts.URL+"/api/me", "application/json", UserBodyDataTest()))

	n.Require().Equal(200, r.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	n.Require().NotEmpty(rd.AccessToken, "should send access_token on login")
}

func AuthorizeRequest(req *http.Request, user *models.User) {
	req.Header.Add("Authorization", "Bearer "+models.GenerateToken(user.ID))
}

func (n *NotesTestSuite) TestCreateNote() {
	client := &http.Client{}

	title := "just text"
	body := "another text"

	req, _ := http.NewRequest("POST", n.ts.URL+"/api/me/notes", NoteBodyDataTest(title, body))
	AuthorizeRequest(req, CreateUserTest())

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	n.Require().Nil(models.GetDB().First(&models.Note{}, "title = ?", title).Error)
}

func (n *NotesTestSuite) TestNotesList() {
	user := CreateUserTest()
	titles := []string{"title 1", "title 2", "another stuff"}
	expectedNotes := []models.Note{}
	for _, title := range titles {
		note := models.Note{Title: title, UserID: user.ID}
		note.Create()
		expectedNotes = append(expectedNotes, note)
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", n.ts.URL+"/api/me/notes", nil)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))
	n.Require().True(rd.Success, rd.Message)
	n.Require().NotEmpty(rd.Pagination, "this response should have pagination info")

	n.Require().ElementsMatch(expectedNotes, rd.Notes)
}

func (n *NotesTestSuite) TestPublishedNotesList() {
	user1 := &models.User{
		Username: "user1",
		Password: "password",
		Notes: []models.Note{
			{Title: "title 1", Published: true},
			{Title: "title not published"},
		},
	}
	user1.Create()

	user2 := &models.User{
		Username: "user2",
		Password: "password",
		Notes: []models.Note{
			{Title: "title 2", Published: true},
		},
	}
	user2.Create()

	expectedNotes := []models.Note{user1.Notes[0], user2.Notes[0]}

	resp := Must(http.Get(n.ts.URL + "/api/notes"))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))
	n.Require().True(rd.Success, rd.Message)
	n.Require().NotEmpty(rd.Pagination, "this response should have pagination info")

	n.Require().ElementsMatch(expectedNotes, rd.Notes)
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
	n.Require().True(rd.Success, rd.Message)

	n.Assert().Equal(expectedNote.Title, rd.Note.Title)
	n.Assert().Equal(expectedNote.Body, rd.Note.Body)
	n.Assert().Equal(expectedNote.ID, rd.Note.ID)

	req.URL.RawQuery = "no_body=true"
	resp = Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd = &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))
	n.Require().True(rd.Success, rd.Message)
	n.Assert().Empty(rd.Note.Body)
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

	req, _ := http.NewRequest("DELETE", fmt.Sprintf(n.ts.URL+"/api/me/notes/%d", expectedNote.ID), nil)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))
	n.Require().True(rd.Success, rd.Message)

	n.Require().NotNil(models.GetDB().First(&models.Note{}, expectedNote.ID).Error)
}

func (n *NotesTestSuite) TestNoteUpdate() {
	user := CreateUserTest()
	expectedTitle := "wow nice title bruh"
	note := &models.Note{
		Title:  expectedTitle,
		Body:   "text text text text text",
		UserID: user.ID,
	}
	note.Create()

	expectedBody := "another body"
	patchNote := &models.Note{
		Body: expectedBody,
	}

	client := &http.Client{}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(patchNote)
	req, _ := http.NewRequest("PUT", fmt.Sprintf(n.ts.URL+"/api/me/notes/%d", note.ID), b)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))
	n.Require().True(rd.Success, rd.Message)

	models.GetDB().Take(note)
	n.Require().Equal(expectedBody, note.Body, "body should change")
	n.Require().Equal(expectedTitle, note.Title, "title should not change")
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
	n.Assert().Equal("<hashed>", rd.User.Password)
}

func (n *NotesTestSuite) TestPagination() {
	user := CreateUserTest()
	for cnt := 0; cnt < Cfg.PerPage+1; cnt++ {
		models.GetDB().Create(&models.Note{UserID: user.ID, Title: fmt.Sprintf("title %d", cnt)})
	}
	client := &http.Client{}

	req, _ := http.NewRequest("GET", n.ts.URL+"/api/me/notes?page=2", nil)
	AuthorizeRequest(req, user)

	resp := Must(client.Do(req))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))
	n.Require().True(rd.Success, rd.Message)
	n.Require().NotEmpty(rd.Pagination, "this response should have pagination info")
	n.Require().Equal(2, rd.Pagination["max_page"])

	actualIDs := []uint{}
	for _, n := range rd.Notes {
		actualIDs = append(actualIDs, n.ID)
	}

	n.Require().ElementsMatch(actualIDs, []uint{1})
}

func (n *NotesTestSuite) TestNotePublishedDetail() {
	user := CreateUserTest()
	note := &models.Note{
		Title:     "not matters",
		Body:      "not empty",
		UserID:    user.ID,
		Published: true,
	}
	note.Create()

	resp := Must(http.Get(fmt.Sprintf(n.ts.URL+"/api/notes/%v", note.ID)))
	n.Require().Equal(200, resp.StatusCode)

	resp = Must(http.Get(fmt.Sprintf(n.ts.URL+"/api/notes/%v?no_body=true", note.ID)))
	n.Require().Equal(200, resp.StatusCode)
	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))
	n.Require().True(rd.Success, rd.Message)
	n.Assert().Empty(rd.Note.Body)

	note.Published = false
	note.Save()

	resp = Must(http.Get(fmt.Sprintf(n.ts.URL+"/api/notes/%v", note.ID)))
	n.Require().Equal(404, resp.StatusCode)
}

func (n *NotesTestSuite) TestUnauth() {
	qs := []struct {
		method, url string
	}{
		{"GET", "/api/me/notes"},
		{"POST", "/api/me/notes"},
		{"GET", "/api/me/notes/1"},
		{"DELETE", "/api/me/notes/1"},
		{"PUT", "/api/me/notes/1"},
		{"GET", "/api/me"},
	}

	client := &http.Client{}
	for _, x := range qs {
		req, _ := http.NewRequest(x.method, n.ts.URL+x.url, nil)
		resp := Must(client.Do(req))
		n.Assert().Equal(403, resp.StatusCode)
	}
}

func (n *NotesTestSuite) TestUnauthContentType() {
	resp := Must(http.Get(n.ts.URL + "/api/me"))
	n.Require().Equal("application/json", resp.Header.Get("Content-Type"))
}

func (n *NotesTestSuite) TestBrokenDB() {
	user := CreateUserTest()
	note := &models.Note{
		Title:     "not matters",
		UserID:    user.ID,
		Published: true,
	}
	note.Create()

	// break db
	db, err := models.GetDB().DB()
	n.Require().Nil(err)
	db.Close()

	resp := Must(http.Get(fmt.Sprintf(n.ts.URL+"/api/notes/%v", note.ID)))
	n.Require().Equal(500, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd), "server should serve with valid json anyway")
	n.Require().False(rd.Success)
	n.Require().NotEmpty(rd.Message)
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

func (n *NotesTestSuite) TestAnotherUserDetail() {
	user := CreateUserTest()

	resp := Must(http.Get(fmt.Sprintf(n.ts.URL+"/api/users/%d", user.ID)))
	n.Require().Equal(200, resp.StatusCode)

	rd := &ResponseData{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd))

	n.Require().Equal(user.Username, rd.User.Username)
	n.Require().Empty(rd.User.Password, "password should not be sent")
}

type Object map[string]interface{}

type ResponseData struct {
	Success     bool           `json:"success"`
	Message     string         `json:"message"`
	Notes       []models.Note  `json:"notes"`
	AccessToken string         `json:"access_token"`
	Note        models.Note    `json:"note"`
	User        models.User    `json:"user"`
	Pagination  map[string]int `json:"pagination"`
}
