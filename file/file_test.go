package file

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFile(t *testing.T) {
	Convey("Given a JSON file", t, func() {
		file, err := Read("test_files/yet-another-configuration-file.yml")
		attributes := file.Attributes

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The extension should be yml", func() {
			So(file.extension, ShouldEqual, "yml")
		})

		Convey("The attributes should have length 6", func() {
			So(attributes, ShouldHaveLength, 6)
		})

		Convey("dependencies should have nested attributes", func() {
			So(attributes["dependencies"], ShouldNotBeEmpty)
		})

		Convey("Given a location with invalid fields", func() {
			wrong := []string{"options.config.version"}

			Convey("When new content is generated", func() {
				_, err := file.generateNewContent("1.0.0", wrong)

				Convey("err should NOT be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("Given a location with extra fields", func() {
			wrong := []string{"title.version"}

			Convey("When new content is generated", func() {
				_, err := file.generateNewContent("1.0.0", wrong)

				Convey("err should NOT be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("Given a correct location", func() {
			correct := []string{"info.extra.version"}

			Convey("When new content is generated", func() {
				m, err := file.generateNewContent("1.0.0", correct)

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The new content should include the new version", func() {
					fmt.Println(m)
					info := m["info"].(map[interface{}]interface{})
					extra := info["extra"].(map[interface{}]interface{})
					version := extra["version"]
					So(version, ShouldNotBeEmpty)
				})
			})
		})
	})
}
