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
	Email       string  `json:"email" gorm:"unique;not null" valid:"email,required,unique~email: Email not unique"`
	Password    string  `json:"password" valid:"ascii,required"`
	RePassword  string  `gorm:"-" json:"repassword" valid:"ascii,required,passmatch~repassword: Passwords do not match"`
	Role        string  `json:"role" valid:"required"`
	Token       string  `json:"token"`
	Salt        string  `json:"-"`
	CheckToken  string  `json:"-"`
	CallBackUrl string  `gorm:"-"`
	Profile     Profile `json:"profile"`
	ProfileID   int     `json:"profileID"`
}

type UserUpdate struct {
	Email   string  `json:"email" valid:"email"`
	Role    string  `json:"role"`
	Profile Profile `json:"profile"`
}

type Login struct {
	Email    string `valid:"email,required" json:"email"`
	Password string `valid:"ascii,required" json:"password"`
}

type Confirm struct {
	CheckToken string `json:"checkToken" valid:"required"`
}

type ResetRequest struct {
	Email       string `json:"email" valid:"email,required"`
	CallBackUrl string `json:"callBackUrl"`
}

type Reset struct {
	CheckToken string `json:"checkToken" valid:"required"`
	Newpass    string `json:"newpass" valid:"required"`
	NewpassRe  string `json:"newpassRe" valid:"required"`
}

func init() {
	govalidator.CustomTypeTagMap.Set("passmatch", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
		switch v := context.(type) { // this validates a field against the value in another field, i.e. dependent validation
		case User:
			if i == v.Password {
				return true
			}
		}
		return false
	}))
}

func Configure(a core.App) {
	App = a

	App.DB.AutoMigrate(&User{}, &Profile{})

	createAdmin()
	createTestUser()

	//public actions
	App.R.HandleFunc("/users/register", actionRegister).Methods("POST")
	App.R.HandleFunc("/users/login", actionLogin).Methods("POST")
	App.R.HandleFunc("/users/confirm", actionConfirm).Methods("POST")
	App.R.HandleFunc("/users/resetrequest", actionResetrequest).Methods("POST")
	App.R.HandleFunc("/users/reset", actionReset).Methods("POST")

	App.R.HandleFunc("/users/{id}/profile", actionGetProfile).Methods("GET")

	//protect actions
	App.R.HandleFunc("/users", App.Protect(actionGetAll, []string{"admin"})).Methods("GET")
	App.R.HandleFunc("/users/{id}", App.Protect(actionGetOne, []string{"admin"})).Methods("GET")
	App.R.HandleFunc("/users", App.Protect(actionCreate, []string{"admin"})).Methods("POST")
	App.R.HandleFunc("/users/{id}", App.Protect(actionUpdate, []string{"admin"})).Methods("PATCH")
	App.R.HandleFunc("/users/{id}", App.Protect(actionDelete, []string{"admin"})).Methods("DELETE")
}

func actionGetOne(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		rsp  = core.Response{Data: &user, Req: r}
	)

	vars := mux.Vars(r)
	App.DB.Preload("Profile").First(&user, vars["id"])

	if user.ID == 0 {
		rsp.Errors.Add("ID", "User not found")
	} else {
		rsp.Data = &user
	}

	w.Write(rsp.Make())
}

func actionGetAll(w http.ResponseWriter, r *http.Request) {
	var (
		users Users
		rsp   = core.Response{Data: &users, Req: r}
		all   = r.FormValue("all")
		id    = r.FormValue("id")
		email = r.FormValue("email")
		role  = r.FormValue("role")
		sort  = r.FormValue("sort")
		db    = App.DB
	)

	if all != "" {
		db = db.Where("id LIKE ?", "%"+all+"%")
		db = db.Or("email LIKE ?", "%"+all+"%")
		db = db.Or("role LIKE ?", "%"+all+"%")
	}

	if id != "" {
		db = db.Where("id LIKE ?", "%"+id+"%")
	}

	if email != "" {
		db = db.Where("email LIKE ?", "%"+email+"%")
	}

	if role != "" {
		db = db.Where("role LIKE ?", "%"+role+"%")
	}

	if sort != "" {
		db = db.Order(sort)
	}

	db.Preload("Profile").Find(&users)

	rsp.Data = &users

	w.Write(rsp.Make())
}

func actionCreate(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		rsp  = core.Response{Data: &user, Req: r}
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

func actionUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		data UserUpdate
		user User
		rsp  = core.Response{Data: &data, Req: r}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {

			vars := mux.Vars(r)
			App.DB.First(&user, vars["id"])

			if user.ID == 0 {
				rsp.Errors.Add("ID", "User not found")
			} else {
				App.DB.Model(&user).Updates(data)
			}
		}
	}

	rsp.Data = &user

	w.Write(rsp.Make())
}

func actionDelete(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		rsp  = core.Response{Data: &user, Req: r}
	)

	vars := mux.Vars(r)
	App.DB.First(&user, vars["id"])

	if user.ID == 0 {
		rsp.Errors.Add("ID", "User not found")
	} else {
		if App.IsTest {
			App.DB.Unscoped().Delete(&user)
		} else {
			App.DB.Delete(&user)
		}
	}

	rsp.Data = &user

	w.Write(rsp.Make())
}

func actionLogin(w http.ResponseWriter, r *http.Request) {
	var (
		data Login
		user User
		rsp  = core.Response{Data: &data, Req: r}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			App.DB.Where("email = ?", data.Email).First(&user)
			if user.ID == 0 || user.Password != App.ToSum256(data.Password+user.Salt) {
				rsp.Errors.Add("email", "User not found or wrong password")
			} else if user.Role == "" || user.Role == "candidate" {
				rsp.Errors.Add("email", "You have not verified your email address")
			} else {
				idstring := fmt.Sprintf("%d", user.ID)
				token, err := App.GenToken(&idstring, &user.Email, &user.Role)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					rsp.Errors.Add("email", "Error generating JWT token: "+err.Error())
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

func actionRegister(w http.ResponseWriter, r *http.Request) {
	var (
		user User
		rsp  = core.Response{Data: &user, Req: r}
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

func actionConfirm(w http.ResponseWriter, r *http.Request) {
	var (
		data Confirm
		user User
		rsp  = core.Response{Data: &data, Req: r}
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

func actionResetrequest(w http.ResponseWriter, r *http.Request) {
	var (
		data ResetRequest
		user User
		rsp  = core.Response{Data: &data, Req: r}
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

func actionReset(w http.ResponseWriter, r *http.Request) {
	var (
		data Reset
		user User
		rsp  = core.Response{Data: &data, Req: r}
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

		fmt.Printf("admin password: %s\n", user.Password)

		passhash := App.ToSum256(user.Password + passsalt)

		user.Password = passhash
		user.Salt = passsalt
		user.Role = "admin"
		App.DB.Create(&user)
	}
}

func createTestUser() {
	var (
		user User
	)

	user.Email = "testuser@test.t"

	App.DB.Where("email = ?", user.Email).First(&user)
	if user.ID == 0 {
		curtime := fmt.Sprintf("%x", time.Now())
		passsalt := App.ToSum256(curtime)

		if App.IsTest {
			user.Password = "testpass"
		} else {
			user.Password = curtime[1:16]
		}

		fmt.Printf("testuser password: %s\n", user.Password)

		passhash := App.ToSum256(user.Password + passsalt)

		user.Password = passhash
		user.Salt = passsalt
		user.Role = "user"
		App.DB.Create(&user)
	}
}

func actionGetProfile(w http.ResponseWriter, r *http.Request) {
	var (
		user    User
		profile Profile
		rsp     = core.Response{Data: &profile, Req: r}
	)

	vars := mux.Vars(r)
	App.DB.First(&user, vars["id"]).Related(&profile)

	if user.ID == 0 {
		rsp.Errors.Add("ID", "User not found")
	} else {
		rsp.Data = &profile
	}

	w.Write(rsp.Make())
}
