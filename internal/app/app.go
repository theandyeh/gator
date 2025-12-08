package app

import (
	"github.com/theandyeh/gator/internal/config"
	"github.com/theandyeh/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}
