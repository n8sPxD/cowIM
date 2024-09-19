package mongo

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/mon"
)

type DB struct {
	*mon.Model
}

func MustNewMongo(url string, opts mon.Option) *mon.Model {
	db, err := mon.NewModel(url, "im_server_db", "messages", opts)
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}
	fmt.Println("Connected to MongoDB")
	return db
}
