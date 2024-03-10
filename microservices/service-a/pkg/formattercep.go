package pkg

import (
	"fmt"
	"regexp"
)

func CepFormatted(input *string) bool {
	matched, err := regexp.MatchString(`^[0-9]{8}$`, *input)
	if err != nil {
		fmt.Printf("Erro ao verificar a string: %v\n", err)
		return false
	}
	return matched
}
