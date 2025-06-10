package model

import "time"

type StoreUserResponse struct {
	Data    StoreData `json:"data"`
	Message string    `json:"message"`
	Status  bool      `json:"status"`
}

type StoreData struct {
	Store
	User     UserData     `json:"user"`
	Category CategoryData `json:"category"`
}

type UserData struct {
	Address     string     `json:"address"`
	CreateAt    time.Time  `json:"createAt"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	ID          int        `json:"id"`
	IsActive    bool       `json:"isActive"`
	IsBanned    bool       `json:"is_banned"`
	LastName    string     `json:"lastName"`
	PhoneNumber string     `json:"phoneNumber"`
	Stores      []StoreRef `json:"stores"`
	UpdateAt    time.Time  `json:"updateAt"`
}
type CategoryData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type StoreRef struct {
	ID int `json:"id"`
}

func (s *StoreUserResponse) GetStore() Store {
	store := s.Data.Store
	store.CategoryName = s.Data.Category.Name
	return store
}
