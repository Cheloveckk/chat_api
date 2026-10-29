package userrep

type db interface {
	CreateMessage(string, ...any)
	DeleteMessageByID(id int)
	ChangeMessage(id int)
}
type MessageRepository struct {
	Db db
}
