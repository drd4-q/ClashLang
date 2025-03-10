package strings

type StringModule struct{}

func NewStringModule() *StringModule {
	return &StringModule{}
}

func (s *StringModule) HandleLen(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(string); ok {
		*lastResult = utf8.RuneCountInString(val)
	} else {
		*lastResult = 0
		fmt.Println("Ошибка: переменная не является строкой")
	}
}

func (s *StringModule) HandleSubstr(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	startStr := params["start"]
	lengthStr := params["length"]
	start, _ := strconv.Atoi(startStr)
	length, _ := strconv.Atoi(lengthStr)
	if val, ok := variables[varName].(string); ok {
		if start >= 0 && start < len(val) && length >= 0 {
			if start+length > len(val) {
				length = len(val) - start
			}
			*lastResult = val[start : start+length]
		} else {
			*lastResult = ""
			fmt.Println("Ошибка: неверные индексы")
		}
	} else {
		*lastResult = ""
		fmt.Println("Ошибка: переменная не является строкой")
	}
}

func (s *StringModule) HandleFind(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	strVar := params["str"]
	subVar := params["sub"]
	if str, ok := variables[strVar].(string); ok {
		if sub, ok := variables[subVar].(string); ok {
			*lastResult = strings.Index(str, sub)
		} else {
			*lastResult = -1
			fmt.Println("Ошибка: подстрока не является строкой")
		}
	} else {
		*lastResult = -1
		fmt.Println("Ошибка: строка не является строкой")
	}
}

func (s *StringModule) HandleReplace(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	strVar := params["str"]
	oldVar := params["old"]
	newVar := params["new"]
	if str, ok := variables[strVar].(string); ok {
		if old, ok := variables[oldVar].(string); ok {
			if new, ok := variables[newVar].(string); ok {
				*lastResult = strings.Replace(str, old, new, -1)
			} else {
				*lastResult = str
				fmt.Println("Ошибка: новое значение не является строкой")
			}
		} else {
			*lastResult = str
			fmt.Println("Ошибка: старое значение не является строкой")
		}
	} else {
		*lastResult = ""
		fmt.Println("Ошибка: строка не является строкой")
	}
}

func (s *StringModule) HandleSplit(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	strVar := params["str"]
	sepVar := params["sep"]
	if str, ok := variables[strVar].(string); ok {
		if sep, ok := variables[sepVar].(string); ok {
			*lastResult = strings.Split(str, sep)
		} else {
			*lastResult = []string{str}
			fmt.Println("Ошибка: разделитель не является строкой")
		}
	} else {
		*lastResult = []string{}
		fmt.Println("Ошибка: строка не является строкой")
	}
}

func (s *StringModule) HandleJoin(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	sliceVar := params["slice"]
	sepVar := params["sep"]
	if slice, ok := variables[sliceVar].([]string); ok {
		if sep, ok := variables[sepVar].(string); ok {
			*lastResult = strings.Join(slice, sep)
		} else {
			*lastResult = strings.Join(slice, "")
			fmt.Println("Ошибка: разделитель не является строкой")
		}
	} else {
		*lastResult = ""
		fmt.Println("Ошибка: переменная не является срезом строк")
	}
}

func (s *StringModule) HandleLower(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(string); ok {
		*lastResult = strings.ToLower(val)
	} else {
		*lastResult = ""
		fmt.Println("Ошибка: переменная не является строкой")
	}
}

func (s *StringModule) HandleText(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	expr := params["expr"]
	*lastResult = parseTextExpression(expr, variables)
}

func (s *StringModule) HandleTextOut(cmd Command, params map[string]string, variables map[string]interface{}, lastResult interface{}) {
	varName := params["var"]
	variables[varName] = lastResult
}

func (s *StringModule) HandleTextLength(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(string); ok {
		*lastResult = utf8.RuneCountInString(val)
	} else {
		*lastResult = 0
		fmt.Println("Ошибка: переменная не является текстом")
	}
}

func (s *StringModule) HandleTextUpper(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	if val, ok := variables[varName].(string); ok {
		*lastResult = strings.ToUpper(val)
	} else {
		*lastResult = ""
		fmt.Println("Ошибка: переменная не является текстом")
	}
}