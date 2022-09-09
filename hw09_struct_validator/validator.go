package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrInvalidLen    = errors.New("invalid Length")
	ErrInvalidIn     = errors.New("invalid In")
	ErrInvalidMax    = errors.New("invalid Max")
	ErrInvalidMin    = errors.New("invalid Min")
	ErrInvalidRegexp = errors.New("invalid Regexp")
)

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for _, err := range v {
		fmt.Fprintf(&sb, "[f: %s, e: %v] ", err.Field, err.Err)
	}

	return sb.String()
}

type StructValidator interface {
	Validate(interface{}) []error
}

type Validator struct{}

type IntValidation struct {
	min int64
	max int64
	in  []int64
}

type StringValidation struct {
	len    int64
	regexp string
	in     []string
}

// PrepareIntValidation разделяем числовой validate тэг на контрольные значения min/max/in.
func (i Validator) PrepareIntValidation(tag string) (*IntValidation, []error) {
	terms := strings.Split(tag, "|")
	var valTerms IntValidation
	var validationErrors []error
	for _, term := range terms {
		splitTag := strings.Split(term, ":")
		if len(splitTag) < 2 {
			continue
		}
		tagExp := splitTag[0]
		tagValue := splitTag[1]

		switch {
		case tagExp == "min":
			val, err := strconv.Atoi(tagValue)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
			valTerms.min = int64(val)

		case tagExp == "max":
			val, err := strconv.Atoi(tagValue)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
			valTerms.max = int64(val)

		case tagExp == "in":
			var inValues []int64
			for _, val := range strings.Split(tagValue, ",") {
				intVal, err := strconv.Atoi(val)
				if err != nil {
					validationErrors = append(validationErrors, err)
					continue
				}
				inValues = append(inValues, int64(intVal))
			}
			valTerms.in = inValues
		}
	}
	return &valTerms, validationErrors
}

// PrepareStringValidation разделяем текстовый validate тэг на контрольные значения len/regexp/in.
func (i Validator) PrepareStringValidation(tag string) (*StringValidation, []error) {
	terms := strings.Split(tag, "|")
	var valTerms StringValidation
	var validationErrors []error
	for _, term := range terms {
		splitTag := strings.Split(term, ":")
		if len(splitTag) < 2 {
			continue
		}
		tagExp := splitTag[0]
		tagValue := splitTag[1]

		switch {
		case tagExp == "len":
			val, err := strconv.Atoi(tagValue)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
			valTerms.len = int64(val)

		case tagExp == "regexp":
			valTerms.regexp = tagValue

		case tagExp == "in":
			var inValues []string
			inValues = append(inValues, strings.Split(tagValue, ",")...)
			valTerms.in = inValues
		}
	}
	return &valTerms, validationErrors
}

// ValidateValue определяем validate тэг по типу int/string и проверяем значения полей.
func (i Validator) ValidateValue(fieldValue reflect.Value, fieldType reflect.StructField, vErr *ValidationErrors) {
	// Поле типа int?
	if fieldValue.Kind() == reflect.Int {
		valTerms, PrepareErrors := i.PrepareIntValidation(fieldType.Tag.Get("validate"))
		for _, err := range PrepareErrors {
			*vErr = append(*vErr, ValidationError{
				Field: fieldType.Name,
				Err:   err,
			})
		}
		// Проверяем int по заданным диапазонам.
		valTerms.validateMin(fieldValue.Int(), ValidationError{Field: fieldType.Name, Err: nil}, vErr)
		valTerms.validateMax(fieldValue.Int(), ValidationError{Field: fieldType.Name, Err: nil}, vErr)
		valTerms.validateIn(fieldValue.Int(), ValidationError{Field: fieldType.Name, Err: nil}, vErr)
		// Поле типа string?
	} else if fieldValue.Kind() == reflect.String {
		valTerms, PrepareErrors := i.PrepareStringValidation(fieldType.Tag.Get("validate"))
		for _, err := range PrepareErrors {
			*vErr = append(*vErr, ValidationError{
				Field: fieldType.Name,
				Err:   err,
			})
		}
		// Проверяем string по заданным критериям.
		valTerms.validateLen(fieldValue.String(), ValidationError{Field: fieldType.Name, Err: nil}, vErr)
		valTerms.validateRegexp(fieldValue.String(), ValidationError{Field: fieldType.Name, Err: nil}, vErr)
		valTerms.validateIn(fieldValue.String(), ValidationError{Field: fieldType.Name, Err: nil}, vErr)
	}
}

// Validate основная функция валидации структуры с разделеним по тегам.
func (i Validator) Validate(structToValidate interface{}) error {
	var vErr ValidationErrors
	// Анализируем поля структуры.
	Value := reflect.ValueOf(structToValidate)
	valueType := Value.Type()
	// Проверям тип/kind каждого поля и запускаем проверку на тип int/string.
	for d := 0; d < valueType.NumField(); d++ {
		if valueType.Field(d).Tag.Get("validate") != "" && Value.Field(d).Kind() != reflect.Slice {
			i.ValidateValue(Value.Field(d), valueType.Field(d), &vErr)
		} else if Value.Field(d).Kind() == reflect.Slice {
			for sl := 0; sl < Value.Field(d).Len(); sl++ {
				i.ValidateValue(Value.Field(d).Index(sl), valueType.Field(d), &vErr)
			}
		}
	}
	if len(vErr) != 0 {
		return vErr
	}
	return vErr
}

// validateMin проверка нижней границы int диапазона с обновлением списка ошибок.
func (i IntValidation) validateMin(val int64, vErr ValidationError, vErrs *ValidationErrors) {
	if val < i.min && i.min > 0 {
		vErr.Err = ErrInvalidMin
		*vErrs = append(*vErrs, vErr)
	}
}

// validateMax проверка верхней границы int диапазона с обновлением списка ошибок.
func (i IntValidation) validateMax(val int64, vErr ValidationError, vErrs *ValidationErrors) {
	if val > i.max && i.max > 0 {
		vErr.Err = ErrInvalidMax
		*vErrs = append(*vErrs, vErr)
	}
}

// validateIn проверка вхождения int в диапазон значений с обновлением списка ошибок.
func (i IntValidation) validateIn(val int64, vErr ValidationError, vErrs *ValidationErrors) {
	var ValueIn bool
	if len(i.in) < 1 {
		return
	}
	for _, value := range i.in {
		if value == val {
			ValueIn = true
			break
		}
	}
	if !ValueIn {
		vErr.Err = ErrInvalidIn
		*vErrs = append(*vErrs, vErr)
	}
}

// validateLen проверка длины строки с обновлением списка ошибок.
func (i StringValidation) validateLen(val string, vErr ValidationError, vErrs *ValidationErrors) {
	if int64(len(val)) != i.len && i.len > 0 {
		vErr.Err = ErrInvalidLen
		*vErrs = append(*vErrs, vErr)
	}
}

// validateRegexp проверка соответствия строки регулярному выражению с обновлением списка ошибок.
func (i StringValidation) validateRegexp(val string, vErr ValidationError, vErrs *ValidationErrors) {
	matched, err := regexp.Match(i.regexp, []byte(val))
	if err != nil {
		return
	}
	if !matched {
		vErr.Err = ErrInvalidRegexp
		*vErrs = append(*vErrs, vErr)
	}
}

// validateIn проверка вхождения строки в диапазон значений с обновлением списка ошибок.
func (i StringValidation) validateIn(val string, vErr ValidationError, vErrs *ValidationErrors) {
	var ValueIn bool
	if len(i.in) < 1 {
		return
	}
	for _, validVal := range i.in {
		if validVal == val {
			ValueIn = true
			break
		}
	}
	if !ValueIn {
		vErr.Err = ErrInvalidIn
		*vErrs = append(*vErrs, vErr)
	}
}
