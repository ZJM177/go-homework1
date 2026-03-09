package homework01

import (
	"sort"
	"strconv"
	"strings"
)

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {
	mp := map[int]int{}
	for _, num := range nums {
		if mp[num] == 0 {
			mp[num] = 1
		} else {
			mp[num] = 2
		}
	}
	for k, v := range mp {
		if v == 1 {
			return k
		}
	}
	return 0
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {
	//负数和非0个位数不是回文数
	if x == 0 {
		return true
	}
	//负数和非0个位数不是回文数
	if x < 0 || (x <= 10 && x != 0) {
		return false
	}

	//将x转换为字符串
	s := strconv.Itoa(x)
	//双指针法，从两端开始比较
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	length := len(s) / 2
	for i := 0; i < length; i++ {
		s = strings.Replace(s, "{}", "", -1)
		s = strings.Replace(s, "[]", "", -1)
		s = strings.Replace(s, "()", "", -1)
	}
	return len(s) == 0
}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		for strings.Index(strs[i], prefix) != 0 {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix

}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	var sumString string
	for _, v := range digits {
		sumString += strconv.FormatInt(int64(v), 10)
	}

	s, _ := strconv.Atoi(sumString)
	s++
	ss := strconv.Itoa(s)
	res := make([]int, 0)
	for i := 0; i < len(ss); i++ {
		s1, _ := strconv.Atoi(string(ss[i]))
		res = append(res, s1)
	}
	return res
}

// 6. 删除有序数组中的重复项
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {
	if len(nums) < 2 {
		return len(nums)
	}
	var j int = 0
	for i := 1; i < len(nums); i++ {
		if nums[j] != nums[i] {
			j += 1
			nums[j] = nums[i]
		}

	}
	return j + 1
}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func Merge(intervals [][]int) [][]int {
	if len(intervals) == 0 || len(intervals) == 1 {
		return intervals
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0] // 这里比较的是每个子切片的第一个元素
	})
	res := [][]int{}
	for i := 0; i < len(intervals)-1; i++ {
		if intervals[i][1] > intervals[i+1][0] {
			z := intervals[i+1][0]
			intervals[i+1][0] = intervals[i][1]
			intervals[i][1] = z
			res = append(res, intervals[i])
			if i == len(intervals)-2 {
				res = append(res, intervals[len(intervals)-1])
			}
		} else if intervals[i][1] == intervals[i+1][0] {
			res[i][0] = intervals[i][0]
			res[i][1] = intervals[i+1][1]
			res = append(res, intervals[i])
		} else {
			res = append(res, intervals[i])
			if i == len(intervals)-2 {
				res = append(res, intervals[len(intervals)-1])
			}
		}
	}

	return res
}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	rsult := make([]int, 2)
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				rsult[0] = i
				rsult[1] = j

				return rsult
			}
		}
	}
	return nil
}
