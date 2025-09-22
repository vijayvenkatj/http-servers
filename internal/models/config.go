package models

import (
	"sync/atomic"

	"github.com/vijayvenkatj/http-servsers/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	PlainTextReqs  atomic.Int32
	DbQueries	*database.Queries
	JWT_SECRET	string
	POLKA_KEY	string
}