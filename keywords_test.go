package users

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/icrowley/fake"
)

var (
	U               UserData
	TestKeywordID   uint
	TestKeywordName string
	Apiurl          = "http://localhost/api"
)

func TestGetAdminToken(t *testing.T) {

	url := "http://localhost/api/users/login"
	var userJson = `{"email":"admin@admin.a", "password":"adminpass"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	U.Read(resp)

	return
}

func Test_actionKeywordCreate(t *testing.T) {
	tests := []struct {
		name string
		args UserKeyword
	}{
		{
			"normal",
			UserKeyword{
				Name:        fake.Words(),
				Description: fake.Words(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data UserKeywordData
			iurl := Apiurl + "/users/keywords"

			userJson, err := json.MarshalIndent(tt.args, " ", " ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(iurl, U.Data.Token, string(userJson))

			resp := doRequest(iurl, "POST", string(userJson), U.Data.Token)

			if resp.StatusCode != 200 {
				t.Errorf("Wrong Response status = %s, want %v", resp.Status, 200)
			}

			data.Read(resp)

			if data.Data.ID == 0 {
				t.Errorf("Wrong ID not created")
			} else {
				TestKeywordID = data.Data.ID
				TestKeywordName = data.Data.Name
			}
		})
	}
}

func Test_actionKeywordUpdate(t *testing.T) {
	tests := []struct {
		name string
		args UserKeyword
	}{
		{
			"normal",
			UserKeyword{
				Description: fake.Words(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data UserKeywordData
			iurl := Apiurl + "/users/keywords/" + fmt.Sprintf("%d", TestKeywordID)
			iproto := "PATCH"

			userJson, err := json.MarshalIndent(tt.args, " ", " ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(iurl, U.Data.Token, string(userJson))

			resp := doRequest(iurl, iproto, string(userJson), U.Data.Token)

			if resp.StatusCode != 200 {
				t.Errorf("Wrong Response status = %s, want %v", resp.Status, 200)
			}

			data.Read(resp)

			fmt.Println("Description in updated model: ", data.Data.Description)
		})
	}
}

func Test_actionKeywordGetAll(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{
			"find all",
			"",
		},
		{
			"find from list by name",
			"&all=" + TestKeywordName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data UserKeywordsData
			iurl := Apiurl + "/users/keywords/" + tt.args
			iproto := "GET"

			resp := doRequest(iurl, iproto, "", U.Data.Token)

			if resp.StatusCode != 200 {
				t.Errorf("Wrong Response status = %s, want %v", resp.Status, 200)
			}

			data.Read(resp)

			if len(data.Data) == 0 {
				t.Errorf("No data in list")
			}
		})
	}
}

func Test_actionKeywordDelete(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data UserKeywordData
			iurl := Apiurl + "/users/keywords/" + fmt.Sprintf("%d", TestKeywordID)
			iproto := "DELETE"

			resp := doRequest(iurl, iproto, "", U.Data.Token)

			if resp.StatusCode != 200 {
				t.Errorf("Wrong Response status = %s, want %v", resp.Status, 200)
			}

			data.Read(resp)
		})
	}
}
