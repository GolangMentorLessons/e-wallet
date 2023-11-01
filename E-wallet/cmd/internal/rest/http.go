package rest

import (
	repo "E-wallet/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Service interface {
	CreateWallet(wallet repo.Wallet) (int, error)
}
type Router struct {
	log     *logrus.Entry
	router  *gin.Engine
	service Service
}

func NewRouter(log *logrus.Logger, service Service) *Router {
	r := &Router{
		log:     log.WithField("transport", "e-wallet"),
		router:  gin.Default(),
		service: service,
	}
	r.router.POST("/wallet", r.createWallet)

	return r

}

func (r *Router) createWallet(c *gin.Context) {
	var input repo.Wallet
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	id, err := r.service.CreateWallet(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
