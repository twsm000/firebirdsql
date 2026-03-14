# Correção: Escrita de BLOB no Firebird 2.0.1

## O Problema

Ao tentar fazer UPDATE em colunas BLOB (binário ou texto) no Firebird 2.0.1, o servidor retorna:

```
feature is not supported
BLOB and array data types are not supported for move operation
```

A leitura funciona normalmente — apenas a escrita falha.

## Causa Raiz

O driver serializa os parâmetros na função `paramsToBlr()` em `wireprotocol.go`. Quando recebe um `[]byte` (dados binários), ele decide como enviar com base no **tamanho**:

- **< 32.767 bytes**: envia inline como **VARCHAR** (BLR tipo 14)
- **≥ 32.767 bytes**: cria um objeto **BLOB** via `createBlob()` (BLR tipo 9)

O Firebird ≥ 2.5 aceita a conversão implícita VARCHAR → BLOB, mas o **Firebird 2.0 não suporta** essa operação de "move" entre tipos.

Como nossos dados de teste são pequenos (8 bytes de header PNG), o driver envia como VARCHAR, e o servidor rejeita.

## A Correção

Forçar que `[]byte` **sempre** seja enviado como BLOB, independente do tamanho. Isso é semanticamente correto — `[]byte` em Go representa dados binários, que correspondem a BLOB no Firebird.

A mudança fica no `case []byte:` dentro de `paramsToBlr()`:

```diff
 case []byte:
-    if len(f) < MAX_CHAR_LENGTH {
-        blr, v = _bytesToBlr(f)
-    } else {
-        v, _ = p.createBlob(f, transHandle)
-        blr = []byte{9, 0}
-    }
+    v, _ = p.createBlob(f, transHandle)
+    blr = []byte{9, 0}
```

Isso garante compatibilidade com **todas as versões** do Firebird, sem quebrar nada nas versões mais novas.
