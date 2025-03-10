package math

import (
	"math"
	"math/rand"
)

type MathModule struct{}

func NewMathModule() *MathModule {
	return &MathModule{}
}

func (m *MathModule) HandleAbs(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Abs(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandleSqrt(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Sqrt(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandlePow(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	baseStr := params["base"]
	expStr := params["exponent"]
	base, _ := strconv.ParseFloat(baseStr, 64)
	exp, _ := strconv.ParseFloat(expStr, 64)
	*lastResult = math.Pow(base, exp)
}

func (m *MathModule) HandleRound(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Round(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandleSin(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Sin(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandleCos(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Cos(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandleTan(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Tan(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandleLog(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Log(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandleLog10(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(float64); ok {
		*lastResult = math.Log10(val)
	} else {
		*lastResult = 0.0
		fmt.Println("Ошибка: переменная не является числом")
	}
}

func (m *MathModule) HandleRandom(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	*lastResult = rand.Float64()
}

func (m *MathModule) HandleRandInt(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	minStr := params["min"]
	maxStr := params["max"]
	min, _ := strconv.Atoi(minStr)
	max, _ := strconv.Atoi(maxStr)
	*lastResult = rand.Intn(max-min+1) + min
}

func (m *MathModule) HandleSolve(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	expr := params["expr"]
	*lastResult = parseExpression(expr, variables)
}

func (m *MathModule) HandleSolveOut(cmd Command, params map[string]string, variables map[string]interface{}, lastResult interface{}) {
	varName := params["var"]
	variables[varName] = lastResult
}