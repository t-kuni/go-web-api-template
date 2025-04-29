//go:generate go tool mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE

package system

type IUuidGenerator interface {
	Generate() (string, error)
}
