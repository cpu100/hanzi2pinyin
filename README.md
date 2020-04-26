一个字典
------------
```go
var Map = map[rune]string{
    '汉': "han",
    '语': "yu",
    '拼': "pin",
    '音': "yin",
    '..': "...",
}
```


### usage

```
go get github.com/cpu100/hanzi2pinyin@gbk
或者
go get github.com/cpu100/hanzi2pinyin@gb2312 // 仅包含 GB2312 共 6763 个字符，占更少内存
```

```go
import "github.com/cpu100/hanzi2pinyin"

fmt.Println(hanzi2pinyin.Map['漢']) // han
fmt.Println(hanzi2pinyin.Map['語']) // yu
fmt.Println(hanzi2pinyin.Map['拼']) // pin
fmt.Println(hanzi2pinyin.Map['音']) // yin
```


### 链接
https://github.com/cao-guang/pinyin  
https://github.com/chain-zhang/pinyin
