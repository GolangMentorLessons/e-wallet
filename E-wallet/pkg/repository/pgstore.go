package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
func (pg *PG) UpdateWallet(id int , newBalance float64) error {

	query := `UPDATE wallet SET balance=?, UpdatedAt=? WHERE id=?`
	_, err := pg.db.Exec(query, newBalance, time.Now(), id)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Println("Update is success")

	return nil

}

func (pg *PG) GetAllWallets()(wallets []Wallet, err error){
	
	query := `SELECT * FROM wallet`

	err = pg.db.Select(&wallets, query)

	if err != nil{
		return wallets, fmt.Errorf("gets wallets failed: %w", err )
	}

	return
}

func (pg *PG) GetWallet(id int) (w Wallet, err error) {
	query := `SELECT * FROM wallet WHERE id=?`
	
	err = pg.db.Get(&w,query,id)
	if err != nil{
		return w, fmt.Errorf("wallet not found %w", err)
	}

	return
}

func (pg *PG) DeleteWallet(id int) (int, error){

	  query := `Delete FROM wallet WHERE id=?`
	_, err := pg.db.Exec(query,id)

	if err != nil {
		return 0, fmt.Errorf("delete failed: %w", err)
	}

	fmt.Println("Delete is success")

	return id, nil
}


