package formula

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormula_checkValid(t *testing.T) {
	s := "()*+/0123456789a"
	require.EqualError(t, checkValid(s), "Invalid string")

	s = "()*+-/0123456789"
	require.NoError(t, checkValid(s))

	s = "())*+-/0123456789"
	require.EqualError(t, checkValid(s), "formal error")

}

func TestFormulaCalculator(t *testing.T) {
	formulaStr := "2+3*4"
	result, err := Calculator(formulaStr)
	require.Equal(t, result, float32(14))

	formulaStr = "2*3+4"
	result, err = Calculator(formulaStr)
	require.Equal(t, result, float32(10))

	formulaStr = "(2+13)*4-44/(2-2)"
	result, err = Calculator(formulaStr)
	require.EqualError(t, err, "division by zero")

	formulaStr = "(2+3)*4"
	result, err = Calculator(formulaStr)
	require.Equal(t, result, float32(20))

	formulaStr = "(1+(3-2))*4"
	result, err = Calculator(formulaStr)
	require.Equal(t, result, float32(8))

	formulaStr = "(1+(3-2))*((4-2)*2)"
	result, err = Calculator(formulaStr)
	require.Equal(t, result, float32(8))

	formulaStr = "10-(1+(3-2))*((4-2)*2)"
	result, err = Calculator(formulaStr)
	require.Equal(t, result, float32(2))

	formulaStr = "10-(1+(3-1*2))*((4-1-1)*2)"
	result, err = Calculator(formulaStr)
	require.Equal(t, result, float32(2))

	formulaStr = "(1/3+2/3)"
	result, err = Calculator(formulaStr)
	fmt.Println(result)

}

func BenchmarkFormulaCalculator(b *testing.B) {
	formulaStr := "10-(1+(3-1*2))*((4-1-1)*2)"
	for i := 0; i < b.N; i++ {
		Calculator(formulaStr)
	}
}
