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

type Transaction struct{
	//name fileds
	id,sum,fromId,toID,amount,Date,operation 	

}

type PG struct {
	log *logrus.Entry
	db  *sqlx.DB
	dsn string
}
func (pg *PG) Close{
	if err :=pg.db.Close();err != nil{
		pg.log.Error("err closing pg connection: %w",err)
	}
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

func (pg *PG) Transfer(senderId int, recieverId int, amount float64) error {
	tx, err := pg.db.BeginTx(nil, nil)
	defer func() {
		if err = tx.Rollback(); err != nil {
			pg.log.Error("err rolling back transaction")
		}
	}()
	query := `UPDATE wallet SET balance = balance - $1 WHERE id  = $2 RETURNING balance`

	_, err = tx.Exec(query, amount, senderId)

	if err != nil {
		return err
	}

	query = `UPDATE wallet set balance = balance + $1 where id = $2 returning balance`

	`INSERT into transactions`

	_, err = tx.Exec(query, amount, recieverId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {

		return fmt.Errorf("err comminting the transaction")
	}
	return nil

}

//Withdrawal
