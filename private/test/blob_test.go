package test

import (
	"context"
	"testing"
)

func TestBlobBin(t *testing.T) {
	ctx := context.Background()
	sampleData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header

	// 1. Leitura inicial — blob_bin deve ser NULL
	var blobBin []byte
	row := testDB.QueryRowContext(ctx, `SELECT blob_bin FROM TEST_TABLE WHERE id = 1`)
	if err := row.Scan(&blobBin); err != nil {
		t.Fatal("falha ao ler blob_bin:", err)
	}
	if blobBin != nil {
		t.Fatal("blob_bin deveria ser nil na leitura inicial")
	}
	t.Log("blob_bin é nil conforme esperado")

	// 2. UPDATE — gravar dados binários no blob_bin
	_, err := testDB.ExecContext(ctx, `UPDATE TEST_TABLE SET blob_bin = ? WHERE id = 1`, sampleData)
	if err != nil {
		t.Fatal("falha ao atualizar blob_bin:", err)
	}
	t.Logf("blob_bin atualizado com %d bytes", len(sampleData))

	// 3. Releitura — verificar se os dados foram gravados
	var blobBinAfter []byte
	row = testDB.QueryRowContext(ctx, `SELECT blob_bin FROM TEST_TABLE WHERE id = 1`)
	if err := row.Scan(&blobBinAfter); err != nil {
		t.Fatal("falha ao reler blob_bin após update:", err)
	}
	if blobBinAfter == nil {
		t.Fatal("blob_bin não deveria ser nil após update")
	}
	if len(blobBinAfter) != len(sampleData) {
		t.Fatalf("tamanho diferente: esperado %d, obtido %d", len(sampleData), len(blobBinAfter))
	}
	t.Logf("blob_bin relido com sucesso: %d bytes", len(blobBinAfter))
}

func TestBlobText(t *testing.T) {
	ctx := context.Background()
	sampleText := "Texto de teste para blob text — com acentuação: ção, à, é, ü"

	// 1. Leitura inicial — blob_text deve ser NULL
	var blobText *string
	row := testDB.QueryRowContext(ctx, `SELECT blob_text FROM TEST_TABLE WHERE id = 1`)
	if err := row.Scan(&blobText); err != nil {
		t.Fatal("falha ao ler blob_text:", err)
	}
	if blobText != nil {
		t.Fatal("blob_text deveria ser nil na leitura inicial")
	}
	t.Log("blob_text é nil conforme esperado")

	// 2. UPDATE — gravar texto no blob_text
	_, err := testDB.ExecContext(ctx, `UPDATE TEST_TABLE SET blob_text = ? WHERE id = 1`, sampleText)
	if err != nil {
		t.Fatal("falha ao atualizar blob_text:", err)
	}
	t.Logf("blob_text atualizado com %d caracteres", len(sampleText))

	// 3. Releitura — verificar se o texto foi gravado
	var blobTextAfter *string
	row = testDB.QueryRowContext(ctx, `SELECT blob_text FROM TEST_TABLE WHERE id = 1`)
	if err := row.Scan(&blobTextAfter); err != nil {
		t.Fatal("falha ao reler blob_text após update:", err)
	}
	if blobTextAfter == nil {
		t.Fatal("blob_text não deveria ser nil após update")
	}
	if *blobTextAfter != sampleText {
		t.Fatalf("texto diferente: esperado %q, obtido %q", sampleText, *blobTextAfter)
	}
	t.Logf("blob_text relido com sucesso: %d caracteres", len(*blobTextAfter))
}
