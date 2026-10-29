package chatrep

type db interface {
	CreateChat(string, ...any)
	DeleteChatByID(id int)
	ChangeChatSettings(id int)
	CreateMember()
}
type ChatRepository struct {
	Db db
}
