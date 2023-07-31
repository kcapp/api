package models

import (
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

// TestIsCheckout will check that the given visit is checkout
func TestIsCheckout(t *testing.T) {
	visit := Visit{FirstDart: &Dart{Value: null.IntFrom(20), Multiplier: 1},
		SecondDart: &Dart{Value: null.IntFrom(20), Multiplier: 1},
		ThirdDart:  &Dart{Value: null.IntFrom(20), Multiplier: 1}}
	assert.Equal(t, visit.IsCheckout(60, OUTSHOTANY), true, "should be checkout")

	visit = Visit{FirstDart: &Dart{Value: null.IntFrom(20), Multiplier: 1},
		SecondDart: &Dart{Value: null.IntFrom(20), Multiplier: 1},
		ThirdDart:  &Dart{Value: null.IntFrom(20), Multiplier: 1}}
	assert.Equal(t, visit.IsCheckout(60, OUTSHOTDOUBLE), false, "should not be checkout")

	visit = Visit{FirstDart: &Dart{Value: null.IntFrom(20), Multiplier: 3}, SecondDart: &Dart{}, ThirdDart: &Dart{}}
	assert.Equal(t, visit.IsCheckout(60, OUTSHOTMASTER), true, "should be checkout")

	visit = Visit{FirstDart: &Dart{}, SecondDart: &Dart{Value: null.IntFrom(20), Multiplier: 2}, ThirdDart: &Dart{}}
	assert.Equal(t, visit.IsCheckout(40, OUTSHOTDOUBLE), true, "should be checkout")
}

// TestIsCheckout_Master will check that the given visit is checkout
func TestIsCheckout_Master(t *testing.T) {
	visit := Visit{FirstDart: &Dart{Value: null.IntFrom(20), Multiplier: 2},
		SecondDart: &Dart{Value: null.IntFrom(10), Multiplier: 1},
		ThirdDart:  &Dart{Value: null.IntFrom(5), Multiplier: 2}}
	assert.Equal(t, visit.IsCheckout(60, OUTSHOTMASTER), true, "should be checkout")

	visit = Visit{FirstDart: &Dart{Value: null.IntFrom(14), Multiplier: 1}, SecondDart: &Dart{Value: null.IntFrom(7), Multiplier: 3}, ThirdDart: &Dart{}}
	assert.Equal(t, visit.IsCheckout(35, OUTSHOTMASTER), true, "should be checkout")

	visit = Visit{FirstDart: &Dart{Value: null.IntFrom(3), Multiplier: 3}, SecondDart: &Dart{}, ThirdDart: &Dart{}}
	assert.Equal(t, visit.IsCheckout(9, OUTSHOTMASTER), true, "should be checkout")
}
