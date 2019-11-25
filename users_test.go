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
var UProfileName string
var Murl = "http://localhost/api/users"

type TestUsers struct {
	Errors []core.ErrorMsg `json:"errors"`
	Data   users.Users     `json:"data"`
}

type TestUser struct {
	Errors []core.ErrorMsg `json:"errors"`
	Data   users.User      `json:"data"`
}

type TestProfile struct {
	Errors []core.ErrorMsg `json:"errors"`
	Data   users.Profile   `json:"data"`
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

func readProfileBody(r *http.Response, t *testing.T) TestProfile {
	var p TestProfile
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &p)
	defer r.Body.Close()
	return p
}

//"/api/users/register", srvRegister).Methods("POST")
func TestRegister(t *testing.T) {

	UEmail = fake.EmailAddress()

	url := Murl + "/register"
	userJson := `{"email":"sldfjsdlfeusdlfjsdlfj", "password":"343223423423"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("email type validation dont work")
	}

	userJson = `{"email":"sdlfjldjflsdf@sldfjsdlf.eu"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("require validation dont work")
	}

	userJson = `{"email":"` + UEmail + `", "password":"343223423423"}`

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

	url := Murl + "/confirm"
	var userJson = `{"checkToken":"wrongtoken"}`

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

	userJson = `{"checkToken":"testchecktoken"}`

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

	url := Murl + "/login"
	var userJson = `{"email":"sdlf@eusdlfjsdlfj.com", "password":"dddd343223423423"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("password check fail")
	}

	userJson = `{"email":"` + UEmail + `"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("require validation dont work")
	}

	userJson = `{"email":"` + UEmail + `", "password":"343223423423"}`

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

	url := Murl + "/resetrequest"
	var userJson = `{"email":"` + UEmail + `", "callBackUrl":"http://test.ttt"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	return
}

//"/api/users/reset", srvReset).Methods("POST")
func TestReset(t *testing.T) {

	url := Murl + "/reset"
	var userJson = `{
		"checkToken":"testchecktoken",
		"newpass":"newpass",
		"newpassRe":"newpass1"
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
		"checkToken":"testchecktoken",
		"newpass":"newpass",
		"newpassRe":"newpass"
	}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	return
}

func TestAdminLogin(t *testing.T) {

	url := Murl + "/login"
	var userJson = `{"email":"admin@admin.a", "password":"adminpass"}`

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
	UProfileName = fake.FirstName()
	url := Murl
	userJson := `{
		"email":"` + UNewEmail + `",
		"password":"newuserpass",
		"role":"user",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + fake.LastName() + `",
			"phone":"` + fake.Phone() + `"
		}
	}`

	resp := doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	Uidnew = u.Data.ID

	fmt.Println(u.Data.Profile.Firstname)

	if u.Data.Profile.Firstname != UProfileName {
		t.Fatal("wrong user profile firstname")
	}

	return
}

//"/api/users/get-all", App.Protect(srvGetAll, []string{"admin"})).Methods("GET")
func TestGetAll(t *testing.T) {
	// get count
	url := Murl

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

	url1 := Murl + "?email=" + UNewEmail

	resp1 := doRequest(url1, "GET", "", AdminToken)

	if resp1.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp1.StatusCode)
	}

	u1 := readUsersBody(resp1, t)

	if len(u1.Errors) != 0 {
		t.Fatal(u1.Errors)
	}

	if len(u1.Data) != 1 {
		fmt.Println(u1.Data)
		t.Errorf("Expected one element, giwen - : %d", len(u1.Data))
	}

	//---------------

	url2 := Murl + "?sort=ID"

	resp2 := doRequest(url2, "GET", "", AdminToken)

	if resp2.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp2.StatusCode)
	}

	u2 := readUsersBody(resp2, t)

	if len(u2.Errors) != 0 {
		t.Fatal(u2.Errors)
	}

	if !(u2.Data[0].ID < u2.Data[1].ID && u2.Data[1].ID < u2.Data[2].ID) {
		t.Fatal("sorting dont work")
	}

	return
}

//"/api/users/get-one/{id}", App.Protect(srvGetOne, []string{"admin"})).Methods("GET")
func TestGetOne(t *testing.T) {
	url := Murl + "/0"
	resp := doRequest(url, "GET", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("element not found dont work")
	}

	url = fmt.Sprintf("%s%s%d", Murl, "/", Uidnew)

	resp = doRequest(url, "GET", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if u.Data.Email != UNewEmail {
		t.Fatal("wrong email get")
	}

	if u.Data.Profile.Firstname != UProfileName {
		t.Fatal("wrong user profile firstname")
	}

	return
}

func TestGetOneProfile(t *testing.T) {
	url := fmt.Sprintf("%s%s%d%s", Murl, "/", Uidnew, "/profile")

	resp := doRequest(url, "GET", "", "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	p := readProfileBody(resp, t)

	if len(p.Errors) != 0 {
		t.Fatal(p.Errors)
	}

	if p.Data.Firstname != UProfileName {
		t.Fatal("wrong user profile firstname")
	}

	return
}

//"/api/users/update", App.Protect(srvUpdate, []string{"admin"})).Methods("POST")
func TestUpdate(t *testing.T) {
	url := fmt.Sprintf("%s%s%d", Murl, "/", Uid)
	userJson := `{"email":"sldfjsdlfeusdlfjsdlfj", "password":"343223423423"}`

	resp := doRequest(url, "PATCH", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("email type validation dont work")
	}

	url = fmt.Sprintf("%s%s%d", Murl, "/", Uidnew)
	userJson = `{"role":"admin"}`

	resp = doRequest(url, "PATCH", userJson, AdminToken)

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
	url := fmt.Sprintf("%s%s%d", Murl, "/", 0)

	resp := doRequest(url, "DELETE", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("wrong id validation dont work")
	}

	url = fmt.Sprintf("%s%s%d", Murl, "/", Uid)

	resp = doRequest(url, "DELETE", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	url = fmt.Sprintf("%s%s%d", Murl, "/", Uidnew)

	resp = doRequest(url, "DELETE", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}
