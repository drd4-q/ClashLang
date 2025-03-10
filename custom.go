package main
package custom

type CustomModule struct{}

func NewCustomModule() *CustomModule {
	return &CustomModule{}
}

func (c *CustomModule) HandleDef(cmd Command, params map[string]string, functions map[string][]string) {
	funcName := params["name"]
	functions[funcName] = []string{}
}

func (c *CustomModule) HandleFunctionCall(cmd Command, params map[string]string, functions map[string][]string, execute func(string)) {
	funcName := params["func"]
	if cmds, ok := functions[funcName]; ok {
		for _, cmd := range cmds {
			execute(cmd)
		}
	}
}

func (c *CustomModule) HandleJump(cmd Command, params map[string]string, functions map[string][]string, execute func(string)) {
	funcName := params["func"]
	if cmds, ok := functions[funcName]; ok {
		for _, cmd := range cmds {
			execute(cmd)
		}
	}
}
