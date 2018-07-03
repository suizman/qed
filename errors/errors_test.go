/*
   Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package balloon implements the tree interface to interact with both hyper
// and history trees.
package errors_test

import (
	"fmt"
	"testing"

	"github.com/bbva/qed/errors"
)

func TestSingleError(t *testing.T) {
	err1 := errors.E(errors.Op("Get"), "Unable to get object")
	fmt.Println(err1)
}
