package connrepos

import (
	"chat/api/internal/mydb"
)

type Repository struct {
	Db mydb.Db
}
