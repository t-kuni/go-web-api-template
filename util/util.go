package util

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

func Ptr[T any](v T) *T {
	return &v
}

func UuidToStrfmtUUID(u uuid.UUID) (strfmt.UUID, error) {
	var strfmtUUID strfmt.UUID
	if err := strfmtUUID.Scan(u.String()); err != nil {
		return "", err
	}
	return strfmtUUID, nil
}

func StringToStrfmtUUID(s string) (strfmt.UUID, error) {
	var strfmtUUID strfmt.UUID
	if err := strfmtUUID.Scan(s); err != nil {
		return "", err
	}
	return strfmtUUID, nil
}
