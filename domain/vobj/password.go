package vobj

import (
	"ToDo/config"
	"ToDo/packages/errors"
	"database/sql"
	"database/sql/driver"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Password string

func NewPassword(password string) (*Password, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), config.BcryptHashCost)
	if err != nil {
		return nil, errors.NewUnexpected(err)
	}

	value := Password(hashedPassword)
	return &value, nil
}

func (p *Password) IsValid(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(*p), []byte(password)) == nil
}

// sql

func (p *Password) Scan(value interface{}) error {
	nullString := &sql.NullString{}
	err := nullString.Scan(value)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	*p = Password(nullString.String)

	return nil
}

func (p *Password) Value() (driver.Value, error) {
	return string(*p), nil
}

// GormDataType gorm common data type
func (p *Password) GormDataType() string {
	return "password"
}

// GormDBDataType gorm db data type
func (p *Password) GormDBDataType(_ *gorm.DB, _ *schema.Field) string {
	return "longtext"
}

// json

func (p *Password) MarshalJSON() ([]byte, error) {
	return []byte("\"" + *p + "\""), nil
}

func (p *Password) UnmarshalJSON(b []byte) error {
	*p = Password(b)
	return nil
}
