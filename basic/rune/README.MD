##go rune
Go 语言的字符有以下两种： 
1. uint8 类型，或者叫 byte 型，代表了 ASCII 码的一个字符。 
2. rune 类型，代表一个 UTF-8 字符。 
当需要处理中文、日文或者其他复合字符时，则需要用到 rune 类型。
rune 类型实际是一个 int32。 Go 使用了特殊的 rune 类型来处理 Unicode，
让基于 Unicode 的文本处理更为方便，也可 以使用 byte 型进行默认字符串处理，性能和扩展性都有照顾。
