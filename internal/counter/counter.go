package counter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	maxContador = 9999999
)

func Avanzar(actual int) int {
	actual++
	if actual > maxContador {
		return 1
	}
	return actual
}

func LeerUltimo(configPath string) int {
	if configPath == "" {
		return 0
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return 0
	}

	num, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return num
}

func GuardarUltimo(configPath string, num int) error {
	if configPath == "" {
		return fmt.Errorf("configPath no inicializado")
	}
	return os.WriteFile(configPath, fmt.Appendf(nil, "%d", num), 0644)
}