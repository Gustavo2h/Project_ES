package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Função generateID
func generateID(index int, prefix string) string {
	return fmt.Sprintf("%s%d", prefix, index)
}

// Definição do tipo product e sua validação
type product struct {
	Name     string
	Type     string
	Quantity int
}

func validateProduct(p product) bool {
	return p.Name != "" && p.Quantity > 0
}

// Definição do tipo customer e sua validação
type customer struct {
	Name  string
	Email string
	CPF   string
}

func validateCustomer(c customer) bool {
	// Validar CPF: apenas verifica se o CPF tem 11 caracteres numéricos
	return c.Name != "" && isValidEmail(c.Email) && len(c.CPF) == 11
}

// Função para verificar se o email é válido
func isValidEmail(email string) bool {
	// Usando expressão regular para validar um email simples
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Definição do tipo seller e sua validação
type seller struct {
	Name  string
	CNPJ  string
	Email string
}

func validateSeller(s seller) bool {
	// Validar CNPJ: apenas verifica se o CNPJ tem 14 caracteres numéricos
	return s.Name != "" && len(s.CNPJ) == 14 && isValidEmail(s.Email)
}

// Testes

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
		t.Run(fmt.Sprintf("GenerateID_%d_%s", tc.index, tc.prefix), func(t *testing.T) {
			result := generateID(tc.index, tc.prefix)
			assert.Equal(t, tc.expected, result, "GenerateID failed for index: %d, prefix: %s", tc.index, tc.prefix)
		})
	}
}

func TestValidateProduct(t *testing.T) {
	validProduct := product{Name: "Product", Type: "Type", Quantity: 10}
	invalidProduct := product{Name: "", Type: "Type", Quantity: 0}

	t.Run("Valid Product", func(t *testing.T) {
		assert.True(t, validateProduct(validProduct))
	})

	t.Run("Invalid Product", func(t *testing.T) {
		assert.False(t, validateProduct(invalidProduct))
	})
}

func TestValidateCustomer(t *testing.T) {
	validCustomer := customer{Name: "Customer", Email: "test@example.com", CPF: "12345678900"}
	invalidCustomer := customer{Name: "", Email: "invalid", CPF: "123"}

	t.Run("Valid Customer", func(t *testing.T) {
		assert.True(t, validateCustomer(validCustomer))
	})

	t.Run("Invalid Customer", func(t *testing.T) {
		assert.False(t, validateCustomer(invalidCustomer))
	})
}

func TestValidateSeller(t *testing.T) {
	validSeller := seller{Name: "Seller", CNPJ: "12345678000190", Email: "seller@example.com"}
	invalidSeller := seller{Name: "", CNPJ: "invalid", Email: ""}

	t.Run("Valid Seller", func(t *testing.T) {
		assert.True(t, validateSeller(validSeller))
	})

	t.Run("Invalid Seller", func(t *testing.T) {
		assert.False(t, validateSeller(invalidSeller))
	})
}
