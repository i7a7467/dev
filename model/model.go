package model

import "github.com/uptrace/bun"

type Person struct {
	bun.BaseModel `bun:"table:person,alias:p"`

	ID int `json:"id"`
	Account string `json:"account"`
	DisplayName string `json:"displayName"`
	CreatedAt string `json:"createdAt"`
}

type Accounts struct {
	// bun.BaseModel `bun:"table:person,alias:as"`

	Persons []Person
}