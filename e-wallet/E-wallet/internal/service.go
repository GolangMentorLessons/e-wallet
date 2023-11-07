package internal

import (
	repo "E-wallet/pkg/repository"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Storage interface {
	CreateWallet(wallet repo.Wallet) (int, error)
	UpdateWallet(id int, wallet repo.Wallet) (repo.Wallet, error)
	GetWallet(id int) (wallet repo.Wallet, err error)
	GetAllWallets(owner string) (wallets []repo.Wallet, err error)
	DeleteWallet(id int) (int, error)

	//transaction
	Transfer(transaction repo.Transaction) (int, error)
	Withdraw(transaction repo.Transaction) (int, error)
	CheckBalance(id int) (wallet repo.Wallet, err error)
}

type Exchange interface {
	Conversion(currency string, amount float64) (float64, error)
}

type Service struct {
	log      *logrus.Entry
	store    Storage
	exchange Exchange
}

func NewService(log *logrus.Logger, store Storage) *Service {
	return &Service{
		log: log.WithField("service", "e-wallet"),
	}
}

func (s *Service) CreateWallet(wallet repo.Wallet) (int, error) {
	id, err := s.store.CreateWallet(wallet)
	if err != nil {
		return 0, fmt.Errorf("err creating wallet: %v", err)
	}
	return id, nil
}

func (s *Service) UpdateWallet(id int, wallet repo.Wallet) (repo.Wallet, error) {
	wallet, err := s.store.UpdateWallet(id, wallet)
	if err != nil {
		return repo.Wallet{}, fmt.Errorf("err updating wallet: %v", err)
	}

	return wallet, nil
}

func (s *Service) GetWallet(id int, currency string) (repo.Wallet, error) {
	wallet, err := s.store.GetWallet(id)
	if err != nil {
		return repo.Wallet{}, fmt.Errorf("err getWalet with id: %v", err)
	}
	if currency != "" {
		wallet.Balance, err = s.exchange.Conversion(currency, wallet.Balance)
		if err != nil {
			return repo.Wallet{}, err
		}
	}

	return wallet, nil
}

func (s *Service) GetAllWallets(owner string) ([]repo.Wallet, error) {
	wallet, err := s.store.GetAllWallets(owner)
	if err != nil {
		return []repo.Wallet{}, fmt.Errorf("error gets all wallets owner:%v", err)
	}
	return wallet, nil
}

func (s *Service) DeleteWallet(id int) (int, error) {
	id, err := s.store.DeleteWallet(id)
	if err != nil {
		return 0, fmt.Errorf("error gets all wallets owner:%v", err)
	}
	return id, nil
}

func (s *Service) Transfer(transaction repo.Transaction) (int, error) {

	txId, err := s.store.Transfer(transaction)

	if err != nil {
		return 0, fmt.Errorf("error transfer amount: %v", err)
	}

	return txId, nil
}

func (s *Service) Withdraw(transaction repo.Transaction) (int, error) {
	txId, err := s.store.Withdraw(transaction)
	if err != nil {
		return 0, fmt.Errorf("error transfer amount: %v", err)
	}

	return txId, nil
}

func (s *Service) CheckBalance(id int) (wallet repo.Wallet, err error) {
	wallet, err = s.store.CheckBalance(id)

	if err != nil {
		return wallet, fmt.Errorf("error transfer amount: %v", err)
	}

	return

}
