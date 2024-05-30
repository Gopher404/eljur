package models

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Group    int8   `json:"group"`
}

const (
	PermBlock int32 = iota
	PermStudent
	PermAdmin
)
