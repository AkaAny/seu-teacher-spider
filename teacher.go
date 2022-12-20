package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strings"
)

type Teacher struct {
	Name string
	Href string
}

func (x *Teacher) GetName() string {
	return x.Name
}

func (x *Teacher) GetInfo() *TeacherInfo {
	var teacherInfo = &TeacherInfo{}
	var pageURL = fmt.Sprintf("http://jszy.seu.edu.cn%s", x.Href)
	_, err := url.Parse(pageURL)
	if err != nil {
		teacherInfo.Name = x.Name
		return teacherInfo
	}
	var bodyData = make([]byte, 0)
	err = NewGoutClient().GET(pageURL).BindBody(&bodyData).Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bodyData))
	var doc = BuildGoQueryDocumentFromBodyData(bodyData)
	{
		var name = doc.Find(
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table.zlink > tbody > tr > td",
		).Text()
		name = strings.TrimSpace(name)
		fmt.Println(name)
		teacherInfo.Name = name
	}
	{
		var matches = parseSelectedMatchesWithRegexp(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table.llink > tbody > tr > td",
			regexp.MustCompile("(.+) (.+)"))
		if len(matches) == 0 { //这个老师没有填写信息
			return teacherInfo
		}
		title, department := matches[1], matches[2]
		teacherInfo.Title = title
		teacherInfo.Department = department
	}
	//brief
	{
		teacherInfo.PhoneNumber = parseSingleValue(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1)",
			"联系电话")
		teacherInfo.OfficeTime = parseSingleValue(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(2)",
			"办公时间")
		teacherInfo.Fax = parseSingleValue(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(2) > td:nth-child(1)",
			"传 真")
		teacherInfo.HomepageURL = parseSingleValue(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(2) > td:nth-child(2)",
			"主页网址")
		teacherInfo.OfficeLocation = parseSingleValue(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(3) > td:nth-child(1)",
			"办公地点")
		teacherInfo.Email = parseSingleValue(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(3) > td:nth-child(2)",
			"电子邮箱")
		teacherInfo.Address = parseSingleValue(doc,
			"#container_content > table > tbody > tr > td > table:nth-child(1) > tbody > tr > td:nth-child(4) > table.zlink > tbody > tr > td > div > div > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(4) > td:nth-child(1)",
			"通讯地址")
	}
	{
		var detailText = doc.Find(
			"#container_content > table > tbody > tr > td > table:nth-child(3) > tbody > tr > td:nth-child(2) > div",
		).Text()
		detailText = regexp.MustCompile("\n+ +").ReplaceAllString(detailText, "\n")
		detailText = regexp.MustCompile("\n+").ReplaceAllString(detailText, "\n")
		fmt.Println(detailText)
		teacherInfo.Detail = detailText
	}
	return teacherInfo
}

func parseSingleValue(doc *goquery.Document, selector string, fieldName string) string {
	var exprStr = fmt.Sprintf("%s：(.+)", fieldName)
	var expr = regexp.MustCompile(exprStr)
	var matches = parseSelectedMatchesWithRegexp(doc,
		selector,
		expr)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

func parseSelectedMatchesWithRegexp(doc *goquery.Document, selector string, expr *regexp.Regexp) []string {
	var rawStr = doc.Find(
		selector,
	).Text()
	rawStr = strings.TrimSpace(rawStr)
	var matches = expr.FindStringSubmatch(rawStr)
	return matches
}

type TeacherInfo struct {
	Name           string
	Title          string
	Department     string
	PhoneNumber    string
	OfficeTime     string
	Fax            string
	HomepageURL    string
	OfficeLocation string
	Email          string
	Address        string
	Detail         string
}
