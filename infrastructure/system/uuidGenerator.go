package system

import (
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
)

type UuidGenerator struct {
}

func NewUuidGenerator() system.IUuidGenerator {
	return &UuidGenerator{}
}

func (u UuidGenerator) Generate() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", eris.Wrap(err, "")
	}
	return id.String(), nil
}
