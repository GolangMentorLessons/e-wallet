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
}

type Service struct {
	log   *logrus.Entry
	store Storage
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

func(s *Service) UpdateWallet(id int, wallet repo.Wallet) (repo.Wallet, error){
	wallet, err := s.store.UpdateWallet(id, wallet)
	if err != nil{
		return repo.Wallet{}, fmt.Errorf("err updating wallet: %v",err)
	}

	return wallet, nil
}

func(s *Service) GetWallet(id int) (repo.Wallet,error){
	wallet, err := s.store.GetWallet(id)
	if err != nil{
		return repo.Wallet{}, fmt.Errorf("err getWalet with id: %v",err)
	}
	return wallet, nil
}

func(s *Service) GetAllWallets(owner string) ([]repo.Wallet,error){
	wallet,err := s.store.GetAllWallets(owner)
	if err != nil {
		return []repo.Wallet{},fmt.Errorf("error gets all wallets owner:%v",err)
	}
	return wallet,nil
}

func(s *Service) DeleteWallet(id int) (int,error){
	id, err := s.store.DeleteWallet(id)
	if err != nil {
		return 0, fmt.Errorf("error gets all wallets owner:%v",err)
	}
	return id,nil
}