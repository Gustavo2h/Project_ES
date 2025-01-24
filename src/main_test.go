package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
    testCases := []struct {
        index    int
        prefix   string
        expected string
    }{
        {0, "A", "A0"},
        {1, "B", "B1"},
        {25, "Z", "Z25"},
    }

    for _, tc := range testCases {
        t.Run(tc.expected, func(t *testing.T) {
            result := generateID(tc.index, tc.prefix)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func TestListProducts(t *testing.T) {
    router := gin.Default()
    router.GET("/products", listProducts)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/products", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response []product
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response)
}

func TestLogin(t *testing.T) {
    router := gin.Default()
    router.POST("/login", login)

    loginPayload := `{"email":"joao@example.com","password":"12345"}`
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/login", strings.NewReader(loginPayload))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response["token"])
}

func TestCreateProduct(t *testing.T) {
    router := gin.Default()
    router.POST("/products", createProduct)

    newProduct := `{"name":"New Product","type":"New Type","quantity":10}`
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/products", strings.NewReader(newProduct))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var response product
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "New Product", response.Name)
}

func TestDeleteProduct(t *testing.T) {
    router := gin.Default()
    router.DELETE("/products/:id", deleteProduct)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/products/P1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Product deleted", response["message"])
}

func TestListCustomers(t *testing.T) {
    router := gin.Default()
    router.GET("/customers", listCustomers)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/customers", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response []customer
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response)
}

func TestListSellers(t *testing.T) {
    router := gin.Default()
    router.GET("/sellers", listSellers)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/sellers", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response []seller
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response)
}

func TestAddCustomer(t *testing.T) {
    router := gin.Default()
    router.POST("/customers", addCustomer)

    newCustomer := `{"name":"New Customer","email":"newcustomer@example.com","password":"password","cpf":"12345678900","phone":"11999999999","address":"New Address"}`
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/customers", strings.NewReader(newCustomer))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var response customer
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "New Customer", response.Name)
}

func TestAddSeller(t *testing.T) {
    router := gin.Default()
    router.POST("/sellers", addSeller)

    newSeller := `{"name":"New Seller","cnpj":"98765432100001","phone":"11988888888","email":"newseller@example.com","password":"password","address":"New Address"}`
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/sellers", strings.NewReader(newSeller))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var response seller
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "New Seller", response.Name)
}
