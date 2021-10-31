package store

import (
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStoreMem(t *testing.T) {
	store := NewMemStore()
	Convey("Test Store Mem", t, func() {
		key := uuid.New().String()
		value := []byte(key)

		err := store.Set(key, value)
		So(err, ShouldBeNil)

		v, err := store.Get(key)
		So(err, ShouldBeNil)
		So(v, ShouldResemble, value)

		err = store.Remove(key)
		So(err, ShouldBeNil)

		v, err = store.Get(key)
		So(err, ShouldBeError)
		So(v, ShouldBeNil)
	})
}
