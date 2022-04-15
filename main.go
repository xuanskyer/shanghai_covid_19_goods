package main

import (
	"fmt"
	"github.com/yanyiwu/gojieba"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	ReplaceWords = map[string]string{
		" ":  "",
		"，":  "",
		",":  "",
		"。":  "",
		".":  "",
		"、":  "",
		"\\": "",
		"/":  "",
		"（":  "",
		"）":  "",
		"!":  "",
		"！":  "",
		"？":  "",
		"?":  "",
		"；":  "",
		";":  "",
		"：":  "",
		":":  "",
		"“":  "",
		"”":  "",
		"‘":  "",
		"’":  "",
		"《":  "",
		"》":  "",
		"\"": "",
		"\n": "",
		"\r": "",
		"\t": "",
		"1":  "",
		"2":  "",
		"3":  "",
		"4":  "",
		"5":  "",
		"6":  "",
		"7":  "",
		"8":  "",
		"9":  "",
		"0":  "",
	}
	words = []string{}
	foodsFile = "./data/foods.txt"
	wordStatFile = "./data/foods_split.txt"
)

//一个jieba分词的demo
func main() {
	fp, err := os.Open(foodsFile) // 获取文件指针
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()
	fileInfo, err := fp.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = fp.Read(buffer) // 读取文件内容
	if err != nil {
		fmt.Println(err)
		return
	}
	foods := string(buffer)
	foods = stringReplaces(foods, ReplaceWords)
	fmt.Println("过滤非法字符之后的foods: ", foods)
	x := gojieba.NewJieba() // 分词器
	defer x.Free()

	words = x.CutAll(foods)
	wordCount := make(map[string]int, 0)
	for _, word := range words {
		wordCount[word]++
	}

	wordStat := sortSlice{}
	for w, c := range wordCount {
		wordStat = append(wordStat, Worder{Word: w, Count: c})
	}
	fw, err := os.OpenFile(wordStatFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err.Error())
	}
	defer fw.Close()
	sort.Sort(wordStat)

	fmt.Println("jieba分词 && 排序 后 TOP20 结果: ")
	for index, w := range wordStat {
		if index < 20 {
			fmt.Println(w.Word, w.Count)
		}
		fw.WriteString(w.Word + " " + strconv.Itoa(w.Count) + "\n")
	}
	fmt.Println("统计结果写文件：", wordStatFile)
	if err != nil {
		log.Println(err.Error())
	}
}

type Worder struct {
	Word  string
	Count int
}
type sortSlice []Worder

func (l sortSlice) Len() int           { return len(l) }
func (l sortSlice) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l sortSlice) Less(i, j int) bool { return l[i].Count > l[j].Count }

func stringReplaces(s string, replaces map[string]string) string {
	for k, v := range replaces {
		s = strings.ReplaceAll(s, k, v)
	}
	return s

}
