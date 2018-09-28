package users_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/go-rest-framework/core"
	"github.com/go-rest-framework/users"
	"github.com/icrowley/fake"
)

var Uid uint
var Uidnew uint
var UEmail string
var UNewEmail string
var AdminToken string
var Murl = "http://gorest.ga/api/users/"

type TestUsers struct {
	Errors []core.ErrorMsg
	Data   users.Users
}

type TestUser struct {
	Errors []core.ErrorMsg
	Data   users.User
}

func readUsersBody(r *http.Response, t *testing.T) TestUsers {
	var u TestUsers
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	return u
}

func doRequest(url, proto, userJson, token string) *http.Response {
	reader := strings.NewReader(userJson)
	request, err := http.NewRequest(proto, url, reader)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func readUserBody(r *http.Response, t *testing.T) TestUser {
	var u TestUser
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	defer r.Body.Close()
	return u
}

//"/api/users/register", srvRegister).Methods("POST")
func TestRegister(t *testing.T) {

	UEmail = fake.EmailAddress()

	url := Murl + "register"
	userJson := `{"Email":"sldfjsdlfeusdlfjsdlfj", "Password":"343223423423"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("email type validation dont work")
	}

	userJson = `{"Email":"sdlfjldjflsdf@sldfjsdlf.eu"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("require validation dont work")
	}

	userJson = `{"Email":"` + UEmail + `", "Password":"343223423423"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	Uid = u.Data.ID

	return
}

//"/api/users/confirm", srvConfirm).Methods("POST")
func TestConfirm(t *testing.T) {

	url := Murl + "confirm"
	var userJson = `{"CheckToken":"wrongtoken"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("token check fail")
	}

	userJson = `{}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("require validation dont work")
	}

	userJson = `{"CheckToken":"testchecktoken"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

//"/api/users/login", srvLogin).Methods("POST")
func TestLogin(t *testing.T) {

	url := Murl + "login"
	var userJson = `{"Email":"sdlf@eusdlfjsdlfj.com", "Password":"dddd343223423423"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("password check fail")
	}

	userJson = `{"Email":"` + UEmail + `"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("require validation dont work")
	}

	userJson = `{"Email":"` + UEmail + `", "Password":"343223423423"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

//"/api/users/resetrequest", srvResetrequest).Methods("POST")
func TestResetrequest(t *testing.T) {

	url := Murl + "resetrequest"
	var userJson = `{"Email":"` + UEmail + `", "CallBackUrl":"http://test.ttt"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	return
}

//"/api/users/reset", srvReset).Methods("POST")
func TestReset(t *testing.T) {

	url := Murl + "reset"
	var userJson = `{
		"CheckToken":"testchecktoken",
		"Newpass":"newpass",
		"NewpassRe":"newpass1"
	}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("check equal passwords fail")
	}

	userJson = `{
		"CheckToken":"testchecktoken",
		"Newpass":"newpass",
		"NewpassRe":"newpass"
	}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	return
}

func TestAdminLogin(t *testing.T) {

	url := Murl + "login"
	var userJson = `{"Email":"admin@admin.a", "Password":"adminpass"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	AdminToken = u.Data.Token

	return
}

//"/api/users/create", App.Protect(srvCreate, []string{"admin"})).Methods("POST")
func TestCreate(t *testing.T) {
	UNewEmail = fake.EmailAddress()
	url := Murl + "create"
	userJson := `{
		"Email":"` + UNewEmail + `",
		"Password":"newuserpass",
		"Role":"user"
	}`

	resp := doRequest(url, "GET", userJson, AdminToken)

	if resp.StatusCode == 200 {
		t.Fatal("POST check dont work")
	}

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	Uidnew = u.Data.ID

	return
}

//"/api/users/get-all", App.Protect(srvGetAll, []string{"admin"})).Methods("GET")
func TestGetAll(t *testing.T) {
	// get count
	url := Murl + "get-all"

	resp := doRequest(url, "GET", "", "  ")

	if resp.StatusCode == 200 {
		t.Fatal("require validation dont work")
	}

	resp = doRequest(url, "GET", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUsersBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

//"/api/users/get-one/{id}", App.Protect(srvGetOne, []string{"admin"})).Methods("GET")
func TestGetOne(t *testing.T) {
	url := Murl + "get-one/" + "0"
	resp := doRequest(url, "GET", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("element not found dont work")
	}

	url = fmt.Sprintf("%s%s%d", Murl, "get-one/", Uid)

	resp = doRequest(url, "GET", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if u.Data.Email != UEmail {
		t.Fatal("wrong email get")
	}

	return
}

//"/api/users/update", App.Protect(srvUpdate, []string{"admin"})).Methods("POST")
func TestUpdate(t *testing.T) {
	url := Murl + "update"
	userJson := `{"Email":"sldfjsdlfeusdlfjsdlfj", "Password":"343223423423"}`

	resp := doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("email type validation dont work")
	}

	userJson = `{"Role":"admin"}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("require validation dont work")
	}

	userJson = fmt.Sprintf("%s%d%s", `{"ID":`, Uidnew, `,"Role":"admin"}`)

	fmt.Printf("%s\n", userJson)

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

//"/api/users/delete", App.Protect(srvDelete, []string{"admin"})).Methods("POST")
func TestDelete(t *testing.T) {
	url := Murl + "delete"
	userJson := `{"ID":"0"}`

	resp := doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("wrong id validation dont work")
	}

	userJson = fmt.Sprintf("%s%d%s", `{"ID":`, Uid, `}`)

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	userJson = fmt.Sprintf("%s%d%s", `{"ID":`, Uidnew, `}`)

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}
