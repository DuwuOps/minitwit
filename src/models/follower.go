package models

type Follower struct {
	WhoID  int `db:"who_id"`  
	WhomID int `db:"whom_id"` 
}
