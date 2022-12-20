package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/syucel96/simplebank/db/sqlc"
	"github.com/syucel96/simplebank/token"
	"github.com/syucel96/simplebank/util"
)

// Server that serves http requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	var tokenMaker token.Maker
	var err error
	if config.TokenMaker {
		tokenMaker, err = token.NewJWTMaker(config.TokenSymmetricKey)
	} else {
		tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)
	}

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("amount", validAmount)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	// router.GET("/users/:username", server.getUser)
	router.POST("/users/login", server.loginUser)
	// router.GET("/users", server.listUsers)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfer", server.CreateTransfer)
	/* 	router.PUT("/accounts", server.updateAccount)
	   	router.DELETE("/accounts/:id", server.deleteAccount) */

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
