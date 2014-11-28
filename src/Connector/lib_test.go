package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Lib(t *testing.T) {
	Convey("convertToByte", t, func() {
		Convey("Smaller than 256", func() {
			result := convertToByte(uint16(12))

			So(result[0], ShouldEqual, 0)
			So(result[1], ShouldEqual, 12)
		})

		Convey("Bigger than 256", func() {
			result := convertToByte(uint16(280))

			So(result[0], ShouldEqual, 1)
			So(result[1], ShouldEqual, 24)
		})

	})

	Convey("convertToUint16", t, func() {
		Convey("Smaller than 256", func() {
			result := convertToUint16([]byte{0, 12})

			So(result, ShouldEqual, 12)
		})

		Convey("Biger than 256", func() {
			result := convertToUint16([]byte{1, 24})

			So(result, ShouldEqual, 280)
		})
	})
}
