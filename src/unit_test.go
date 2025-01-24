package main

import (
	"testing"

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
