package users

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-rest-framework/core"
)

type UserData struct {
	Errors []core.ErrorMsg `json:"errors"`
	Data   User            `json:"data"`
}

func (u UserData) Read(r *http.Response) UserData {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	defer r.Body.Close()
	return u
}
