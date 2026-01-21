package models

type Message struct {
	ChatID int64  `json:"chat_id"`
	UserID int64  `json:"from_id"`
	Text   string `json:"text"`
}


type User struct{
	ChatID int64 `json:"chat_id"`
}

// // состояния для валидации и сохранения
// type UserState struct {
// 	CurrentStep string
// 	Valentine *Message
// }

// func NewUserState(userID int64) *UserState {
// 	return &UserState{
// 		CurrentStep: "idle",
// 		Valentine:   &Message{
// 			FromID: userID,
// 		},
// 	}
// }
