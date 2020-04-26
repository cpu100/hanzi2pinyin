package hanzi2pinyin

import (
    "bufio"
    "bytes"
    "encoding/hex"
    "fmt"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
    "io"
    "io/ioutil"
    "os"
    "sort"
    "strconv"
    "testing"
)

// 由 pinyin.txt 生成 pinyin.go
func TestPinyinGo(t *testing.T) {
    file, err := os.Open("pinyin.txt")
    if nil != err {
        panic(err)
    }
    defer file.Close()
    br := bufio.NewReader(file)
    var mapRunePinYin = make(map[rune][]byte)
    for {
        line, isPrefix, err := br.ReadLine()
        if nil != err {
            if io.EOF == err {
                break
            } else {
                panic(err)
            }
        }
        if isPrefix {
            panic("isPrefix")
        }
        sli := bytes.Split(line, []byte("=>"))
        r, err := strconv.ParseInt(string(sli[0]), 16, 32)
        if nil != err {
            panic(err)
        }
        if len(sli) != 2 {
            panic("Split")
        }
        // 为什么 []byte(unicode) 经传递之后就丢失 rune 信息，不能再还原回 utf-8 字符串？
        // todo 用 []byte 传递 unicode 字符串还安全吗？
        rs := []rune(string(bytes.TrimSpace(sli[1])))
        bs := make([]byte, len(rs))
        for i, r := range rs {
            if r > 'z' {
                bs[i] = mapTuneChar[r]
                if 0 == bs[i] {
                    panic(string(r) + " missing in mapTuneChar")
                }
            } else {
                bs[i] = byte(r)
            }
        }
        mapRunePinYin[rune(r)] = bs
    }
    fmt.Printf("0x%X %s %s \n", '龍', string('龍'), mapRunePinYin['龍'])
    fmt.Printf("0x%X %s %s \n", '馬', string('馬'), mapRunePinYin['馬'])
    fmt.Printf("0x%X %s %s \n", '精', string('精'), mapRunePinYin['精'])
    fmt.Printf("0x%X %s %s \n", '神', string('神'), mapRunePinYin['神'])
    file2, err := os.Create("pinyin.go")
    if nil != err {
        panic(err)
    }
    defer file2.Close()
    bw := bufio.NewWriter(file2)
    bw.WriteString("package hanzi2pinyin\n")
    bw.WriteString("var Map = map[rune]string{\n")
    var rs []rune
    for r := range mapRunePinYin {
        rs = append(rs, r)
    }
    // map是无序的
    // 为了让生成的 pinyin.go 内容有固定的顺序，便于 git 管理
    sort.Slice(rs, func(i, j int) bool {
        // 小的放前面
        return rs[i] < rs[j]
    })
    for _, r := range rs {
        fmt.Fprintf(bw, "'%s':\"%s\",\n", string(r), mapRunePinYin[r])
    }
    bw.WriteString("}\n")
    bw.Flush()
}

// https://github.com/cao-guang/pinyin/blob/master/pinyin.go
var mapTuneChar = map[rune]byte{
    'ā': 'a',
    'á': 'a',
    'ǎ': 'a',
    'à': 'a',
    'ō': 'o',
    'ó': 'o',
    'ǒ': 'o',
    'ò': 'o',
    'ē': 'e',
    'é': 'e',
    'ě': 'e',
    'è': 'e',
    'ī': 'i',
    'í': 'i',
    'ǐ': 'i',
    'ì': 'i',
    'ū': 'u',
    'ú': 'u',
    'ǔ': 'u',
    'ù': 'u',
    'ǖ': 'v',
    'ǘ': 'v',
    'ǚ': 'v',
    'ǜ': 'v',
    // 'ü': 'u',
    'ü': 'v',
    'ń': 'n',
    'ň': 'n',
    'ǹ': 'n',
    'ḿ': 'm',
}

// -7565=>è
// +7565=>lüè
// 纠正了 略 的拼音
// 不排除 pinyin.txt 还有其他错误

func TestUincode(t *testing.T) {
    // Unicode 确定了符号的编码
    // UTF-8 确定了编码的存储方式
    fmt.Printf("%s %X %X \n", string(0x4E25), []rune("严"), "严")
    // Rune 是一个符号的 Unicode 编码
    // UTF8 按它的规则存储之后，字节数值上不是直接的 Unicode 编码

    // UTF8的存储方式
    // https://www.cnblogs.com/tsingke/p/10853936.html

    // 0xxxxxxx
    // 110xxxxx 10xxxxxx
    // 1110xxxx 10xxxxxx 10xxxxxx
    // 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx

    // 对于单字节的符号，字节的第一位设为0，后面7位为这个符号的 Unicode 码。因此对于a-z0-9，UTF-8 编码和 ASCII 码是相同的;
    // 对于多字节的符号，第一个字节连续1后跟零表示字节数; x是有效存储位; 从第二个字节开始，前两位10是前缀;
}

// http://mengqi.info/html/2015/201507071345-using-golang-to-convert-text-between-gbk-and-utf-8.html
func Gb2312ToUtf8(s []byte) ([]byte) {
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewDecoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        panic(e)
    }
    return d
}

func Utf8ToGb2312(s []byte) ([]byte) {
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewEncoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        panic(e)
    }
    return d
}

func TestHZGB2312(t *testing.T) {
    // 不太理解这个 HZGB2312 和 GB2312 的关系
    fmt.Printf("%X\n", Utf8ToGb2312([]byte("连连看")))
    fmt.Println(string(Gb2312ToUtf8(Utf8ToGb2312([]byte("连连看")))))

    fmt.Printf("%X\n", Utf8ToGb2312([]byte("宝矿力水特")))
    fmt.Println(string(Gb2312ToUtf8(Utf8ToGb2312([]byte("宝矿力水特")))))

    for i := 0x347C; i <= 0x349C; i++ {
        // 这种编码总是带有 7E7B 前缀，不知道跟 GB2312 有什么关系
        bs, _ := hex.DecodeString(fmt.Sprintf("7E7B%X", i))
        fmt.Printf("%X %d %s \n", i, i, Gb2312ToUtf8(bs))
    }
}

func GbkToUtf8(s []byte) ([]byte) {
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        panic(e)
    }
    return d
}

func Utf8ToGbk(s []byte) ([]byte) {
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        panic(e)
    }
    return d
}

func TestGbk(t *testing.T) {
    // file,_:=os.Open("a.txt")
    // bs,_ := ioutil.ReadAll(file)
    // fmt.Printf("%X\n", bs) // B0A1

    // fmt.Printf("%X\n", Utf8ToGbk([]byte("啊")))
    // bs, _ := hex.DecodeString("FEFE")

    // 01-09区收录除汉字外的682个字符。
    // 10-15区为空白区，没有使用。
    // 16-55区收录3755个一级汉字，按拼音排序。
    // 56-87区收录3008个二级汉字，按部首/笔画排序。
    // 88-94区为空白区，没有使用。
    // https://www.qqxiuzi.cn/zh/hanzi-gb2312-bianma.php

    file2, err := os.Create("pinyin-gb2312.go")
    if nil != err {
        panic(err)
    }
    defer file2.Close()
    bw := bufio.NewWriter(file2)
    bw.WriteString("package hanzi2pinyin\n")
    bw.WriteString("var Map = map[rune]string{\n")

    fmt.Println("汉字个数", -5+(87-15)*94)

    for i := 16; i <= 87; i++ {
        for j := 1; j <= 94; j++ {
            bs := []byte{0xA0 + byte(i), 0xA0 + byte(j)}
            rs := []rune(string(GbkToUtf8(bs)))
            // fmt.Println(string(GbkToUtf8(bs)), Map[rs[0]])
            if py, ok := Map[rs[0]]; !ok {
                // 55区的最后5个是空的没错
                fmt.Println(string(rs), rs, bs[0]-0xA0, bs[1]-0xA0)
            } else {
                fmt.Fprintf(bw, "'%s':\"%s\",\n", string(rs), py)
            }
        }
    }

    bw.WriteString("}\n")
    bw.Flush()

}
