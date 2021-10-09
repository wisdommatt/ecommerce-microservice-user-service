package users

import "time"

type User struct {
	ID          string    `json:"id" bson:"id,omitempty"`
	FullName    string    `json:"fullName" bson:"fullName,omitempty"`
	Email       string    `json:"email" bson:"email,omitempty"`
	Password    string    `json:"password" bson:"password,omitempty"`
	Country     string    `json:"country" bson:"country,omitempty"`
	TimeAdded   time.Time `json:"timeAdded" bson:"timeAdded,omitempty"`
	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated,omitempty"`
}
