package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlobBin(t *testing.T) {
	ctx := context.Background()
	sampleData := Blob{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header

	// 1. Leitura inicial — blob_bin deve ser NULL
	var blobBin Blob
	row := testDB.QueryRowContext(ctx, `SELECT blob_bin FROM TEST_TABLE WHERE id = 1`)
	require.NoError(t, row.Scan(&blobBin))
	assert.Nil(t, blobBin, "blob_bin deveria ser nil na leitura inicial")

	// 2. UPDATE — gravar dados binários no blob_bin
	_, err := testDB.ExecContext(ctx, `UPDATE TEST_TABLE SET blob_bin = ? WHERE id = 1`, sampleData)
	require.NoError(t, err, "falha ao atualizar blob_bin")
	t.Logf("blob_bin atualizado com %d bytes", len(sampleData))

	// 3. Releitura — verificar se os dados foram gravados
	var blobBinAfter Blob
	row = testDB.QueryRowContext(ctx, `SELECT blob_bin FROM TEST_TABLE WHERE id = 1`)
	require.NoError(t, row.Scan(&blobBinAfter))
	require.NotNil(t, blobBinAfter, "blob_bin não deveria ser nil após update")
	assert.Equal(t, sampleData, blobBinAfter)
	t.Logf("blob_bin relido com sucesso: %d bytes", len(blobBinAfter))
}

func TestBlobText(t *testing.T) {
	ctx := context.Background()
	sampleText := BlobText("Texto de teste para blob text — com acentuação: ção, à, é, ü")

	// 1. Leitura inicial — blob_text deve ser NULL
	var blobText BlobText
	row := testDB.QueryRowContext(ctx, `SELECT blob_text FROM TEST_TABLE WHERE id = 1`)
	require.NoError(t, row.Scan(&blobText))
	assert.Empty(t, blobText, "blob_text deveria ser vazio na leitura inicial")

	// 2. UPDATE — gravar texto no blob_text
	_, err := testDB.ExecContext(ctx, `UPDATE TEST_TABLE SET blob_text = ? WHERE id = 1`, sampleText)
	require.NoError(t, err, "falha ao atualizar blob_text")
	t.Logf("blob_text atualizado com %d caracteres", len(sampleText))

	// 3. Releitura — verificar se o texto foi gravado
	var blobTextAfter BlobText
	row = testDB.QueryRowContext(ctx, `SELECT blob_text FROM TEST_TABLE WHERE id = 1`)
	require.NoError(t, row.Scan(&blobTextAfter))
	require.NotEmpty(t, blobTextAfter, "blob_text não deveria ser vazio após update")
	assert.Equal(t, sampleText, blobTextAfter)
	t.Logf("blob_text relido com sucesso: %d caracteres", len(blobTextAfter))
}
