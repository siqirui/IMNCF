package Mregexp

import (
	"regexp"
)
type regexpHandle interface{
	FindPhone(string)[][]string
	FindeMail(string)[]string
}
type Regexp struct {
	rePhone string
	reMail string
	reLink string
	reIDNumber string
	reImg string
}
func NewMregexp() *Regexp{
	mregexp := Regexp{
		`(1[3-9][0-9])([0-9]{4})([0-9]{4})`,
		`\w+@\w+\.[a-zA-Z]{2,3}(\.[a-zA-Z]{2,3})?`,
		`<a[\s\S]+?href="(http[\s\S]+?)"`,
		`[1-9]\d{5}((19\d{2})|(20((0\d))|(1[0-9])))((0[1-9])|(1[012]))((0[1-9])|([12]\d)|(3[01]))\d{3}[\dX|x]`,
		`<img[\s\S]+?src="(http[\s\S]+?)"`,
	}
	return &mregexp
}
func (mregexp * Regexp) FindPhone(html string) [][]string {
	re := regexp.MustCompile(mregexp.rePhone)
	allstring := re.FindAllStringSubmatch(html,-1)
	return allstring
	//fmt.Println(allstring)
}
func (mregexp * Regexp) FindImg(html string)[][]string{
	re := regexp.MustCompile(mregexp.reImg)
	allstring := re.FindAllStringSubmatch(html,-1)
	return allstring
}

func (mregexp * Regexp) FindMail(html string)[]string{
	re := regexp.MustCompile(mregexp.reMail)
	allstring := re.FindAllString(html,-1)
	return allstring
}
func(mregexp * Regexp) FindLink(html string)[][]string{
	re := regexp.MustCompile(mregexp.reLink)
	allstring := re.FindAllStringSubmatch(html,-1)
	return allstring
}
func (mregexp * Regexp) FindIDMember(html string)[][]string{
	re := regexp.MustCompile(mregexp.reIDNumber)
	allstring := re.FindAllStringSubmatch(html,-1)
	return allstring
}

