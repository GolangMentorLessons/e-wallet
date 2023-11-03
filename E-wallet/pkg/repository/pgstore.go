package repository

import (
	"E-wallet/pkg/metrics"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Wallet struct {
	ID            int       `json:"id"`
	Owner         string    `json:"owner"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AccountNumber string    `json:"account_number"`
}

// id,sum,fromId,toID,amount,Date,operation
type Transaction struct {
	ID         int       `json:"id"`
	FromID     int       `json:"from_id"`
	ToID       int       `json:"to_id"`
	FromWallet int       `json:"from_wallet"`
	ToWallet   int       `json:"to_wallet"`
	Amount     int       `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	Operation  string    `json:"operation"`
}

type Withdraw struct {
	Amount   float64
	WalletID int
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
	started := time.Now()
	defer func() {
		metrics.MetricDBRequestDuration.WithLabelValues("CreateWallet").Observe(time.Since(started).Seconds())
	}()
	query := `INSERT INTO wallet(owner,balance,updated_at,account_number) VALUES ($1,$2,$3,$4) returning id`
	var id int
	row := pg.db.QueryRow(query, wallet.Owner, wallet.Balance, wallet.UpdatedAt, wallet.AccountNumber)
	if err := row.Scan(&id); err != nil {
		metrics.MetricErrCount.WithLabelValues("CreateWallet").Inc()
		return 0, fmt.Errorf("err creating wallet: %w", err)
	}

	return id, nil
}

// Update(````),Delete(```“),Get(``````“)
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

func (pg *PG) createTransaction(transaction Transaction) (int, error) {
	query := `INSERT INTO transaction(from_wallet,to_wallet,amount,created_at,operation) VALUES ($1,$2,$3,$4,$5) returning id`
	var id int
	row := pg.db.QueryRow(query, transaction.FromWallet, transaction.ToWallet, transaction.Amount, time.Now(), transaction.Operation)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("err creating transaction: %w", err)
	}

	return id, nil
}

func (pg *PG) Transfer(transaction Transaction) (int, error) {
	//TODO: change nils, do not pass a nil Context, even if a function permits it; pass context.TODO if you are unsure about which Context to use (SA1012)
	tx, errRoll := pg.db.BeginTx(nil, nil)
	defer func() {
		if errRoll = tx.Rollback(); errRoll != nil {
			pg.log.Error("err rolling back transaction")
		}
	}()

	query := `UPDATE wallet SET balance = balance - $1 WHERE id = $2 RETURNING balance`
	_, err := pg.db.Exec(query, transaction.Amount, transaction.FromID)

	if err != nil {
		return 0, fmt.Errorf("error with update from id balance:%w", err)
	}

	query = `UPDATE wallet SET balance = balance + $1 WHERE id = $2 RETURNING balance`
	_, err = pg.db.Exec(query, transaction.Amount, transaction.ToID)
	if err != nil {
		return 0, fmt.Errorf("error with update to id balance:%w", err)
	}

	newTxId, err := pg.createTransaction(transaction)
	if err != nil {
		return 0, fmt.Errorf("error with create transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {

		return 0, fmt.Errorf("err comminting the transaction")
	}

	return newTxId, nil
}

func (pg *PG) Withdraw(transaction Transaction) (int, error) {
	tx, errRoll := pg.db.BeginTx(nil, nil)
	defer func() {
		if errRoll = tx.Rollback(); errRoll != nil {
			pg.log.Error("err rolling back transaction")
		}
	}()
	query := `UPDATE wallet set balance = balance - $1 WHERE id = $2`

	_, err := pg.db.Exec(query, transaction.Amount, transaction.FromID)

	if err != nil {
		return 0, fmt.Errorf("error in withdraw:%w", err)
	}
	//TODO: при вводе toId и toWallet что с ними или передается null?
	newIdTx, err := pg.createTransaction(transaction)
	if err != nil {
		return 0, fmt.Errorf("error with create transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {

		return 0, fmt.Errorf("err comminting the transaction")

	}
	return newIdTx, nil
}

func (pg *PG) Close() {
	if err := pg.db.Close(); err != nil {
		pg.log.Errorf("err closing pg connection :%v", err)
	}
}
