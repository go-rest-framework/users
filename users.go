package users

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-rest-framework/core"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var App core.App

type Users []User

type User struct {
	gorm.Model
	Email       string `gorm:"unique;not null" valid:"email,required,unique~Email not unique"`
	Password    string `valid:"ascii,required"`
	Role        string
	Token       string
	Salt        string `json:"-"`
	CheckToken  string `json:"-"`
	CallBackUrl string `gorm:"-"`
}

type UserUpdate struct {
	ID       uint
	Email    string `valid:"email"`
	Password string `valid:"ascii"`
	Role     string
}

type Login struct {
	Email    string `valid:"email,required"`
	Password string `valid:"ascii,required"`
}

type Confirm struct {
	CheckToken string `valid:"required"`
}

type ResetRequest struct {
	Email       string `valid:"email,required"`
	CallBackUrl string
}

type Reset struct {
	CheckToken string `valid:"required"`
	Newpass    string `valid:"required"`
	NewpassRe  string `valid:"required"`
}

func Configure(a core.App) {
	App = a

	App.DB.AutoMigrate(&User{})

	createAdmin()

	//public actions
	App.R.HandleFunc("/api/users/register", srvRegister).Methods("POST")
	App.R.HandleFunc("/api/users/login", srvLogin).Methods("POST")
	App.R.HandleFunc("/api/users/confirm", srvConfirm).Methods("POST")
	App.R.HandleFunc("/api/users/resetrequest", srvResetrequest).Methods("POST")
	App.R.HandleFunc("/api/users/reset", srvReset).Methods("POST")

	//protect actions
	App.R.HandleFunc("/api/users/get-all", App.Protect(srvGetAll, []string{"admin"})).Methods("GET")
	App.R.HandleFunc("/api/users/get-one/{id}", App.Protect(srvGetOne, []string{"admin"})).Methods("GET")
	App.R.HandleFunc("/api/users/create", App.Protect(srvCreate, []string{"admin"})).Methods("POST")
	App.R.HandleFunc("/api/users/update", App.Protect(srvUpdate, []string{"admin"})).Methods("POST")
	App.R.HandleFunc("/api/users/delete", App.Protect(srvDelete, []string{"admin"})).Methods("POST")
}

func srvGetOne(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		rsp  = core.Response{Data: &user}
	)

	vars := mux.Vars(r)
	App.DB.First(&user, vars["id"])

	if user.ID == 0 {
		rsp.Errors.Add("ID", "User not found")
	} else {
		rsp.Data = &user
	}

	w.Write(rsp.Make())
}

func srvGetAll(w http.ResponseWriter, r *http.Request) {
	var (
		users Users
		rsp   = core.Response{Data: &users}
	)

	App.DB.Find(&users)

	rsp.Data = &users

	w.Write(rsp.Make())
}

func srvCreate(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		rsp  = core.Response{Data: &user}
	)

	govalidator.TagMap["unique"] = govalidator.Validator(func(str string) bool {
		App.DB.Where("email = ?", str).First(&user)
		if user.ID != 0 {
			return false
		}
		return true
	})

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			curtime := fmt.Sprintf("%x", time.Now())
			passsalt := App.ToSum256(curtime)
			passhash := App.ToSum256(user.Password + passsalt)
			user.Password = passhash
			user.Salt = passsalt
			App.DB.Create(&user)
		}
	}

	rsp.Data = &user

	w.Write(rsp.Make())
}

func srvUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		data UserUpdate
		user User
		rsp  = core.Response{Data: &data}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			if data.ID == 0 {
				rsp.Errors.Add("ID", "ID not set")
			} else {
				App.DB.First(&user, data.ID)
				if user.ID == 0 {
					rsp.Errors.Add("ID", "User not found")
				} else {
					App.DB.Model(&user).Updates(data)
				}
			}
		}
	}

	rsp.Data = &user

	w.Write(rsp.Make())
}

func srvDelete(w http.ResponseWriter, r *http.Request) {
	var (
		data User
		user User
		rsp  = core.Response{Data: &data}
	)

	if rsp.IsJsonParseDone(r.Body) {
		App.DB.First(&user, data.ID)
		if user.ID == 0 {
			rsp.Errors.Add("ID", "User not found")
		} else {
			if App.IsTest {
				App.DB.Unscoped().Delete(&user)
			} else {
				App.DB.Delete(&user)
			}
		}
	}

	rsp.Data = &user

	w.Write(rsp.Make())
}

func srvLogin(w http.ResponseWriter, r *http.Request) {
	var (
		data Login
		user User
		rsp  = core.Response{Data: &data}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			App.DB.Where("email = ?", data.Email).First(&user)
			if user.ID == 0 || user.Password != App.ToSum256(data.Password+user.Salt) {
				rsp.Errors.Add("Email", "User not found or wrong password")
			} else if user.Role == "" || user.Role == "candidate" {
				rsp.Errors.Add("Email", "You have not verified your email address")
			} else {
				token, err := App.GenToken(&user.Email, &user.Role)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					rsp.Errors.Add("Email", "Error generating JWT token: "+err.Error())
				} else {
					w.Header().Set("Authorization", "Bearer "+token)
					w.WriteHeader(http.StatusOK)
					user.Token = token
				}
			}
		}
	}

	user.Password = ""
	rsp.Data = &user

	w.Write(rsp.Make())
}

func srvRegister(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		rsp  = core.Response{Data: &user}
	)

	govalidator.TagMap["unique"] = govalidator.Validator(func(str string) bool {
		App.DB.Where("email = ?", str).First(&user)
		if user.ID != 0 {
			return false
		}
		return true
	})

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			var checktoken string
			curtime := fmt.Sprintf("%x", time.Now())
			passsalt := App.ToSum256(curtime)
			passhash := App.ToSum256(user.Password + passsalt)
			if App.IsTest {
				checktoken = "testchecktoken"
			} else {
				checktoken = App.ToSum256(fmt.Sprintf("%s.%x", user.Email, time.Now()))
			}
			user.Password = passhash
			user.Salt = passsalt
			user.Role = "candidate"
			user.CheckToken = checktoken
			log.Println(checktoken)
			App.DB.Create(&user)
			App.Mail.Send(
				user.Email,
				"Registration confirm",
				"To confirm the registration, go to the link "+user.CallBackUrl+"?token="+checktoken,
			)
		}
	}

	user.Password = ""
	rsp.Data = &user

	w.Write(rsp.Make())
}

func srvConfirm(w http.ResponseWriter, r *http.Request) {
	var (
		data Confirm
		user User
		rsp  = core.Response{Data: &data}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			App.DB.Where("check_token = ?", data.CheckToken).First(&user)
			if user.ID == 0 {
				rsp.Errors.Add("CheckToken", "User not found")
			} else if user.Role != "" && user.Role != "candidate" {
				rsp.Errors.Add("CheckToken", "You have already verified your email")
			} else {
				res := App.DB.Model(&user).Updates(map[string]interface{}{"role": "user", "check_token": ""})
				if res.Error != nil {
					w.WriteHeader(http.StatusInternalServerError)
					rsp.Errors.Add("CheckToken", "Data saving error")
					log.Println("Data saving error: " + res.Error.Error())
				}
			}
		}
	}

	user.Password = ""
	rsp.Data = &user

	w.Write(rsp.Make())
}

func srvResetrequest(w http.ResponseWriter, r *http.Request) {
	var (
		data ResetRequest
		user User
		rsp  = core.Response{Data: &data}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			var checktoken string
			App.DB.Where("email = ?", data.Email).First(&user)
			if user.ID == 0 {
				rsp.Errors.Add("Email", "User not found")
			} else if user.Role == "" || user.Role == "candidate" {
				rsp.Errors.Add("Email", "You have not confirmed your email")
			} else {
				if App.IsTest {
					checktoken = "testchecktoken"
				} else {
					checktoken = App.ToSum256(fmt.Sprintf("%s.%x", user.Email, time.Now()))
				}
				res := App.DB.Model(&user).Update("checktoken", checktoken)
				if res.Error != nil {
					w.WriteHeader(http.StatusInternalServerError)
					rsp.Errors.Add("Email", "Data saving error")
					log.Println("Data saving error: " + res.Error.Error())
				} else {
					App.Mail.Send(
						user.Email,
						"Password reset request",
						"To reset your password, go to the link "+data.CallBackUrl+"?token="+checktoken,
					)
					log.Println(checktoken)
				}
			}
		}
	}

	rsp.Data = &user

	w.Write(rsp.Make())
}

func srvReset(w http.ResponseWriter, r *http.Request) {
	var (
		data Reset
		user User
		rsp  = core.Response{Data: &data}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			App.DB.Where("check_token = ?", data.CheckToken).First(&user)
			if user.ID == 0 {
				rsp.Errors.Add("Email", "User not found")
			} else if user.Role == "" || user.Role == "candidate" {
				rsp.Errors.Add("Email", "You have already verified your email")
			} else if &data.Newpass != &data.NewpassRe {
				rsp.Errors.Add("Newpass", "New password and new password repeat must be equal")
			} else {
				passhash := App.ToSum256(data.Newpass + user.Salt)
				App.DB.Model(&user).Updates(map[string]interface{}{"password": passhash, "check_token": ""})
			}
		}
	}

	rsp.Data = &user

	w.Write(rsp.Make())
}

func createAdmin() {
	var (
		user User
	)

	user.Email = "admin@admin.a"

	App.DB.Where("email = ?", user.Email).First(&user)
	if user.ID == 0 {
		curtime := fmt.Sprintf("%x", time.Now())
		passsalt := App.ToSum256(curtime)

		if App.IsTest {
			user.Password = "adminpass"
		} else {
			user.Password = curtime[1:16]
		}

		fmt.Printf("%s\n", user.Password)

		passhash := App.ToSum256(user.Password + passsalt)

		user.Password = passhash
		user.Salt = passsalt
		user.Role = "admin"
		App.DB.Create(&user)
	}
}
