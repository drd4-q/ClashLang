package system

import (
	"os"
	"time"
)

type SystemModule struct{}

func NewSystemModule() *SystemModule {
	return &SystemModule{}
}

func (s *SystemModule) HandleTime(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	*lastResult = time.Now().Format("15:04:05")
}

func (s *SystemModule) HandleDate(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	*lastResult = time.Now().Format("2006-01-02")
}

func (s *SystemModule) HandleEnv(cmd Command, params map[string]string, variables map[string]interface{}, lastResult *interface{}) {
	varName := params["var"]
	*lastResult = os.Getenv(varName)
}