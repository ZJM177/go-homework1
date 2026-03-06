package main
import ("fmt" ;"strings";"strconv";"sort")
func main() {
  // longestCommonPrefix([]string{"flower","flow","flight"})
	//plusOne([]int{1,2,3})
	merge([][]int{{1,3},{2,6},{8,10},{15,18}})
}
func merge(intervals [][]int) [][]int {
	if len(intervals)==0||len(intervals)==1 {return intervals }
    sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0] // 这里比较的是每个子切片的第一个元素
	})
    res:= [][]int{}
    for i:=0;i<len(intervals)-1;i++{
      if (intervals[i][1]>intervals[i+1][0]){
              z:= intervals[i+1][0]
              intervals[i+1][0]=intervals[i][1]
              intervals[i][1]=z
              res=append(res,intervals[i])
      }else if(intervals[i][1]==intervals[i+1][0]){
          res[i][0]=intervals[i][0]
          res[i][1]=intervals[i+1][1]
		  res=append(res,intervals[i])
      }else{
            res=append(res,intervals[i])
      }
    
   
}
 return res
}
func plusOne(digits []int) []int {
    var sumString string ;
    for _,v:= range digits{
        sumString+=strconv.FormatInt(int64(v),10)
    }
    
    s,_ := strconv.Atoi(sumString)
    s++
    ss:= strconv.Itoa(s)
    res:=make([]int,0) 
	fmt.Println("prefix=%v",res)
    for i:=0;i<len(ss);i++{
        s1,_ := strconv.Atoi(string(ss[i]))
        res=append(res,s1)
    }
    return res
}
func longestCommonPrefix(strs []string) string {
   length:= len(strs)
    if length==0||length==1 {return "" }
     prefix:=strs[0]
     for i:=1;i<length;i++{
        if !strings.HasPrefix(strs[i],prefix) {
            prefix= strings.TrimSuffix(prefix, string(prefix[len(prefix)-1]))
            strs[0]=prefix
			fmt.Println("prefix=%v,strs=%v",prefix,strs)
            prefix= longestCommonPrefix(strs)
			fmt.Println("prefixy=%v,",prefix)
        }
     }
     return prefix

}
