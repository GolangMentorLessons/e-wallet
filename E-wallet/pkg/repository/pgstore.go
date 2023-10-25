package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type Wallet struct {
	Owner         string    `json:"owner"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AccountNumber string    `json:"account_number"`
}

type PG struct {
	log *logrus.Entry
	db  *sqlx.DB
	dsn string
}

var (
	ErrWalletNotFound = fmt.Errorf("err wallet not found")
)

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
	query := `INSERT INTO wallet(owner,balance,updated_at,account_number) VALUES ($1,$2,$3,$4) returning id`
	var id int
	row := pg.db.QueryRow(query, wallet.Owner, wallet.Balance, wallet.UpdatedAt, wallet.AccountNumber)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("err creating wallet: %w", err)
	}
	return id, nil
}

//Update(````),Delete(`````),Get(````````)

// UpdateWallet(id int,balance float)(......){
//}

func (pg *PG) UpdateWallet(id int, wallet Wallet) (Wallet, error) {

	query := `UPDATE wallet SET balance = ?, updatedAt = ? WHERE id = ?`
	row := pg.db.QueryRow(query, wallet.Balance, time.Now(), id)
	err := row.Scan(&wallet)
	if err != nil {
		return Wallet{}, ErrWalletNotFound
	}
	return wallet, nil
}

func (pg *PG) GetWallet(id int) (wallet Wallet, err error) {
	query := `SELECT * FROM wallet WHERE id=?`

	err = pg.db.Get(&wallet, query, id)
	if err != nil {
		return wallet, fmt.Errorf("wallet not found %w", err)
	}

	return
}

func (pg *PG) GetAllWallets(owner string) (wallets []Wallet, err error) {

	query := `SELECT * FROM wallet where owner = ?`

	err = pg.db.Select(&wallets, query, owner)

	if err != nil {
		return wallets, fmt.Errorf("gets wallets failed: %w", err)
	}

	return
}

func (pg *PG) DeleteWallet(id int) (int, error) {

	query := `Delete FROM wallet WHERE id=?`
	_, err := pg.db.Exec(query, id)

	if err != nil {
		return 0, fmt.Errorf("delete failed: %w", err)
	}

	fmt.Println("Delete is success")

	return id, nil
}
