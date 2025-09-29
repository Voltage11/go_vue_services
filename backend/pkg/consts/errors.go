package consts

import "errors"

var (
	ErrNotFound       = errors.New("не найдено")
	ErrBadData        = errors.New("некорректные данные")
	ErrAlreadyExists  = errors.New("уже существует")
	ErrNoRowsAffected = errors.New("ни одна запись не была обработана")
)
