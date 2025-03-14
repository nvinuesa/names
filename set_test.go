// Copyright 2013-2018 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package names_test

import (
	"github.com/juju/testing"
	gc "gopkg.in/check.v1"

	"github.com/juju/names/v6"
)

type tagSetSuite struct {
	testing.IsolationSuite

	foo  names.Tag
	bar  names.Tag
	baz  names.Tag
	bang names.Tag
}

var _ = gc.Suite(&tagSetSuite{})

func (s *tagSetSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)

	var err error

	s.foo, err = names.ParseTag("unit-wordpress-0")
	c.Assert(err, gc.IsNil)

	s.bar, err = names.ParseTag("unit-rabbitmq-server-0")
	c.Assert(err, gc.IsNil)

	s.baz, err = names.ParseTag("unit-mongodb-0")
	c.Assert(err, gc.IsNil)

	s.bang, err = names.ParseTag("machine-0")
	c.Assert(err, gc.IsNil)
}

func (tagSetSuite) TestEmpty(c *gc.C) {
	t := names.NewSet()
	c.Assert(t.Size(), gc.Equals, 0)
}

func (s tagSetSuite) TestInitialValues(c *gc.C) {
	t := names.NewSet(s.foo, s.bar)
	c.Assert(t.Size(), gc.Equals, 2)
}

func (tagSetSuite) TestInitialStringValues(c *gc.C) {
	t, err := names.NewSetFromStrings("unit-wordpress-0", "unit-rabbitmq-server-0")
	c.Assert(err, gc.IsNil)
	c.Assert(t.Size(), gc.Equals, 2)
}

func (tagSetSuite) TestInitialStringValuesBad(c *gc.C) {
	_, err := names.NewSetFromStrings("not-a-tag")
	c.Assert(err, gc.ErrorMatches, `"not-a-tag" is not a valid tag`)
}

func (tagSetSuite) TestSize(c *gc.C) {
	// Empty sets are empty.
	s := names.NewSet()
	c.Assert(s.Size(), gc.Equals, 0)

	s, err := names.NewSetFromStrings(
		"unit-wordpress-0",
		"unit-rabbitmq-server-0",
	)
	c.Assert(err, gc.IsNil)
	c.Assert(s.Size(), gc.Equals, 2)
}

func (tagSetSuite) TestSizeDuplicate(c *gc.C) {
	// Empty sets are empty.
	s := names.NewSet()
	c.Assert(s.Size(), gc.Equals, 0)

	// Size returns number of unique values.
	s, err := names.NewSetFromStrings(
		"unit-wordpress-0",
		"unit-rabbitmq-server-0",
		"unit-wordpress-0",
	)
	c.Assert(err, gc.IsNil)
	c.Assert(s.Size(), gc.Equals, 2)
}

func (s tagSetSuite) TestIsEmpty(c *gc.C) {
	// Empty sets are empty.
	t := names.NewSet()
	c.Assert(t.IsEmpty(), gc.Equals, true)

	// Non-empty sets are not empty.
	t = names.NewSet(s.foo)
	c.Assert(t.IsEmpty(), gc.Equals, false)

	// Newly empty sets work too.
	t.Remove(s.foo)
	c.Assert(t.IsEmpty(), gc.Equals, true)
}

func (s tagSetSuite) TestAdd(c *gc.C) {
	t := names.NewSet()
	t.Add(s.foo)
	c.Assert(t.Size(), gc.Equals, 1)
	c.Assert(t.Contains(s.foo), gc.Equals, true)
}

func (s tagSetSuite) TestAddDuplicate(c *gc.C) {
	t := names.NewSet()

	t.Add(s.foo)
	t.Add(s.bar)
	t.Add(s.bar)

	c.Assert(t.Size(), gc.Equals, 2)
}

func (s tagSetSuite) TestRemove(c *gc.C) {
	t := names.NewSet(s.foo, s.bar)
	t.Remove(s.foo)

	c.Assert(t.Contains(s.foo), gc.Equals, false)
	c.Assert(t.Contains(s.bar), gc.Equals, true)
}

func (s tagSetSuite) TestContains(c *gc.C) {
	t, err := names.NewSetFromStrings("unit-wordpress-0", "unit-rabbitmq-server-0")
	c.Assert(err, gc.IsNil)

	c.Assert(t.Contains(s.foo), gc.Equals, true)
	c.Assert(t.Contains(s.bar), gc.Equals, true)
	c.Assert(t.Contains(s.baz), gc.Equals, false)
}

func (s tagSetSuite) TestSortedValues(c *gc.C) {
	t := names.NewSet(s.foo, s.bang, s.baz, s.bar)
	values := t.SortedValues()

	c.Assert(values, gc.DeepEquals, []names.Tag{s.bang, s.baz, s.bar, s.foo})
}

func (s tagSetSuite) TestRemoveNonExistent(c *gc.C) {
	t := names.NewSet()
	t.Remove(s.foo)
	c.Assert(t.Size(), gc.Equals, 0)
}

func (s tagSetSuite) TestUnion(c *gc.C) {
	t1 := names.NewSet(s.foo, s.bar)
	t2 := names.NewSet(s.foo, s.baz, s.bang)
	union1 := t1.Union(t2)
	union2 := t2.Union(t1)

	c.Assert(union1.Size(), gc.Equals, 4)
	c.Assert(union2.Size(), gc.Equals, 4)

	c.Assert(union1, gc.DeepEquals, union2)
	c.Assert(union1, gc.DeepEquals, names.NewSet(s.foo, s.bar, s.baz, s.bang))
}

func (s tagSetSuite) TestIntersection(c *gc.C) {
	t1 := names.NewSet(s.foo, s.bar)
	t2 := names.NewSet(s.foo, s.baz, s.bang)

	int1 := t1.Intersection(t2)
	int2 := t2.Intersection(t1)

	c.Assert(int1.Size(), gc.Equals, 1)
	c.Assert(int2.Size(), gc.Equals, 1)

	c.Assert(int1, gc.DeepEquals, int2)
	c.Assert(int1, gc.DeepEquals, names.NewSet(s.foo))
}

func (s tagSetSuite) TestDifference(c *gc.C) {
	t1 := names.NewSet(s.foo, s.bar)
	t2 := names.NewSet(s.foo, s.baz, s.bang)

	diff1 := t1.Difference(t2)
	diff2 := t2.Difference(t1)

	c.Assert(diff1, gc.DeepEquals, names.NewSet(s.bar))
	c.Assert(diff2, gc.DeepEquals, names.NewSet(s.baz, s.bang))
}

func (s tagSetSuite) TestUninitializedPanics(c *gc.C) {
	f := func() {
		var t names.Set
		t.Add(s.foo)
	}
	c.Assert(f, gc.PanicMatches, "uninitalised set")
}
