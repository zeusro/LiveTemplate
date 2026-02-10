// 包 roman_converter 提供了将罗马数字与整数互相转换的功能
package roman_converter

import "strings"

// romanToInt 将罗马数字字符串转换为整数
// 参数 s: 表示罗马数字的字符串
// 返回值: 转换后的整数值
func romanToInt(s string) int {
	romanMap := map[rune]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	n := len(s)
	total := 0
	for i := 0; i < n; i++ {
		value := romanMap[rune(s[i])]
		// 如果当前值小于下一个值，则减去当前值
		if i+1 < n && value < romanMap[rune(s[i+1])] {
			total -= value
		} else {
			// 否则加上当前值
			total += value
		}
	}
	return total
}

// intToRoman 将整数转换为罗马数字字符串
// 参数 num: 表示需要转换的整数
// 返回值: 转换后的罗马数字字符串
func intToRoman(num int) string {
	val := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	var result strings.Builder
	for i := 0; i < len(val); i++ {
		// 按照从大到小的顺序匹配整数值
		for num >= val[i] {
			num -= val[i]
			result.WriteString(symbols[i])
		}
	}
	return result.String()
}
