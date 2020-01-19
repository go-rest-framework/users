//TODO
// - add test for login not active user
// - add test for actions with not active user
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
var UProfileLastname string
var UProfilePhone string
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

	userJson = `{
		"email":"` + UEmail + `",
		"password":"aaAA11..",
		"repassword":"aaAA11.."
	}`

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
func TestConfirmEmail(t *testing.T) {

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
	var userJson string
	url := Murl + "/login"

	userJson = `{"email":"admin@admin.a", "password":"wrongpass"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("No error with wrong password")
	}

	userJson = `{"email":"admin@admin.a", "password":"adminpass"}`

	resp = doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	AdminToken = u.Data.Token

	return
}

//"/api/users/create", App.Protect(srvCreate, []string{"admin"})).Methods("POST")
func TestCreate(t *testing.T) {
	var userJson string
	UNewEmail = fake.EmailAddress()
	UProfileName = fake.FirstName()
	UProfileLastname = fake.LastName()
	UProfilePhone = fake.Phone()
	url := Murl
	//negative empty role
	userJson = `{
		"email":"` + UNewEmail + `",
		"password":"good.PASS123",
		"repassword":"good.PASS123",
		"role":"wrongrole",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + fake.LastName() + `",
			"phone":"` + fake.Phone() + `",
			"avatar":"00001100"
		}
	}`

	resp := doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("no error if wrong role")
	}
	//negative wrong email format
	userJson = `{
		"email":"sdlfjsdflsdfsdfsdf",
		"password":"good.PASS123",
		"repassword":"good.PASS123",
		"role":"user",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + fake.LastName() + `",
			"phone":"` + fake.Phone() + `",
			"avatar":"00001100"
		}
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("no error if wrong email format")
	}
	//negative no repass
	userJson = `{
		"email":"` + UNewEmail + `",
		"password":"good.PASS123",
		"repassword":"",
		"role":"user",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + fake.LastName() + `",
			"phone":"` + fake.Phone() + `",
			"avatar":"00001100"
		}
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)
	if len(u.Errors) == 0 {
		t.Fatal("no error if no repass")
	}
	//negative low password complisity
	userJson = `{
		"email":"` + UNewEmail + `",
		"password":"PASS123",
		"repassword":"PASS123",
		"role":"user",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + fake.LastName() + `",
			"phone":"` + fake.Phone() + `",
			"avatar":"00001100"
		}
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)
	if len(u.Errors) == 0 {
		t.Fatal("no error if low pass complisity")
	}

	//negative no status
	userJson = `{
		"email":"` + UNewEmail + `",
		"password":"good.PASS123",
		"repassword":"good.PASS123",
		"role":"user",
		"status":"wrongstatus",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + fake.LastName() + `",
			"phone":"` + fake.Phone() + `",
			"avatar":"00001100"
		}
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)
	if len(u.Errors) == 0 {
		t.Fatal("no error if wrong status")
	}

	//negative status not from ["active","blocked","draft"]
	userJson = `{
		"email":"` + UNewEmail + `",
		"password":"good.PASS123",
		"repassword":"good.PASS123",
		"role":"user",
		"status":"wrongstatus",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + fake.LastName() + `",
			"phone":"` + fake.Phone() + `",
			"avatar":"00001100"
		}
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)
	if len(u.Errors) == 0 {
		t.Fatal("no error if status not in list active, blocked, draft")
	}

	//positive
	userJson = `{
		"email":"` + UNewEmail + `",
		"password":"good.PASS123",
		"repassword":"good.PASS123",
		"role":"user",
		"status":"active",
		"profile":{
			"firstname":"` + UProfileName + `",
			"middlename":"` + fake.FirstName() + `",
			"lastname":"` + UProfileLastname + `",
			"phone":"` + UProfilePhone + `",
			"avatar":"00001100"
		}
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

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

func doOneSearch(url string, t *testing.T) TestUsers {
	resp := doRequest(url, "GET", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	return readUsersBody(resp, t)
}

//"/api/users/get-all", App.Protect(srvGetAll, []string{"admin"})).Methods("GET")
func TestGetAll(t *testing.T) {
	var u TestUsers
	url := Murl
	resp := doRequest(url, "GET", "", "  ")
	if resp.StatusCode == 200 {
		t.Fatal("require auntifications dont work")
	}
	//positive get all users
	u = doOneSearch(url, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}
	//positive search all email
	u = doOneSearch(Murl+"?all="+UNewEmail, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search all role
	u = doOneSearch(Murl+"?all=user", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 3 {
		t.Errorf("Expected 3 element, giwen - : %d", len(u.Data))
	}
	//positive search all status
	u = doOneSearch(Murl+"?all=blocked", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 0 {
		t.Errorf("Expected 0 elements, giwen - : %d", len(u.Data))
	}
	//positive search all firstname
	u = doOneSearch(Murl+"?all="+UProfileName, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search all phone
	u = doOneSearch(Murl+"?all="+UProfilePhone, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search by filter id
	u = doOneSearch(Murl+"?id="+fmt.Sprintf("%d", Uidnew), t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search by filter email
	u = doOneSearch(Murl+"?email="+UNewEmail, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search by filter role
	u = doOneSearch(Murl+"?role=user", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 3 {
		t.Errorf("Expected 3 element, giwen - : %d", len(u.Data))
	}
	//positive search by filter status
	u = doOneSearch(Murl+"?status=blocked", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 0 {
		t.Errorf("Expected 0 elements, giwen - : %d", len(u.Data))
	}
	//positive search by filter name with firstname
	u = doOneSearch(Murl+"?name="+UProfileName, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search by filter name with firstname + lastname
	u = doOneSearch(Murl+"?name="+UProfileName+"+"+UProfileLastname, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search by filter phone
	u = doOneSearch(Murl+"?phone="+UProfilePhone, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive search by filter with firstname and status
	u = doOneSearch(Murl+"?name="+UProfileName+"&status=active", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 1 {
		t.Errorf("Expected one element, giwen - : %d", len(u.Data))
	}
	//positive sort by id
	u = doOneSearch(Murl+"?sort=id", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if !(u.Data[0].ID < u.Data[1].ID && u.Data[1].ID < u.Data[2].ID) {
		t.Fatal("sorting id dont work")
	}
	//positive sort by id DESC
	u = doOneSearch(Murl+"?sort=-id", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if !(u.Data[0].ID > u.Data[1].ID && u.Data[1].ID > u.Data[2].ID) {
		t.Fatal("sorting id DESC dont work")
	}
	//positive sort by email
	u = doOneSearch(Murl+"?sort=email", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	for k, v := range u.Data {
		fmt.Println(k, v.Email)
	}

	u = doOneSearch(Murl+"?sort=-email", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	for k, v := range u.Data {
		fmt.Println(k, v.Email)
	}
	//positive sort by phone
	u = doOneSearch(Murl+"?sort=-phone", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if u.Data[0].Profile.Phone == "" {
		t.Fatal("sorting phone dont work")
	}
	//positive sort by name
	u = doOneSearch(Murl+"?sort=-name", t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if u.Data[0].Email == "" {
		t.Fatal("sorting email dont work")
	}
	// get count

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
	url := fmt.Sprintf("%s%s%d", Murl, "/", Uidnew)
	userJson := `{
		"status":"blocked",
		"profile":{
			"firstname": "test111",
			"middlename": "test222",
			"lastname": "test333",
			"phone": "12345",
			"avatar": ""
		}
	}`

	resp := doRequest(url, "PATCH", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if u.Data.Status != "blocked" {
		t.Fatal("update dont work", u.Errors, Uidnew)
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

	/*url = fmt.Sprintf("%s%s%d", Murl, "/", Uid)

	resp = doRequest(url, "DELETE", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}*/

	fmt.Println(Uidnew)

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
