package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"username"`
	Password   string             `bson:"password"`
	School     string             `bson:"school"`
	Email      string             `bson:"email"`
	Role       string             `bson:"roles"`
	Source     string             `bson:"source"`
	IsDelete   bool               `bson:"is_delete"`
	CreateTime time.Time          `bson:"create_time"`
	LatestTime time.Time          `bson:"latest_time"`
}

func UserInit() *User {
	return &User{
		Role:     "user",
		Source:   "Competition_Platform",
		IsDelete: false,
	}
}

func (user *User) GetStringID() string {
	return user.ID.Hex()
}

func (user *User) GetInfo() map[string]any {
	return map[string]any{
		"id":       user.GetStringID(),
		"username": user.Name,
		//"password":    user.Password,
		"school":      user.School,
		"email":       user.Email,
		"roles":       user.Role,
		"source":      user.Source,
		"is_delete":   user.IsDelete,
		"create_time": user.CreateTime.Format("2006-01-02 15:04:05"),
		"latest_time": user.LatestTime.Format("2006-01-02 15:04:05"),
	}
}

func (user *User) GetShortInfo() map[string]any {
	return map[string]any{
		"id":       user.GetStringID(),
		"username": user.Name,
		"roles":    user.Role,
		"source":   user.Source,
	}
}

func (user *User) SetValue(name string, password string, school string, email string, role string) {
	if name != "" {
		user.Name = name
	}
	if password != "" {
		user.Password = password
	}
	if school != "" {
		user.School = school
	}
	if email != "" {
		user.Email = email
	}
	if role != "" {
		if role != "admin" {
			role = "user"
		}
		user.Role = role
	}
}
