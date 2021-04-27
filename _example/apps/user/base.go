package user

import "fmt"

type Gender string

const GenderBoy Gender = "boy"
const GenderGirl Gender = "girl"
const GenderOther Gender = "unknown"

type User struct {
	Id int64 `json:"id"`
	Email string `json:"-"`
	Passwd string `json:"-"`
	NickName string `json:"nick_name"`
	Birthday string `json:"birthday"`
	Gender Gender `json:"gender"`
}

func (user *User) Token() string {
	return fmt.Sprintf("token_%d", user.Id)
}


var Users = []*User{
	{ 100, "admin@nutshell.com", "123456", "admin", "1990/01/01", GenderGirl},
	{ 101,"joy@nutshell.com", "123456", "joy", "1991/01/02", GenderBoy},
	{ 102,"john@nutshell.com", "123456", "john", "1992/01/03", GenderBoy},
	{ 103,"jack@nutshell.com", "123456", "jack", "1993/01/04", GenderOther},
	{ 104,"kim@nutshell.com", "123456", "kim", "1980/11/05", GenderGirl},
	{ 105,"ken@nutshell.com", "123456", "ken", "1981/08/21", GenderGirl},
	{ 106,"kaven@nutshell.com", "123456", "kaven", "1982/12/18", GenderGirl},
}

var userIdx = make(map[string]*User, 0)
var userIdIdx = make(map[int64]*User, 0)

func init() {
	for _, user := range Users {
		userIdx[user.Email] = user
		userIdIdx[user.Id] = user
	}
}

