package control

type ControlModule struct{}

func NewControlModule() *ControlModule {
	return &ControlModule{}
}

func (c *ControlModule) HandleFor(cmd Command, params map[string]string, variables map[string]interface{}, execute func(string)) {
	varName := params["var"]
	startStr := params["start"]
	endStr := params["end"]
	start, _ := strconv.Atoi(startStr)
	end, _ := strconv.Atoi(endStr)
	for i := start; i <= end; i++ {
		variables[varName] = i
		execute("{")
	}
}

func (c *ControlModule) HandleWhile(cmd Command, params map[string]string, variables map[string]interface{}, execute func(string), inIfBlock *bool) {
	condVar := params["var"]
	condValueStr := params["value"]
	condValue, _ := strconv.Atoi(condValueStr)
	for {
		if val, ok := variables[condVar].(int); ok && val == condValue {
			execute("{")
			break
		}
	}
}

func (c *ControlModule) HandleDo(cmd Command, params map[string]string, variables map[string]interface{}, execute func(string), inIfBlock *bool) {
	// Пустая реализация, требует доработки
}

func (c *ControlModule) HandleWhileDo(cmd Command, params map[string]string, variables map[string]interface{}, execute func(string), inIfBlock *bool) {
	condVar := params["var"]
	condValueStr := params["value"]
	condValue, _ := strconv.Atoi(condValueStr)
	for {
		execute("{")
		if val, ok := variables[condVar].(int); ok && val == condValue {
			break
		}
	}
}

func (c *ControlModule) HandleSwitch(cmd Command, params map[string]string, variables map[string]interface{}, execute func(string), inIfBlock *bool) {
	*inIfBlock = true
}

func (c *ControlModule) HandleCase(cmd Command, params map[string]string, variables map[string]interface{}, execute func(string), inIfBlock *bool) {
	if *inIfBlock {
		valueStr := params["value"]
		if val, ok := variables["switch_var"].(int); ok {
			caseValue, _ := strconv.Atoi(valueStr)
			if val == caseValue {
				execute("{")
			}
		}
	}
}

func (c *ControlModule) HandleDefault(cmd Command, params map[string]string, variables map[string]interface{}, execute func(string), inIfBlock *bool) {
	if *inIfBlock {
		execute("{")
	}
}

func (c *ControlModule) HandleIf(cmd Command, params map[string]string, variables map[string]interface{}, inIfBlock *bool, ifBlock *[]string) {
	varName := params["var"]
	valueStr := params["value"]
	value, _ := strconv.Atoi(valueStr)
	if val, ok := variables[varName].(int); ok && val == value {
		*inIfBlock = true
	}
}