package validators

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/ronaldalds/base-go-api/internal/handlers"
	"github.com/ronaldalds/base-go-api/internal/i18n"
)

type Validator struct {
	*validator.Validate
}

func NewValidator() *Validator {
	var validate = validator.New()
	return &Validator{
		Validate: validate,
	}
}
func (v *Validator) ValidateStruct(data interface{}) *handlers.ErrHandler {
	// Verifica se o objeto possui erros de validação
	err := v.Struct(data)
	if err != nil {
		// Converte o erro para ValidationErrors, se aplicável
		errors := handlers.NewError()
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				fieldName := err.Field()
				tag := err.Tag()
				errors.AddDetailErr(fieldName, fmt.Sprintf(i18n.ERR_INVALID_FIELD, tag))
			}
			return errors
		}
		// Retorna erro genérico se não for ValidationErrors
		errors.AddDetailErr("error", i18n.ERR_VALIDATE_STRUCTURE_DEFAULT)
		return errors
	}
	// Nenhum erro encontrado
	return nil
}

func (v *Validator) ValidatePassword(password string) *handlers.ErrHandler {
	errors := handlers.NewError()

	// Verificar se contém uma letra maiúscula
	hasUpper := false
	hasSymbol := false

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsSymbol(r) || unicode.IsPunct(r) { // Símbolos e pontuações
			hasSymbol = true
		}
	}

	// Adicionar erros caso os critérios não sejam atendidos
	if !hasUpper {
		errors.AddDetailErr("uppercase", i18n.ERR_PASSWORD_INVALID_UPPERCASE)
	}
	if !hasSymbol {
		errors.AddDetailErr("symbol", i18n.ERR_PASSWORD_INVALID_SYMBOL)
	}

	// Retornar nil se não houver erros
	if hasSymbol || hasUpper {
		return nil
	}

	return errors
}
