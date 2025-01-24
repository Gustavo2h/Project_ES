package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func executeRequest(req *http.Request, router http.Handler) *http.Response {
	// Função auxiliar para executar a requisição e retornar a resposta
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr.Result()
}

func TestLogin(t *testing.T) {
	// Testa o login e a geração do token JWT

	loginRequest := loginRequest{
		Email:    "loja_a@example.com",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(loginRequest)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	resp := executeRequest(req, router)

	assert.Equal(t, 200, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(t, response["token"])
}

func TestCreateProductWithAuth(t *testing.T) {
	// Testa a criação de um produto com autenticação JWT

	createProductRequest := product{
		Name:     "Teclado",
		Type:     "Eletrônico",
		Quantity: 50,
	}
	jsonValue, _ := json.Marshal(createProductRequest)

	// Obter o token com o login
	loginRequest := loginRequest{
		Email:    "loja_a@example.com",
		Password: "password123",
	}
	jsonValueLogin, _ := json.Marshal(loginRequest)
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValueLogin))
	respLogin := executeRequest(reqLogin, router)

	var response map[string]interface{}
	json.NewDecoder(respLogin.Body).Decode(&response)
	token := response["token"].(string)

	// Fazer a requisição para criar o produto com o token JWT
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token) // Adicionando o token no cabeçalho
	resp := executeRequest(req, router)

	assert.Equal(t, 201, resp.StatusCode)
}

func TestDeleteProductWithAuth(t *testing.T) {
	// Testa a exclusão de um produto com autenticação JWT

	productID := "P1"

	// Obter o token com o login
	loginRequest := loginRequest{
		Email:    "loja_a@example.com",
		Password: "password123",
	}
	jsonValueLogin, _ := json.Marshal(loginRequest)
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValueLogin))
	respLogin := executeRequest(reqLogin, router)

	var response map[string]interface{}
	json.NewDecoder(respLogin.Body).Decode(&response)
	token := response["token"].(string)

	// Fazer a requisição para excluir o produto com o token JWT
	req, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := executeRequest(req, router)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestListProducts(t *testing.T) {
	// Testa a listagem de produtos (rota pública)
	req, _ := http.NewRequest("GET", "/products", nil)
	resp := executeRequest(req, router)

	assert.Equal(t, 200, resp.StatusCode)

	var response []product
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Greater(t, len(response), 0) // Verifica se a lista de produtos não está vazia
}

func TestAddCustomer(t *testing.T) {
	// Testa o cadastro de cliente (rota pública)
	newCustomer := customer{
		Name:     "Lucas Costa",
		Email:    "lucas@example.com",
		Password: "password123",
		CPF:      "12345678900",
		Phone:    "11999990000",
		Address:  "Rua XYZ, 123",
	}
	jsonValue, _ := json.Marshal(newCustomer)

	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonValue))
	resp := executeRequest(req, router)

	assert.Equal(t, 201, resp.StatusCode)

	var response customer
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, newCustomer.Name, response.Name)
}

func TestAddSeller(t *testing.T) {
	// Testa o cadastro de vendedor (rota pública)
	newSeller := seller{
		Name:     "Vendedor XYZ",
		Email:    "vendedor@example.com",
		Password: "password123",
		CNPJ:     "12345678000123",
		Phone:    "11988880000",
		Address:  "Avenida ABC, 456",
	}
	jsonValue, _ := json.Marshal(newSeller)

	req, _ := http.NewRequest("POST", "/sellers", bytes.NewBuffer(jsonValue))
	resp := executeRequest(req, router)

	assert.Equal(t, 201, resp.StatusCode)

	var response seller
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, newSeller.Name, response.Name)
}
