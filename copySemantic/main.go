// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go call parameter passing semantic

package main

import "fmt"

type User struct {
	Name string
}

func ModifyPtr(u *User) {
	u = &User{Name: "Ptr"}
}

func Modify(u User) {
	u.Name = "Modify"
}

func ModifyPtrPtr(u **User) {
	*u = &User{Name: "PtrPtr"}
}

func main() {
	u := &User{Name: "Paul"}
	println(u.Name)
	ModifyPtr(u)
	println(u.Name)

	u2 := User{Name: "Paul"}
	println(u2.Name)
	ModifyPtr(&u2)
	println(u2.Name)

	u3 := User{Name: "Paul"}
	println(u3.Name)
	Modify(u3)
	println(u3.Name)

	u4 := &User{Name: "Paul"}
	fmt.Println(u4.Name)
	ModifyPtrPtr(&u4)
	fmt.Println(u4.Name)

}
