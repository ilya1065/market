package inits

import (
	"Product_Service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var DB *sqlx.DB

func ConnectToDB(conf *config.Config) *sqlx.DB {
	db, err := sqlx.Connect(conf.DBDriver, conf.DBDSN)
	if err != nil {
		log.Fatal("не удалось подключится к базе данных ", err)
	}
	return db
}
