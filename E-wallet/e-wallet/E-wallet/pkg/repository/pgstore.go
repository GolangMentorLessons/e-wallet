package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type Wallet struct {
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PG struct {
	log *logrus.Entry
	db  *sqlx.DB
	dsn string
}

func NewRepo(dsn string, log *logrus.Logger) (*PG, error) {
	db, err := sqlx.Connect("sqlx", dsn)
	if err != nil {
		return nil, fmt.Errorf("err connecting to PG: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("err pinging database after connection: %w", err)
	}
	pg := &PG{
		log: log.WithField("database", "NewRepo"),
		db:  db,
		dsn: dsn,
	}

	return pg, nil
}

func (pg *PG) CreateWallet(wallet Wallet) (int, error) {
	query := `INSERT INTO wallet(owner,balance,updated_at) VALUES ($1,$2,$3) returning id`
	var id int
	row := pg.db.QueryRow(query, wallet.Owner, wallet.Balance, wallet.UpdatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("err creating wallet: %w", err)
	}
	return id, nil
}

//Update(````),Delete(`````),Get(````````)
