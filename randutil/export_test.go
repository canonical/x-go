// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2022 Canonical Ltd
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

package randutil

// Needed by TestRandomKernelUUIDReal() to check if access
// to kernel path exist on test host.
var KernelUUIDPath = kernelUUIDPath

func MockKernelUUIDPath(newPath string) (restore func()) {
	kernelUUIDPathDefault := kernelUUIDPath
	kernelUUIDPath = newPath
	return func() {
		kernelUUIDPath = kernelUUIDPathDefault
	}
}
