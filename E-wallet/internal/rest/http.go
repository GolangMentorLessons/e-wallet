package rest

import (
	repo "E-wallet/pkg/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreateWallet(wallet repo.Wallet) (int, error)
	UpdateWallet(id int, wallet repo.Wallet) (repo.Wallet, error)
	GetWallet(id int, currency string) (repo.Wallet, error)
	GetAllWallets(owner string) ([]repo.Wallet, error)
	DeleteWallet(id int) (int, error)

	//transaction
	Transfer(transaction repo.Transaction) (int, error)
	Withdraw(transaction repo.Transaction) (int, error)
	CheckBalance(id int) (wallet repo.Wallet, err error)
}
type Router struct {
	log     *logrus.Entry
	router  *gin.Engine
	service Service
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func NewRouter(log *logrus.Logger, service Service) *Router {
	r := &Router{
		log:     log.WithField("transport", "e-wallet"),
		router:  gin.Default(),
		service: service,
	}

	r.router.GET("/metrics", prometheusHandler())
	g := r.router.Group("service/v1")
	g.GET("/wallet/:id:", r.getWallets)
	g.GET("/wallet/:owner", r.getWallet)
	g.POST("/wallet", r.createWallet)
	g.PUT("/wallet/:id:", r.updateWallet)
	g.DELETE("/wallet/:id", r.deleteWallet)

	g.PUT("/wallet/transfer", r.transfer)
	g.PUT("/wallet/withdraw", r.withdraw)
	g.GET("/wallet/checkbalance/:id", r.checkBalance)

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

func (r *Router) updateWallet(c *gin.Context) {
	var input repo.Wallet
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err = c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	input, err = r.service.UpdateWallet(id, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": input})
}

func (r *Router) getWallet(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	wallet, err := r.service.GetWallet(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

func (r *Router) getWallets(c *gin.Context) {
	owner := c.Param("owner")
	if len(owner) <= 0 {
		c.JSON(http.StatusBadRequest, fmt.Errorf("param is empty"))
		return
	}

	wallets, err := r.service.GetAllWallets(owner)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

func (r *Router) deleteWallet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, err = r.service.DeleteWallet(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id delete": id})
}

func (r *Router) transfer(c *gin.Context) {
	var input repo.Transaction
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	idTX, err := r.service.Transfer(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"transfer": idTX})
}

func (r *Router) withdraw(c *gin.Context) {
	var input repo.Transaction
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	idTX, err := r.service.Withdraw(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"transfer": idTX})
}

func (r *Router) checkBalance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	wallet, err := r.service.CheckBalance(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"walletBalance": wallet})

}
