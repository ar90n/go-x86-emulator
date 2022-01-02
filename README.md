# go-x86-emulator
[自作エミュレータで学ぶx86アーキテクチャ](https://book.mynavi.jp/ec/products/detail/id=41347)のGoによる写経です

## 実行方法
```
$ go run *.go <path to bin file>
```

## メモ
サンプルプログラムをビルドするためにはMakefileを以下のように修正する必要がありました
```diff
<       -I$(Z_TOOLS)/i386-elf-gcc/include -g -fno-stack-protector -m32 -fno-pie
< LDFLAGS += --entry=start --oformat=binary -Ttext 0x7c00  -m elf_i386
---
>       -I$(Z_TOOLS)/i386-elf-gcc/include -g -fno-stack-protector
> LDFLAGS += --entry=start --oformat=binary -Ttext 0x7c00
```