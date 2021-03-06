//             ,%%%%%%%%,
//           ,%%/\%%%%/\%%
//          ,%%%\c "" J/%%%
// %.       %%%%/ o  o \%%%
// `%%.     %%%%    _  |%%%
//  `%%     `%%%%(__Y__)%%'
//  //       ;%%%%`\-/%%%'
// ((       /  `%%%%%%%'
//  \\    .'          |
//   \\  /       \  | |
//    \\/攻城狮保佑) | |
//     \         /_ | |__
//     (___________)))))))                   `\/'
/*
 * 修订记录:
 * long.qian 2021-10-05 12:19 创建
 */

/**
 * @author long.qian
 */

package go_lazy_util

import (
	"regexp"
	"strings"
)

var (
	StringUtil = new(stringUtil)
)

type stringUtil struct {
}

func (receiver *stringUtil) IsEmpty(str string) bool {
	if str == "" {
		return true
	}
	return false
}

func (receiver *stringUtil) IsNotEmpty(str string) bool {
	return !receiver.IsEmpty(str)
}

func (receiver *stringUtil) IsBlank(str string) bool {
	if receiver.IsEmpty(str) {
		return true
	}
	if len(strings.TrimSpace(str)) == 0 {
		return true
	}
	return false
}

func (receiver *stringUtil) IsNotBlank(str string) bool {
	return !receiver.IsBlank(str)
}

// IsMatch 是否匹配指定的正则表达式
func (receiver *stringUtil) IsMatch(str string, pattern string) bool {
	m, err := regexp.Match(pattern, []byte(str))
	if err != nil {
		return false
	}
	return m
}

// IsIp 是否是 IP 地址
func (receiver *stringUtil) IsIp(str string) bool {
	return receiver.IsMatch(str, `((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))`)
}

// MatchingRate 基于编辑距离算法，计算两个字符串的相似度（0 到 1 之间），越靠近 1 代表越相似
func (receiver *stringUtil) MatchingRate(str1, str2 string, ignoreCase bool) float64 {
	if str1 == str2 {
		return 1
	}
	distance := receiver.EditDistanceDP(str1, str2, ignoreCase)
	if distance == 0 {
		return 1
	}
	totalDistance := float64(receiver.getMaxEditDistance(str1, str2, ignoreCase))
	return float64(1) - (float64(distance) / float64(totalDistance))
}

// 获取最大距离
func (receiver *stringUtil) getMaxEditDistance(str1, str2 string, ignoreCase bool) int {
	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}
	var totalDistance int
	if len(str1) == len(str2) || len(str1) > len(str2) {
		totalDistance = receiver.EditDistanceDP(str1, "", ignoreCase)
	} else {
		totalDistance = receiver.EditDistanceDP(str2, "", ignoreCase)
	}
	return totalDistance
}

// EditDistanceDP 编辑距离计算
func (receiver *stringUtil) EditDistanceDP(str1, str2 string, ignoreCase bool) int {
	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}
	d := make([][]int, len(str1)+1)
	for i := range d {
		d[i] = make([]int, len(str2)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(str2); j++ {
		for i := 1; i <= len(str1); i++ {
			if str1[i-1] == str2[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(str1)][len(str2)]
}
