# bmp-parse

BMP 位图解析：读出头信息（宽、高、位深、调色板条目数）并算出像素校验和。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/bmpcheck --sample=padding
go run ./cmd/bmpcheck --sample=topdown
go run ./cmd/bmpcheck --sample=offset
go run ./cmd/bmpcheck --sample=work
```

## 口径

- **行步长**：每行字节数按 4 字节对齐（不足补零），填充字节参与逐字节计数。
- **朝向**：高为正表示自下而上存储（最后一行是最上面一行）；高为负表示自上而下。
- **像素起点**：用头里第 10 字节起的偏移字段，不假定固定值。
- **8 位图**：像素字节是调色板下标，按调色板里的 B+G+R 求和。
- **校验和**：按逻辑行序、逐字节累加 `字节值 × 序号`（序号从 1 起）。

## 不变量

- 解析出的宽高与头字段一致，高的符号保留。
- 校验和只由像素数据与朝向决定，与解析次数无关。
- `scanned`（调色板条目检查次数）不随像素数乘调色板大小放大：两千像素的 8 位图不超过 9000。

## 输出契约

打印一行 JSON：`width`、`height`、`bits`、`palette`、`checksum`、`scanned`。
`--sample=work` 只打印 `{"width": N, "scanned": N}`。
