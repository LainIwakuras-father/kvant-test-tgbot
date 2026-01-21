package interfaces

import (
	"context"

)

type IStorage interface {
	Save(ctx context.Context, userID int64, chatID int64, text string) error
	Count() int
	GetAllUsers(ctx context.Context) (map[int64]string,error)
	UserExists(ctx context.Context, chatID int64) (bool,error)
	AddUser(chatID int64,username string) error
}
