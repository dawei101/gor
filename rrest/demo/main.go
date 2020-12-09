package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gorilla/mux"
	"roo.bo/rlib"
	"roo.bo/rlib/rrest"
	"roo.bo/rlib/rsql"
)

type User struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name" validate:"required,min=10"`
	CreatedAt string `db:"created_at" json:"createdAt"`
}

func (u *User) PK() string {
	return "id"
}

func (u *User) TableName() string {
	return "user"
}

type Profile struct {
	ID        int    `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"userID"`
	Nickaname string `db:"nickname" json:"nickname" validate:"required,max=10"`
	CreatedAt string `db:"created_at" json:"createdAt"`
}

func (u *Profile) PK() string {
	return "id"
}

func (u *Profile) TableName() string {
	return "profile"
}

type Photo struct {
	ID        int    `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"userID"`
	Url       string `db:"url" json:"url" validate:"required,url"`
	CreatedAt string `db:"created_at" json:"createdAt"`
}

func (u *Photo) PK() string {
	return "id"
}

func (u *Photo) TableName() string {
	return "photo"
}

type RichUser struct {
	User
	Profile *Profile `json:"profile" db:"-" relation:"id,user_id" connection:"rrest"`
	Photos  []Photo  `json:"photos" db:"-" relation:"id,user_id" connection:"rrest" ` //one-to-many
}

func main() {

	fmt.Printf("user: %T\n", &User{})

	rlib.DefaultRooboConfig()

	user := &User{}
	rsql.Use("rrest").Model(user).Where("id=?", 2).Get()

	fmt.Printf("user: %v\n", user)

	r := mux.NewRouter()
	res := rrest.NewResource("user", &RichUser{}, "id", "rrest")
	res.Route(r)

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		fmt.Printf("%s %s\n", strings.Join(methods, ","), pathTemplate)
		return nil
	})
	fmt.Println(err)

	serveAt := ":9004"
	rlib.Info(context.Background(), "all routers are setup, start serve at >", serveAt)
	err = rlib.HttpServe(r, serveAt)
	fmt.Printf("%v", err)
}
