package test

import (
	"database/sql/driver"
	"fmt"
)

// Blob é um tipo customizado para BLOB BIN do Firebird.
// Implementa as interfaces sql.Scanner e driver.Valuer para
// compatibilidade com o driver de banco de dados.
type Blob []byte

// Scan implementa sql.Scanner para leitura de BLOB BIN do banco.
func (b *Blob) Scan(value any) error {
	if value == nil {
		*b = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*b = make(Blob, len(v))
		copy(*b, v)
		return nil
	case string:
		*b = Blob(v)
		return nil
	}

	return fmt.Errorf("Blob.Scan: tipo não suportado: %T", value)
}

// Value implementa driver.Valuer para escrita de BLOB BIN no banco.
func (b Blob) Value() (driver.Value, error) {
	// len(b) == 0 cobre tanto nil quanto slices vazios,
	// evitando que o driver envie '' para o Firebird
	if len(b) == 0 {
		return nil, nil
	}
	return []byte(b), nil
}

// Bytes retorna o conteúdo como []byte.
func (b Blob) Bytes() []byte {
	return []byte(b)
}

// BlobText é um tipo customizado para BLOB SUB_TYPE 1 (texto) do Firebird.
// Implementa as interfaces sql.Scanner e driver.Valuer para
// compatibilidade com o driver de banco de dados.
type BlobText string

// Scan implementa sql.Scanner para leitura de BLOB TEXT do banco.
func (bt *BlobText) Scan(value any) error {
	if value == nil {
		*bt = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*bt = BlobText(v)
		return nil
	case string:
		*bt = BlobText(v)
		return nil
	}

	return fmt.Errorf("BlobText.Scan: tipo não suportado: %T", value)
}

// Value implementa driver.Valuer para escrita de BLOB TEXT no banco.
func (bt BlobText) Value() (driver.Value, error) {
	if len(bt) == 0 {
		return nil, nil
	}
	return []byte(bt), nil
}

// String retorna o conteúdo como string.
func (bt BlobText) String() string {
	return string(bt)
}
