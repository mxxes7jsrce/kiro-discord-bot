package bot

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CalcResult holds the result of a calculation.
type CalcResult struct {
	Expression string
	Result     float64
}

// EvaluateSimple evaluates a simple two-operand expression like "3 + 4".
// Supported operators: +, -, *, /, ^
func EvaluateSimple(expr string) (CalcResult, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return CalcResult{}, fmt.Errorf("empty expression")
	}

	ops := []string{"^", "*", "/", "+", "-"}
	for _, op := range ops {
		// find last occurrence to handle negative numbers on left side
		idx := strings.LastIndex(expr, " "+op+" ")
		if idx < 0 {
			continue
		}
		leftStr := strings.TrimSpace(expr[:idx])
		rightStr := strings.TrimSpace(expr[idx+len(op)+2:])

		left, err := strconv.ParseFloat(leftStr, 64)
		if err != nil {
			return CalcResult{}, fmt.Errorf("invalid left operand: %s", leftStr)
		}
		right, err := strconv.ParseFloat(rightStr, 64)
		if err != nil {
			return CalcResult{}, fmt.Errorf("invalid right operand: %s", rightStr)
		}

		var result float64
		switch op {
		case "+":
			result = left + right
		case "-":
			result = left - right
		case "*":
			result = left * right
		case "/":
			if right == 0 {
				return CalcResult{}, fmt.Errorf("division by zero")
			}
			result = left / right
		case "^":
			result = math.Pow(left, right)
		}
		return CalcResult{Expression: expr, Result: result}, nil
	}

	// Try parsing as a bare number
	v, err := strconv.ParseFloat(expr, 64)
	if err != nil {
		return CalcResult{}, fmt.Errorf("cannot evaluate expression: %s", expr)
	}
	return CalcResult{Expression: expr, Result: v}, nil
}

// FormatCalcResult formats a CalcResult for display.
func FormatCalcResult(r CalcResult) string {
	if r.Result == math.Trunc(r.Result) {
		return fmt.Sprintf("🧮 `%s` = **%d**", r.Expression, int64(r.Result))
	}
	return fmt.Sprintf("🧮 `%s` = **%.4g**", r.Expression, r.Result)
}
