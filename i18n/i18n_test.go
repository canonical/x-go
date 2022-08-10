// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2014-2022 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package i18n_test

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/check.v1"

	"github.com/canonical/golib/i18n"
)

func Test(t *testing.T) { check.TestingT(t) }

type i18nSuite struct{}

var _ = check.Suite(&i18nSuite{})

func (ts *i18nSuite) TestDefaults(c *check.C) {
	var call, callDefault reflect.Value

	// i18n.G
	call = reflect.ValueOf(i18n.G)
	callDefault = reflect.ValueOf(i18n.GDefault)
	c.Check(call.Pointer(), check.Equals, callDefault.Pointer(), check.Commentf("expected i18n.G == i18n.GDefault"))
	c.Check(i18n.G("Hello"), check.Equals, "Hello", check.Commentf("expected output unchanged"))

	// i18n.NG
	call = reflect.ValueOf(i18n.NG)
	callDefault = reflect.ValueOf(i18n.NGDefault)
	c.Check(call.Pointer(), check.Equals, callDefault.Pointer(), check.Commentf("expected i18n.NG == i18n.NGDefault"))
	c.Check(i18n.NG("Hello", "Hellos", 0), check.Equals, "Hellos", check.Commentf("expected plural form"))
	c.Check(i18n.NG("Hello", "Hellos", 1), check.Equals, "Hello", check.Commentf("expected singular form"))
	c.Check(i18n.NG("Hello", "Hellos", 2), check.Equals, "Hellos", check.Commentf("expected plural form"))
}

func (ts *i18nSuite) TestOverrides(c *check.C) {
	i18n.G = func(msgid string) string { return "something" }
	i18n.NG = func(msgid, msgid2 string, n int) string { return fmt.Sprintf("%s%d", "something", n) }

	// i18n.G
	c.Check(i18n.G("Hello"), check.Equals, "something", check.Commentf("expected translated output"))

	// i18n.NG
	c.Check(i18n.NG("Hello", "Hellos", 0), check.Equals, "something0", check.Commentf("expected translated output"))
}
