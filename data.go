package data

type DataModule struct{}

func NewDataModule() *DataModule {
	return &DataModule{}
}

func (d *DataModule) HandleArrayCreate(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	sizeStr := params["size"]
	size, _ := strconv.Atoi(sizeStr)
	*lastResult = make([]interface{}, size)
}

func (d *DataModule) HandleArraySet(cmd Command, params map[string]string, variables map[string]interface{}) {
	arrayVar := params["array"]
	indexStr := params["index"]
	valueVar := params["value"]
	index, _ := strconv.Atoi(indexStr)
	if array, ok := variables[arrayVar].([]interface{}); ok {
		if index >= 0 && index < len(array) {
			if val, ok := variables[valueVar]; ok {
				array[index] = val
			}
		}
	}
}

func (d *DataModule) HandleArrayGet(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	arrayVar := params["array"]
	indexStr := params["index"]
	index, _ := strconv.Atoi(indexStr)
	if array, ok := variables[arrayVar].([]interface{}); ok {
		if index >= 0 && index < len(array) {
			*lastResult = array[index]
		} else {
			*lastResult = nil
			fmt.Println("Ошибка: индекс вне диапазона")
		}
	}
}

func (d *DataModule) HandleListCreate(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	*lastResult = []interface{}{}
}

func (d *DataModule) HandleListAppend(cmd Command, params map[string]string, variables map[string]interface{}) {
	listVar := params["list"]
	valueVar := params["value"]
	if list, ok := variables[listVar].([]interface{}); ok {
		if val, ok := variables[valueVar]; ok {
			variables[listVar] = append(list, val)
		}
	}
}

func (d *DataModule) HandleListGet(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	listVar := params["list"]
	indexStr := params["index"]
	index, _ := strconv.Atoi(indexStr)
	if list, ok := variables[listVar].([]interface{}); ok {
		if index >= 0 && index < len(list) {
			*lastResult = list[index]
		} else {
			*lastResult = nil
			fmt.Println("Ошибка: индекс вне диапазона")
		}
	}
}

func (d *DataModule) HandleDictCreate(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	*lastResult = make(map[string]interface{})
}

func (d *DataModule) HandleDictSet(cmd Command, params map[string]string, variables map[string]interface{}) {
	dictVar := params["dict"]
	keyVar := params["key"]
	valueVar := params["value"]
	if dict, ok := variables[dictVar].(map[string]interface{}); ok {
		if key, ok := variables[keyVar].(string); ok {
			if val, ok := variables[valueVar]; ok {
				dict[key] = val
			}
		}
	}
}

func (d *DataModule) HandleDictGet(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	dictVar := params["dict"]
	keyVar := params["key"]
	if dict, ok := variables[dictVar].(map[string]interface{}); ok {
		if key, ok := variables[keyVar].(string); ok {
			if val, exists := dict[key]; exists {
				*lastResult = val
			} else {
				*lastResult = nil
				fmt.Println("Ошибка: ключ не найден")
			}
		}
	}
}