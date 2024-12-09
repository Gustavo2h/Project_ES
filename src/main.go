package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret_key") // Chave para assinar os tokens

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Estruturas dos dados
type seller struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	CNPJ     string `json:"cnpj"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type customer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CPF      string `json:"cpf"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// "Banco de dados" em memória com dados iniciais
var sellers = []seller{
	{"S1", "Loja A", "12345678000101", "11999990000", "loja_a@example.com", "password123", "Endereço A"},
	{"S2", "Loja B", "22345678000101", "11988880000", "loja_b@example.com", "password456", "Endereço B"},
	{"S3", "Loja C", "32345678000101", "11977770000", "loja_c@example.com", "password789", "Endereço C"},
}

var customers = []customer{
	{"C1", "João Silva", "joao@example.com", "12345", "11122233344", "11966660000", "Endereço do João"},
	{"C2", "Maria Oliveira", "maria@example.com", "54321", "22233344455", "11955550000", "Endereço da Maria"},
	{"C3", "Carlos Pereira", "carlos@example.com", "password", "33344455566", "11944440000", "Endereço do Carlos"},
}

var products = []product{
	{"P1", "Notebook", "Eletrônico", 10},
	{"P2", "Smartphone", "Eletrônico", 15},
	{"P3", "Cadeira", "Móveis", 20},
	{"P4", "Livro de Go", "Livro", 30},
	{"P5", "Mesa de Escritório", "Móveis", 5},
}

func main() {
	router := gin.Default()

	// Rotas públicas
	router.GET("/products", listProducts) // Qualquer um pode listar produtos
	router.POST("/login", login)          // Login para obter token

	// Rotas protegidas
	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		// Rotas exclusivas para vendedores
		protected.POST("/products", sellerMiddleware(), createProduct)
		protected.DELETE("/products/:id", sellerMiddleware(), deleteProduct)

		// Rotas gerais (somente autenticadas)
		protected.GET("/customers", listCustomers)
		protected.GET("/sellers", listSellers)
	}

	router.Run("localhost:8080")
}

// Middleware de autenticação
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["userID"])
			c.Set("userType", claims["type"])
		}
		c.Next()
	}
}

// Middleware para verificar se o usuário é vendedor
func sellerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, _ := c.Get("userType")
		if userType != "seller" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Only sellers can perform this action"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Login para gerar tokens
func login(c *gin.Context) {
	var request loginRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Verifica credenciais
	var userID, userType string
	for _, seller := range sellers {
		if seller.Email == request.Email && seller.Password == request.Password {
			userID = seller.ID
			userType = "seller"
			break
		}
	}
	for _, customer := range customers {
		if customer.Email == request.Email && customer.Password == request.Password {
			userID = customer.ID
			userType = "customer"
			break
		}
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	// Cria o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"type":   userType,
		"exp":    time.Now().Add(time.Hour * 1).Unix(), // Expiração de 1 hora pro token
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Operações de Produto
func createProduct(c *gin.Context) {
	var newProduct product
	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	newProduct.ID = generateID(len(products))
	products = append(products, newProduct)
	c.JSON(http.StatusCreated, newProduct)
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func listProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// Funções auxiliares (listagem e gerar IDs)
func listSellers(c *gin.Context) {
	c.JSON(http.StatusOK, sellers)
}

func listCustomers(c *gin.Context) {
	c.JSON(http.StatusOK, customers)
}

func generateID(index int) string {
	return string('A' + index)
}

