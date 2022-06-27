package users

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-rest-framework/core"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type UserKeyword struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
}

type UserKeywordData struct {
	Errors []core.ErrorMsg `json:"errors"`
	Data   UserKeyword     `json:"data"`
}

type UserKeywordsData struct {
	Errors []core.ErrorMsg `json:"errors"`
	Data   []UserKeyword   `json:"data"`
}

func (u *UserKeywordData) Read(r *http.Response) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	defer r.Body.Close()
}

func (u *UserKeywordsData) Read(r *http.Response) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	defer r.Body.Close()
}

func actionKeywordCreate(w http.ResponseWriter, r *http.Request) {
	var (
		model UserKeyword
		rsp   = core.Response{Data: &model, Req: r}
	)

	govalidator.TagMap["unique"] = govalidator.Validator(func(str string) bool {
		App.DB.Where("name = ?", str).First(&model)
		if model.ID != 0 {
			return false
		}
		return true
	})

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			App.DB.Create(&model)
		}
	}

	rsp.Data = &model

	w.Write(rsp.Make())
}

func actionKeywordUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		model UserKeyword
		data  UserKeyword
		rsp   = core.Response{Data: &model, Req: r}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {

			vars := mux.Vars(r)
			App.DB.First(&model, vars["id"])

			if model.ID == 0 {
				rsp.Errors.Add("ID", "Keyword not found")
			} else {
				App.DB.Model(&model).Updates(data)
			}
		}
	}

	rsp.Data = &model

	w.Write(rsp.Make())
}

func actionKeywordDelete(w http.ResponseWriter, r *http.Request) {
	var (
		model UserKeyword
		rsp   = core.Response{Data: &model, Req: r}
	)

	vars := mux.Vars(r)
	App.DB.First(&model, vars["id"])

	if model.ID == 0 {
		rsp.Errors.Add("ID", "Keyword not found")
	} else {
		App.DB.Unscoped().Delete(&model)
	}

	rsp.Data = &model

	w.Write(rsp.Make())
}

func actionKeywordGetAll(w http.ResponseWriter, r *http.Request) {
	var (
		models []UserKeyword
		count  int64
		rsp    = core.Response{Data: &models, Req: r}
		all    = r.FormValue("all")
		sort   = r.FormValue("sort")
		limit  = r.FormValue("limit")
		offset = r.FormValue("offset")
		db     = App.DB
	)

	if all != "" {
		db = db.Where("name LIKE ?", "%"+all+"%")
		db = db.Or("description LIKE ?", "%"+all+"%")
	}

	if sort != "" {
		switch sort {
		case "id":
			db = db.Order("id")
		case "-id":
			db = db.Order("id DESC")
		case "name":
			db = db.Order("name")
		case "-name":
			db = db.Order("name DESC")
		}
	} else {
		db = db.Order("id DESC")
	}

	db.Find(&models).Count(&count)

	if limit != "" {
		db = db.Limit(limit)
	} else {
		db = db.Limit(5)
	}

	if offset != "" {
		db = db.Offset(offset)
	}

	db.Find(&models)

	rsp.Data = &models
	rsp.Count = count

	w.Write(rsp.Make())
}
