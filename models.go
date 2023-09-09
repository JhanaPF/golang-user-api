package main

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Mail     string `json:"mail"`
	Password string `json:"-"`
}

type Relationship struct {
	ID       uint `json:"id"`
	UserID   uint `json:"user_id"`
	FriendID uint `json:"friend_id"`
}

type GameHistory struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}

type CourseChoice struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}
