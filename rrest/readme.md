


```
package course

import (
	"fmt"
	"net/http"
	"roo.bo/appcourse/api/rres"
	"roo.bo/rlib/base"
)

type Course struct {
	ID             int    `db:"id" json:"id"`
	AppID          string `db:"appid" json:"appid"`
	Name           string `db:"name" json:"name" validate:"required"`
	Intro          string `db:"intro" json:"intro"`
	Icon           string `db:"icon" json:"icon" validate:"required"`
	LastModifiedBy string `db:"last_modified_by" json:"lastModifiedBy"`
	IsDeleted      int    `db:"is_deleted" json:"isDeleted"`
	CreatedAt      int    `db:"created_at" json:"createdAt"`
	UpdatedAt      int    `db:"updated_at" json:"updatedAt"`
	CopyFrom       int    `db:"copy_from" json:"copyFrom"`
	BatchNo        string `db:"batchNo" json:"batchNo"`
	Language       string `db:"language" json:"language"`
	IsPublic       int    `db:"is_public" json:"isPublic"`
}

func (*Course) PK() string {
	return "id"
}

func (*Course) TableName() string {
	return "course"
}

func List(w http.ResponseWriter, r *http.Request) {
	rres.List{
		Filters: []string{
			"appid__eq",
		},
		Force: map[string]string{
			"appid": "current appid",
		},
		Model: &[]Course{},
		R:     r,
		W:     w,
	}.Parse()
}

func Create(w http.ResponseWriter, r *http.Request) {
	rres.Create{
		Force: map[string]string{
			"last_modified_by": "current username",
		},
		Validate: map[string]rres.Validator{
			"language": func(k string, form *base.Struct) (string, bool) {
				v, _ := form.GetString(k)
				if v != "en" && v != "cn" {
					return "language only be en or cn", false
				}
				return "", true
			},
		},
		Model: &Course{},
		R:     r,
		W:     w,
	}.Parse()
}

func Detail(w http.ResponseWriter, r *http.Request) {
	rres.Detail{
		Model: &Course{},
		R:     r,
		W:     w,
	}.Parse()
}

func Delete(w http.ResponseWriter, r *http.Request) {
	rres.Delete{
		Force: map[string]string{
			"appid": "could only delete current appid",
		},
		Model: &Course{},
		R:     r,
		W:     w,
	}.Parse()
}

func Update(w http.ResponseWriter, r *http.Request) {
	rres.Update{
		Force: map[string]string{
			"last_modified_by": "current username",
		},
		Validate: map[string]rres.Validator{
			"language": func(k string, form *base.Struct) (string, bool) {
				v, _ := form.GetString(k)
				fmt.Printf("%+v", v)
				if v != "en" && v != "cn" {
					return "language only be en or cn", false
				}
				return "", true
			},
		},
		Model: &Course{},
		R:     r,
		W:     w,
	}.Parse()
}
```
