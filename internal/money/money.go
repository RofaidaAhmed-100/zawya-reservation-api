package money

import (
	"errors"
	"fmt"
)

type Currency string

const (
	EGP Currency = "EGP"
	USD Currency = "USD"
	EUR Currency = "EUR"
	SAR Currency = "SAR"
	KWD Currency = "KWD"
	JPY Currency = "JPY"
)

type CurrencyInfo struct {
	Code          Currency
	MinorUnits    int
	MinorUnitName string
}

var currencyInfo = map[Currency]CurrencyInfo{
	EGP: {Code: EGP, MinorUnits: 2, MinorUnitName: "piasters"},
	USD: {Code: USD, MinorUnits: 2, MinorUnitName: "cents"},
	EUR: {Code: EUR, MinorUnits: 2, MinorUnitName: "cents"},
	SAR: {Code: SAR, MinorUnits: 2, MinorUnitName: "halalas"},
	KWD: {Code: KWD, MinorUnits: 3, MinorUnitName: "fils"},
	JPY: {Code: JPY, MinorUnits: 0, MinorUnitName: ""},
}

type Money struct {
	Amount   uint64   `json:"amount"`
	Currency Currency `json:"currency"`
}

func New(amount uint64, currency Currency) *Money {
	return &Money{
		Amount:   amount,
		Currency: currency,
	}
}

func NewFromFloat(amount float64, currency Currency) (*Money, error) {
	info, exists := currencyInfo[currency]
	if !exists {
		return nil, errors.New("unsupported currency")
	}

	multiplier := uint64(1)
	for i := 0; i < info.MinorUnits; i++ {
		multiplier *= 10
	}

	smallestUnit := uint64(amount*float64(multiplier) + 0.5)

	return &Money{
		Amount:   smallestUnit,
		Currency: currency,
	}, nil
}

func (m *Money) ToFloat() float64 {
	info := currencyInfo[m.Currency]
	divisor := float64(1)
	for i := 0; i < info.MinorUnits; i++ {
		divisor *= 10
	}
	return float64(m.Amount) / divisor
}

func (m *Money) String() string {
	info := currencyInfo[m.Currency]
	if info.MinorUnits == 0 {
		return fmt.Sprintf("%d %s", m.Amount, m.Currency)
	}
	return fmt.Sprintf("%."+fmt.Sprint(info.MinorUnits)+"f %s", m.ToFloat(), m.Currency)
}

func (m *Money) Add(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, errors.New("cannot add different currencies")
	}

	sum := m.Amount + other.Amount
	if sum < m.Amount {
		return nil, errors.New("addition overflow")
	}

	return &Money{
		Amount:   sum,
		Currency: m.Currency,
	}, nil
}
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, errors.New("cannot subtract different currencies")
	}

	if m.Amount < other.Amount {
		return nil, errors.New("insufficient funds")
	}

	return &Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}, nil
}

func (m *Money) Multiply(factor uint64) (*Money, error) {
	result := m.Amount * factor
	if factor != 0 && result/factor != m.Amount {
		return nil, errors.New("multiplication overflow")
	}

	return &Money{
		Amount:   result,
		Currency: m.Currency,
	}, nil
}

func (m *Money) MultiplyFloat(factor float64) (*Money, error) {
	if factor < 0 {
		return nil, errors.New("factor cannot be negative")
	}

	result := uint64(float64(m.Amount)*factor + 0.5)

	return &Money{
		Amount:   result,
		Currency: m.Currency,
	}, nil
}

func (m *Money) Divide(divisor uint64) (*Money, error) {
	if divisor == 0 {
		return nil, errors.New("division by zero")
	}

	return &Money{
		Amount:   m.Amount / divisor,
		Currency: m.Currency,
	}, nil
}

func (m *Money) Equals(other *Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

func (m *Money) GreaterThan(other *Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, errors.New("cannot compare different currencies")
	}
	return m.Amount > other.Amount, nil
}

func (m *Money) LessThan(other *Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, errors.New("cannot compare different currencies")
	}
	return m.Amount < other.Amount, nil
}
