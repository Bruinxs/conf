package gconf

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestIniConfig_envH(t *testing.T) {
	ic := NewIniConfig()
	ic.Set("key1", "string value")
	ic.Set("key2", 10)
	ic.Set("key3", true)

	vs := []struct {
		Item string
		Want string
	}{
		{"$", "$"},
		{"${", "${"},
		{"${}", ""},
		{"${key1}", "string value"},
		{" ${key2}", "10"},
		{" ${ key3}", "true"},
		{"${key1 ,   key2 }", "string value10"},
		{"${key1,   key2,key3}", "string value10true"},
		{" ${  path}", os.Getenv("path")},
	}

	for _, v := range vs {
		rs := ic.envH(v.Item)
		if rs != v.Want {
			t.Errorf("item(%v) return(%v) not equal want(%v)", v.Item, rs, v.Want)
		}
	}
}

func TestIniConfig_parseCommand(t *testing.T) {
	ic := NewIniConfig()
	//illegal argument
	err := ic.parseCommand("fake command")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}
	err = ic.parseCommand("l()")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}
	err = ic.parseCommand("l(http://test.fake.com)")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}
	err = ic.parseCommand("l(https://www.baidu.com)")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}
	if len(ic.dict) != 0 {
		t.Errorf("dict len(%v) is not 0", len(ic.dict))
	}
	err = ic.parseCommand("l(file://test.fake.com)")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}
	err = ic.parseCommand("l(fake file path)")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}

	//http request config
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		str := `
			#comment
			key1=string
			;comment2
			[section1]
			key2=string1,string2,string3,string4,string5
			[section2]
			key3=10
			[section3]
			key4=9223372036854775807
			key5=true
			key6=99.99
		`
		w.Write([]byte(str))
	}))
	err = ic.parseCommand("l(" + srv.URL + ")")
	if err != nil {
		t.Error(err)
	}
	if len(ic.dict) != 6 {
		t.Errorf("dict(%v) len(%v) not equal 6", ic.dict, len(ic.dict))
	}

	key := "key1"
	if ic.String(key) != "string" {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, ic.String(key), "string")
	}

	key = "section1.key2"
	slice := ic.Strings(key)
	if len(slice) != 5 {
		t.Errorf("ic key(%v) return slice len(%v) not equal (%v)", key, len(slice), 5)
	}
	for i, s := range slice {
		want := fmt.Sprintf("string%v", i+1)
		if s != want {
			t.Errorf("slice index(%v) value(%v) not equal (%v)", i, s, want)
		}
	}

	key = "section2.key3"
	iv, err := ic.Int(key)
	if err != nil {
		t.Error(err)
	}
	if iv != 10 {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, iv, 10)
	}

	key = "section3.key4"
	iv64, err := ic.Int64(key)
	if err != nil {
		t.Error(err)
	}
	if iv64 != int64(9223372036854775807) {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, iv64, int64(9223372036854775807))
	}

	key = "section3.key5"
	bv, err := ic.Bool(key)
	if err != nil {
		t.Error(err)
	}
	if bv != true {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, bv, true)
	}

	key = "section3.key6"
	fv, err := ic.Float(key)
	if err != nil {
		t.Error(err)
	}
	if fv != 99.99 {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, fv, 99.99)
	}

	//load config file
	ic = NewIniConfig()
	filename := "tmp.conf"
	err = ioutil.WriteFile(filename, []byte(`
		#comment1
		key1 = 100
		;comment2
		[section]

		key2 = No
	`), os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
	err = ic.parseCommand("l(" + filename + ")")
	if err != nil {
		t.Error(err)
	}
	if len(ic.dict) != 2 {
		t.Errorf("ic dict len(%v) not equal 2", len(ic.dict))
	}
	key = "key1"
	iv, err = ic.Int(key)
	if err != nil {
		t.Error(err)
	}
	if iv != 100 {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, iv, 100)
	}
	key = "section.key2"
	bv, err = ic.Bool(key)
	if err != nil {
		t.Error(err)
	}
	if bv != false {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, bv, false)
	}

	err = os.Remove(filename)
	if err != nil {
		t.Error(err)
	}
}

func TestIniConfig_Parse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			[section]
			httpKey=100
			[section1]
			key1 = s1
		`))
	}))

	data := []string{
		"httpKey=20",
		"[section]",
		"httpKey=30",
		"@l(" + srv.URL + ")",
		"key2=s2",
	}

	filename := "tmp.conf"
	err := ioutil.WriteFile(filename, []byte(strings.Join(data, "\n")), os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			t.Error(err)
		}
	}()

	ic := NewIniConfig()
	_, err = ic.Parse(filename)
	if err != nil {
		t.Error(err)
	}
	if len(ic.dict) != 4 {
		t.Errorf("dict len(%v) not equal (%v)", len(ic.dict), 4)
	}

	key := "httpKey"
	if ic.String(key) != "20" {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, ic.String(key), "20")
	}
	key = "section.httpKey"
	if ic.String(key) != "100" {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, ic.String(key), "100")
	}
	key = "section.key2"
	if ic.String(key) != "s2" {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, ic.String(key), "s2")
	}
	key = "section1.key1"
	if ic.String(key) != "s1" {
		t.Errorf("ic key(%v) return val(%v) not equal (%v)", key, ic.String(key), "s1")
	}
}

func TestIniConfig_Set(t *testing.T) {
	ic := NewIniConfig()

	//string
	ic.Set("k1", nil)
	sv := ic.OptString("k1", "v1")
	if sv != "v1" {
		t.Errorf("key(%v) sv(%v) not equal (%v)", "k1", sv, "v1")
	}
	ic.Set("k1", "s1")
	sv = ic.OptString("k1", "v1")
	if sv != "s1" {
		t.Errorf("key(%v) sv(%v) not equal (%v)", "k1", sv, "s1")
	}

	//strings
	ss := ic.OptStrings("k2", []string{"1", "2", "3", "4"})
	if len(ss) != 4 {
		t.Errorf("key(%v) ss(%v) illegal", "k2", ss)
	}
	ic.Set("k2", "s1,s2,s3")
	ss = ic.OptStrings("k2", []string{"1", "2", "3", "4"})
	if len(ss) != 3 {
		t.Errorf("key(%v) ss(%v) illegal", "k2", ss)
	}

	//int and int64
	ic.Set("k3", "fake1")
	iv := ic.OptInt("k3", 10)
	if iv != 10 {
		t.Errorf("key(%v) iv(%v) not equal (%v)", "k3", iv, 10)
	}
	iv64 := ic.OptInt64("k3", 11)
	if iv64 != int64(11) {
		t.Errorf("key(%v) iv64(%v) not equal (%v)", "k3", iv64, 11)
	}
	ic.Set("k3", 123)
	iv = ic.OptInt("k3", 10)
	if iv != 123 {
		t.Errorf("key(%v) iv(%v) not equal (%v)", "k3", iv, 123)
	}
	iv64 = ic.OptInt64("k3", 11)
	if iv64 != int64(123) {
		t.Errorf("key(%v) iv64(%v) not equal (%v)", "k3", iv64, 123)
	}

	//float
	ic.Set("k4", "fake2")
	fv := ic.OptFloat("k4", 10.10)
	if fv != 10.10 {
		t.Errorf("key(%v) fv(%v) not equal (%v)", "k4", fv, 10.10)
	}
	ic.Set("k4", 11.11)
	fv = ic.OptFloat("k4", 10.10)
	if fv != 11.11 {
		t.Errorf("key(%v) fv(%v) not equal (%v)", "k4", fv, 11.11)
	}

	//bool
	bv := ic.OptBool("k5", true)
	if bv != true {
		t.Errorf("key(%v) bv(%v) not equal (%v)", "k5", bv, true)
	}
	ic.Set("k5", false)
	bv = ic.OptBool("k5", true)
	if bv != false {
		t.Errorf("key(%v) bv(%v) not equal (%v)", "k5", bv, false)
	}
	ic.Set("k5", "False")
	bv = ic.OptBool("k5", true)
	if bv != false {
		t.Errorf("key(%v) bv(%v) not equal (%v)", "k5", bv, false)
	}

	ic.Set("k5", "fake false")
	_, err := ic.Bool("k5")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}
	ic.Set("k5", 1)
	bv, err = ic.Bool("k5")
	if err != nil {
		t.Error(err)
	}
	if bv != true {
		t.Errorf("key(%v) bv(%v) not equal (%v)", "k5", bv, true)
	}
	ic.Set("k5", float64(0.0))
	bv, err = ic.Bool("k5")
	if err != nil {
		t.Error(err)
	}
	if bv != false {
		t.Errorf("key(%v) bv(%v) not equal (%v)", "k5", bv, false)
	}
	ic.Set("k5", 2)
	_, err = ic.Bool("k5")
	if err == nil {
		t.Errorf("err(%v) is nil", err)
	}
}
