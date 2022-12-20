package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

type Department struct {
	Code string
	Name string
	Href string
}

func (x *Department) GetName() string {
	return x.Name
}

func (x *Department) GetTeachers() []*Teacher {
	//看第一页和页数
	//http://jszy.seu.edu.cn/_web/_plugs/teacherProfile/search/searchTeacher.do?
	//isFirst=false&_p=YXQ9MTQmZD01NCZwPTImZj0xJm09TiY_&orgOralp=91&searchType=1
	var firstPageDoc = x.getRaw(1)
	//body > div.framesplit
	var pageInfoStr = firstPageDoc.Find("body > div.framesplit").Text()
	var pageInfoStrExpr = regexp.MustCompile("页码：(\\d+)/(\\d+)")
	var matches = pageInfoStrExpr.FindStringSubmatch(pageInfoStr)
	var totalPageStr = matches[2]
	fmt.Println(totalPageStr)
	totalPage, err := strconv.ParseInt(totalPageStr, 10, 64)
	if err != nil {
		panic(err)
	}
	var teachers = make([]*Teacher, 0)
	for pageIndex := 1; pageIndex <= int(totalPage); pageIndex++ {
		var teachersOnPage = x.GetTeachersOnPage(pageIndex)
		teachers = append(teachers, teachersOnPage...)
	}
	return teachers
}

func (x *Department) getRaw(pageIndexFromOne int) *goquery.Document {
	const commonPrefix = "http://jszy.seu.edu.cn/_web/_plugs/teacherProfile/search/searchTeacher.do" +
		"?isFirst=false&_p=YXQ9MTQmZD01NCZwPTImZj0xJm09TiY_"
	var refererUrl = commonPrefix +
		fmt.Sprintf("&orgOralp=%s&searchType=1",
			x.Code)
	var bodyData = make([]byte, 0)
	err := NewGoutClient().POST(commonPrefix).SetHeader(map[string]interface{}{
		"Origin":  "http://jszy.seu.edu.cn",
		"Referer": refererUrl,
	}).SetWWWForm(map[string]interface{}{
		"keyword":    "",
		"pageIndex":  pageIndexFromOne,
		"orgOralp":   x.Code,
		"isFirst":    "",
		"searchType": "1",
	}).BindBody(&bodyData).Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bodyData))
	var doc = BuildGoQueryDocumentFromBodyData(bodyData)
	return doc
}

func (x *Department) GetTeachersOnPage(pageIndexFromOne int) []*Teacher {
	var doc = x.getRaw(pageIndexFromOne)
	var teachers = x.parseOnDocument(doc)
	return teachers
}

func (x *Department) parseOnDocument(doc *goquery.Document) []*Teacher {
	//body > div:nth-child(4) > div:nth-child(2) > div > div > ul > li:nth-child(1) > a
	var teachers = make([]*Teacher, 0)
	doc.Find("body > div:nth-child(4) > div:nth-child(2) > div > div > ul > li").
		Each(func(i int, liSelection *goquery.Selection) {
			var aSelection = liSelection.Find("a").First()
			var name = aSelection.Text()
			var href = aSelection.AttrOr("href", "")
			fmt.Println(name, href)
			var teacher = &Teacher{
				Name: name,
				Href: href,
			}
			teachers = append(teachers, teacher)
		})
	return teachers
}
