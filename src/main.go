package main

import (
    "net/http"
    "time"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret_key") // Chave para assinar os tokens

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Estruturas dos dados
type Seller struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    CNPJ     string `json:"cnpj"`
    Phone    string `json:"phone"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Address  string `json:"address"`
}

type Customer struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    CPF      string `json:"cpf"`
    Phone    string `json:"phone"`
    Address  string `json:"address"`
}

type Product struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Type     string `json:"type"`
    Quantity int    `json:"quantity"`
}

// "Banco de dados" em memória com dados iniciais
var Sellers = []Seller{
    {"S1", "Loja A", "12345678000101", "11999990000", "loja_a@example.com", "password123", "Endereço A"},
    {"S2", "Loja B", "22345678000101", "11988880000", "loja_b@example.com", "password456", "Endereço B"},
    {"S3", "Loja C", "32345678000101", "11977770000", "loja_c@example.com", "password789", "Endereço C"},
}

var Customers = []Customer{
    {"C1", "João Silva", "joao@example.com", "12345", "11122233344", "11966660000", "Endereço do João"},
    {"C2", "Maria Oliveira", "maria@example.com", "54321", "22233344455", "11955550000", "Endereço da Maria"},
    {"C3", "Carlos Pereira", "carlos@example.com", "password", "33344455566", "11944440000", "Endereço do Carlos"},
}

var Products = []Product{
    {"P1", "Notebook", "Eletrônico", 10},
    {"P2", "Smartphone", "Eletrônico", 15},
    {"P3", "Cadeira", "Móveis", 20},
    {"P4", "Livro de Go", "Livro", 30},
    {"P5", "Mesa de Escritório", "Móveis", 5},
}

func main() {
    router := gin.Default()

    // Rotas públicas
    router.GET("/products", ListProducts) // Qualquer um pode listar produtos
    router.POST("/login", Login)          // Login para obter token
    router.POST("/customers", AddCustomer) // Cadastro de cliente sem autenticação
    router.POST("/sellers", AddSeller)    // Cadastro de vendedor sem autenticação

    // Rotas protegidas
    protected := router.Group("/")
    protected.Use(AuthMiddleware())
    {
        // Rotas exclusivas para vendedores
        protected.POST("/products", SellerMiddleware(), CreateProduct)
        protected.DELETE("/products/:id", SellerMiddleware(), DeleteProduct)

        // Rotas gerais (somente autenticadas)
        protected.GET("/customers", ListCustomers)
        protected.GET("/sellers", ListSellers)
    }

    router.Run("localhost:8080")
}

// Middleware de autenticação
func AuthMiddleware() gin.HandlerFunc {
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
func SellerMiddleware() gin.HandlerFunc {
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
func Login(c *gin.Context) {
    var request LoginRequest
    if err := c.BindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    // Verifica credenciais
    var userID, userType string
    for _, seller := range Sellers {
        if seller.Email == request.Email && seller.Password == request.Password {
            userID = seller.ID
            userType = "seller"
            break
        }
    }
    for _, customer := range Customers {
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
func CreateProduct(c *gin.Context) {
    var newProduct Product
    if err := c.BindJSON(&newProduct); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }
    newProduct.ID = GenerateID(len(Products), "P")
    Products = append(Products, newProduct)
    c.JSON(http.StatusCreated, newProduct)
}

func DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    for i, p := range Products {
        if p.ID == id {
            Products = append(Products[:i], Products[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func ListProducts(c *gin.Context) {
    c.JSON(http.StatusOK, Products)
}

// Funções auxiliares (listagem e gerar IDs)
func ListSellers(c *gin.Context) {
    c.JSON(http.StatusOK, Sellers)
}

func ListCustomers(c *gin.Context) {
    c.JSON(http.StatusOK, Customers)
}

func GenerateID(index int, prefix string) string {
    return prefix + strconv.Itoa(index)
}

// Função para adicionar cliente
func AddCustomer(c *gin.Context) {
    var newCustomer Customer
    if err := c.BindJSON(&newCustomer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }
    newCustomer.ID = GenerateID(len(Customers), "C")
    Customers = append(Customers, newCustomer)
    c.JSON(http.StatusCreated, newCustomer)
}

// Função para adicionar vendedor
func AddSeller(c *gin.Context) {
    var newSeller Seller
    if err := c.BindJSON(&newSeller); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }
    newSeller.ID = GenerateID(len(Sellers), "S")
    Sellers = append(Sellers, newSeller)
    c.JSON(http.StatusCreated, newSeller)
}
