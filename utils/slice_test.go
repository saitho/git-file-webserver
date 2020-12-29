package utils_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/saitho/static-git-file-server/utils"
)

func TestSplitModuleName(t *testing.T) {
	Convey("test array contains", t, func() {
		array := []string{"a", "b"}
		So(utils.Contains(array, "a"), ShouldBeTrue)
		So(utils.Contains(array, "b"), ShouldBeTrue)
		So(utils.Contains(array, "c"), ShouldBeFalse)
	})
	Convey("test unpack", t, func() {
		array := []string{"a", "b", "c", "d"}
		var var1, var2, var3 string
		utils.Unpack(array, &var1, &var2, nil, &var3)
		So(var1, ShouldEqual, "a")
		So(var2, ShouldEqual, "b")
		So(var3, ShouldEqual, "d")
	})
}
