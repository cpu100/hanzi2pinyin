package hanzi2pinyin

import "bytes"

// 汉字转拼音
// 其他除字母、数字、下划线之外全部用连接符(-)替换
func Identify(title []byte) []byte {
    var bss [][]byte
    var afterAlphanumeric bool
    for _, r := range []rune(string(title)) {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || ('_' == r) {
            if afterAlphanumeric {
                bss[len(bss)-1] = append(bss[len(bss)-1], byte(r))
            } else {
                bss = append(bss, []byte{byte(r)})
            }
            afterAlphanumeric = true
        } else {
            if py := Map[r]; len(py) > 0 {
                bss = append(bss, []byte(py))
            }
            afterAlphanumeric = false
        }
    }
    return bytes.Join(bss, []byte("-"))
}
