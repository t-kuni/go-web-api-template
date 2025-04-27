package system

import (
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"time"
)

type Timer struct {
}

func NewTimer() system.ITimer {
	return &Timer{}
}

func (u Timer) Now() time.Time {
	return time.Now()
}
