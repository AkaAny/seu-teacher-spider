package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type School struct {
	Code string
	Name string
	Href string
}

func (x *School) GetName() string {
	return x.Name
}

func (x *School) GetDepartments() []*Department {
	const departmentListUrl = "http://jszy.seu.edu.cn/_web/_plugs/teacherProfile/search/searchTeacher.do?" +
		"isFirst=false&_p=YXQ9MTQmZD01NCZwPTImZj0xJm09TiY_"
	var bodyData = make([]byte, 0)
	err := NewGoutClient().POST(departmentListUrl).SetHeader(map[string]interface{}{
		"Origin":  "http://jszy.seu.edu.cn",
		"Referer": "http://jszy.seu.edu.cn/_web/_plugs/teacherProfile/search/searchIndex.jsp?_p=YXM9MiZ0PTE0JmQ9NTQmcD0yJmY9MSZtPU4m",
	}).SetWWWForm(map[string]interface{}{
		"keyword":    "",
		"orgOralp":   x.Code,
		"searchType": "1",
	}).BindBody(&bodyData).Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bodyData))
	//body > div:nth-child(4) > div > div > div.t2 > ul > li:nth-child(3) > a
	var doc = BuildGoQueryDocumentFromBodyData(bodyData)
	var departments = make([]*Department, 0)
	doc.Find("body > div:nth-child(4) > div > div > div.t2 > ul > li").
		Each(func(i int, liSelection *goquery.Selection) {
			var aSelection = liSelection.Find("a").First()
			var name = aSelection.Text()
			var href = aSelection.AttrOr("href", "")
			fmt.Println(name, href)
			var code = GetOrgOverlapByRegexp(href)
			var department = &Department{
				Code: code,
				Name: name,
				Href: href,
			}
			departments = append(departments, department)
		})
	return departments
}
