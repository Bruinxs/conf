package conf

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIniConfig(t *testing.T) {
	Convey("Tese ini config", t, func() {
		ini := NewIniConfig()
		_, err := ini.ParseFile("./testdata/conf.ini")
		So(err, ShouldBeNil)

		Convey("Get value", func() {
			So(ini.Exists("string_val"), ShouldBeTrue)

			strV := ini.String("string_val")
			So(strV, ShouldEqual, "bar_string_foo_string")

			strsV := ini.Strings("strings_val")
			So(strsV, ShouldResemble, []string{"s1", "s2", "s3"})

			intV, err := ini.Int("int_val")
			So(err, ShouldBeNil)
			So(intV, ShouldEqual, 15)

			floatV, err := ini.Float("float_val")
			So(err, ShouldBeNil)
			So(floatV, ShouldEqual, 3.14)

			boolV, err := ini.Bool("bool_val")
			So(err, ShouldBeNil)
			So(boolV, ShouldBeTrue)
		})

		Convey("Get defalut value", func() {
			So(ini.Exists("string_def_val"), ShouldBeFalse)

			strV := ini.DefaultString("string_def_val", "def")
			So(strV, ShouldEqual, "def")

			strsV := ini.DefaultStrings("strings_def_val", []string{"s"})
			So(strsV, ShouldResemble, []string{"s"})

			intV := ini.DefaultInt("int_def_val", 1)
			So(intV, ShouldEqual, 1)

			intV = ini.DefaultInt("string_val", 2)
			So(intV, ShouldEqual, 2)

			floatV := ini.DefaultFloat("float_def_val", 1.1)
			So(floatV, ShouldEqual, 1.1)

			floatV = ini.DefaultFloat("string_val", 1.2)
			So(floatV, ShouldEqual, 1.2)

			boolV := ini.DefaultBool("bool_def_val", false)
			So(boolV, ShouldBeFalse)

			boolV = ini.DefaultBool("string_val", true)
			So(boolV, ShouldBeTrue)
		})
	})
}
