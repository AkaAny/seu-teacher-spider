package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"regexp"
)

func GetOrgOverlapByRegexp(s string) string {
	var expr = regexp.MustCompile("orgOralp=(\\d+)")
	var orgOverlap = expr.FindStringSubmatch(s)[1]
	return orgOverlap
}

func BuildGoQueryDocumentFromBodyData(bodyData []byte) *goquery.Document {
	var bodyReader = bytes.NewReader(bodyData)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		panic(err)
	}
	return doc
}
