//go:generate go run cmd/gen_classids/generator.go cmd/gen_classids/main.go -classes ./classes.txt -output ./classids.go

package unity

type ClassID uint32
