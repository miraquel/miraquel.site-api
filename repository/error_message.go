package repository

import "errors"

func ErrDuplicate() error {
	return errors.New("record already exists")
}

func ErrNotExists() error {
	return errors.New("row does not exists")
}

func ErrUpdateFailed() error {
	return errors.New("update failed")
}

func ErrDeleteFailed() error {
	return errors.New("delete failed")
}
