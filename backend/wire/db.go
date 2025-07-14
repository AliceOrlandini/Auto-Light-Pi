package wire

import (
	"database/sql"

	"github.com/AliceOrlandini/Auto-Light-Pi/config"
)

func ProvideDB() (*sql.DB, error) {
	err := config.InitDB()
	if err != nil {
		return nil, err
	}
	return config.DB, nil
}