// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2020 Canonical Ltd
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

package randutil_test

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"

	"github.com/canonical/x-go/randutil"
)

type cryptoRandutilSuite struct{}

var _ = Suite(&cryptoRandutilSuite{})

func (s *cryptoRandutilSuite) TestCryptoTokenBytes(c *C) {
	x, err := randutil.CryptoTokenBytes(5)
	c.Assert(err, IsNil)
	c.Check(x, HasLen, 5)
}

func (s *cryptoRandutilSuite) TestCryptoToken(c *C) {
	x, err := randutil.CryptoToken(5)
	c.Assert(err, IsNil)

	b, err := base64.RawURLEncoding.DecodeString(x)
	c.Assert(err, IsNil)
	c.Check(b, HasLen, 5)
}

var (
	KernelTestUUID        = "1031319a-b661-4c01-aafa-6def8a118944"
	KernelUUIDPathDefault = randutil.KernelUUIDPath
)

func (s *cryptoRandutilSuite) TestRandomKernelUUIDNoFile(c *C) {
	tmpDir := c.MkDir()
	uuidPath := filepath.Join(tmpDir, "no-file")

	randutil.KernelUUIDPath = uuidPath
	defer func() {
		randutil.KernelUUIDPath = KernelUUIDPathDefault
	}()

	value, err := randutil.RandomKernelUUID()
	c.Check(value, Equals, "")
	c.Check(err, ErrorMatches, "*.no such file or directory")
}

func (s *cryptoRandutilSuite) TestRandomKernelUUIDNoPerm(c *C) {
	tmpDir := c.MkDir()
	uuidPath := filepath.Join(tmpDir, "no-perm")
	err := ioutil.WriteFile(uuidPath, []byte(KernelTestUUID), 0)
	if err != nil {
		c.Errorf("Write failed with: %s", err.Error())
	}

	randutil.KernelUUIDPath = uuidPath
	defer func() {
		randutil.KernelUUIDPath = KernelUUIDPathDefault
	}()

	value, err := randutil.RandomKernelUUID()
	c.Check(value, Equals, "")
	c.Assert(err, ErrorMatches, "*.permission denied")
}

func (s *cryptoRandutilSuite) TestRandomKernelUUID(c *C) {
	for _, uuid := range []string{
		KernelTestUUID,
		" \t\n " + KernelTestUUID + " \n\t\r\n",
	} {
		// Create new path on each iteration because we cannot
		// reuse previous path to read-only (0444) file.
		tmpDir := c.MkDir()
		uuidPath := filepath.Join(tmpDir, "uuid")
		err := ioutil.WriteFile(uuidPath, []byte(uuid), 0444)
		if err != nil {
			c.Errorf("Write failed with: %s", err.Error())
		}

		randutil.KernelUUIDPath = uuidPath
		defer func() {
			randutil.KernelUUIDPath = KernelUUIDPathDefault
		}()

		value, err := randutil.RandomKernelUUID()
		c.Check(value, Equals, KernelTestUUID)
		c.Assert(err, IsNil)
	}
}

func (s *cryptoRandutilSuite) TestRandomKernelUUIDReal(c *C) {
	if _, err := os.Stat(randutil.KernelUUIDPath); err != nil {
		c.Skip("Kernel UUID procfs file is not accessible in the current test environment")
	}

	value, err := randutil.RandomKernelUUID()
	c.Check(value, Not(Equals), "")
	// https://www.rfc-editor.org/rfc/rfc4122#section-3
	// We are not testing the kernel here, so minimal check:
	// UUID should be 36 bytes in length exactly.
	c.Check(value, HasLen, 36)
	c.Assert(err, IsNil)
}
