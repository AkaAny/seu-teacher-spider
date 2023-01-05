package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"path/filepath"
)

func main() {
	const orgListPageURL = "http://jszy.seu.edu.cn/_web/_plugs/teacherProfile/search/fetchOrg.do?" +
		"oid=5&_p=YXM9MiZ0PTE0JmQ9NTUmcD0xJm09TiY_"
	var bodyData = make([]byte, 0)
	err := NewGoutClient().SetHeader(map[string]interface{}{
		"Referer": "http://jszy.seu.edu.cn/",
	}).GET(orgListPageURL).BindBody(&bodyData).Do()
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(bodyData))
	//body > table > tbody > tr:nth-child(1) > td > a:nth-child(3)
	var doc = BuildGoQueryDocumentFromBodyData(bodyData)
	var schools = make([]*School, 0)
	doc.Find("body > table > tbody > tr:nth-child(1) > td > a").
		Each(func(i int, aSelection *goquery.Selection) {
			var name = aSelection.Text()
			var href = aSelection.AttrOr("href", "")
			var code = GetOrgOverlapByRegexp(href)
			var school = &School{
				Code: code,
				Name: name,
				Href: href,
			}
			schools = append(schools, school)
		})
	rootPath, err := filepath.Abs("data")
	if err != nil {
		panic(err)
	}
	for _, school := range schools {
		basePath := CreateDirectory(school, rootPath)
		var departments = school.GetDepartments()
		for _, department := range departments {
			basePath := CreateDirectory(department, basePath)
			var teachers = department.GetTeachers()
			for _, teacher := range teachers {
				var teacherInfo = teacher.GetInfo()
				var savePath = filepath.Join(basePath, teacher.Name+".json")
				fmt.Println(savePath)
				rawData, err := json.MarshalIndent(teacherInfo, "", "    ")
				if err != nil {
					panic(err)
				}
				err = os.WriteFile(savePath, rawData, 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
