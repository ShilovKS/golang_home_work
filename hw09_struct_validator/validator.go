package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError представляет ошибку валидации конкретного поля.
type ValidationError struct {
	Field string
	Err   error
}

// ValidationErrors представляет набор ошибок валидации.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	errs := make([]string, len(ve))
	for i, v := range ve {
		errs[i] = fmt.Sprintf("%s: %s", v.Field, v.Err)
	}
	return strings.Join(errs, ", ")
}

// CommonError используется для ошибок программирования (например, неверный формат тега).
type CommonError struct {
	Message string
}

func (e CommonError) Error() string {
	return e.Message
}

// Фабричные функции для создания ошибок.
func NewErrParsing(def, rule string) error {
	return CommonError{Message: fmt.Sprintf("can't parse %s: %s", def, rule)}
}

func NewErrInvalidRuleFormat(rule string) error {
	return CommonError{Message: fmt.Sprintf("invalid rule format: %s", rule)}
}

func NewErrUnsupportedType(kind reflect.Kind) error {
	return CommonError{Message: fmt.Sprintf("unsupported type: %v", kind)}
}

func NewErrUnsupportedValidator(validatorType string) error {
	return CommonError{Message: fmt.Sprintf("unsupported validator: %s", validatorType)}
}

func NewErrValueOutOfRange(comp string, val int) error {
	return fmt.Errorf("value must be %s than or equal to %d", comp, val)
}

func NewErrValueNotInSet(set string) error {
	return fmt.Errorf("value must be one of %s", set)
}

func NewErrStringLengthMismatch(expected int) error {
	return fmt.Errorf("string length must be exactly %d", expected)
}

func NewErrStringDoesNotMatchPattern(pattern string) error {
	return fmt.Errorf("string does not match pattern: %s", pattern)
}

var ErrInput = errors.New("input must be a struct")

// Validate проверяет публичные поля входной структуры v, используя тег `validate`.
// Если v не является структурой, возвращается ErrInput.
// Если для поля указан тег validate, то для каждого правила (разделённых символом '|')
// применяется соответствующий валидатор.
// Все ошибки валидации накапливаются и возвращаются в виде ValidationErrors.
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if val.Kind() != reflect.Struct {
		return ErrInput
	}

	var validationErrors ValidationErrors

	// Перебираем все поля структуры.
	for i := 0; i < typ.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")

		if fieldType.PkgPath != "" || tag == "" {
			continue
		}

		// Разбиваем тег по символу '|' для объединения валидаторов.
		rules := strings.Split(tag, "|")
		for _, rule := range rules {
			err := applyValidation(rule, fieldVal)
			// Если обнаружена программная ошибка – немедленно вернуть её.
			var ce CommonError
			if errors.As(err, &ce) {
				return ce
			}
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: fieldType.Name, Err: err})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

// applyValidation применяет одно правило валидации (формата "validator:parameter")
// к значению поля field.
func applyValidation(rule string, field reflect.Value) error {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return NewErrInvalidRuleFormat(rule)
	}

	validatorType := parts[0]
	validatorParam := parts[1]

	switch field.Kind() {
	case reflect.Int:
		return validateInt(int(field.Int()), validatorType, validatorParam)
	case reflect.String:
		return validateString(field.String(), validatorType, validatorParam)
	case reflect.Slice:
		// Поддерживаем слайсы int и string.
		switch field.Type().Elem().Kind() {
		case reflect.Int:
			for i := 0; i < field.Len(); i++ {
				if err := validateInt(int(field.Index(i).Int()), validatorType, validatorParam); err != nil {
					return err
				}
			}
		case reflect.String:
			for i := 0; i < field.Len(); i++ {
				if err := validateString(field.Index(i).String(), validatorType, validatorParam); err != nil {
					return err
				}
			}
		default:
			return NewErrUnsupportedType(field.Kind())
		}
	default:
		return NewErrUnsupportedType(field.Kind())
	}
	return nil
}

// validateString применяет валидаторы для строк.
func validateString(value string, validatorType, validatorParam string) error {
	switch validatorType {
	case "len":
		expectedLen, err := strconv.Atoi(validatorParam)
		if err != nil {
			return NewErrParsing("len", validatorParam)
		}
		if len(value) != expectedLen {
			return NewErrStringLengthMismatch(expectedLen)
		}
	case "regexp":
		re, err := regexp.Compile(validatorParam)
		if err != nil {
			return NewErrParsing("regexp", validatorParam)
		}
		if !re.MatchString(value) {
			return NewErrStringDoesNotMatchPattern(validatorParam)
		}
	case "in":
		inSet := strings.Split(validatorParam, ",")
		found := false
		for _, v := range inSet {
			if value == v {
				found = true
				break
			}
		}
		if !found {
			return NewErrValueNotInSet(validatorParam)
		}
	default:
		return NewErrUnsupportedValidator(validatorType)
	}
	return nil
}

// validateInt применяет валидаторы для чисел.
func validateInt(value int, validatorType, validatorParam string) error {
	switch validatorType {
	case "min":
		minVal, err := strconv.Atoi(validatorParam)
		if err != nil {
			return NewErrParsing("min", validatorParam)
		}
		if value < minVal {
			return NewErrValueOutOfRange("greater", minVal)
		}
	case "max":
		maxVal, err := strconv.Atoi(validatorParam)
		if err != nil {
			return NewErrParsing("max", validatorParam)
		}
		if value > maxVal {
			return NewErrValueOutOfRange("less", maxVal)
		}
	case "in":
		inSet := strings.Split(validatorParam, ",")
		found := false
		for _, v := range inSet {
			num, err := strconv.Atoi(v)
			if err != nil {
				return NewErrParsing("in", validatorParam)
			}
			if value == num {
				found = true
				break
			}
		}
		if !found {
			return NewErrValueNotInSet(validatorParam)
		}
	default:
		return NewErrUnsupportedValidator(validatorType)
	}
	return nil
}
