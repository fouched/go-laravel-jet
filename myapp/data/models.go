package data

import (
	"database/sql"
	"fmt"
	upperdb "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
	"os"
)

var db *sql.DB
var upper upperdb.Session

// Models any models created here (and in the New function)
// are easily accessible throughout the entire application
type Models struct {
	Users  User
	Tokens Token
}

func New(databasePool *sql.DB) Models {
	db = databasePool

	if os.Getenv("DATABASE_TYPE") == "mysql" || os.Getenv("DATABASE_TYPE") == "mariadb" {
		upper, _ = mysql.New(databasePool)
	} else {
		upper, _ = postgresql.New(databasePool)
	}

	return Models{
		Users:  User{},
		Tokens: Token{},
	}
}

func getInsertID(i upperdb.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}
	return i.(int)
}
