/*
 * git_test.go
 *
 * Copyright 2021 Bill Zissimopoulos
 */
/*
 * This file is part of Hubfs.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

package git

import (
	"fmt"
	"os"
	"testing"

	"github.com/billziss-gh/golib/keyring"
)

const remote = "https://github.com/billziss-gh/hubfs"
const refName = "refs/heads/master"

const hash0 = "90f898ae1f8d3c976f9224d92e3b08d7813e961e"
const hash1 = "609d3b892764952ef69676e653e06b2ca904be18"

var token string

func TestGetRefs(t *testing.T) {
	repository, err := OpenRepository(remote, token)
	if nil != err {
		t.Error(err)
	}
	defer repository.Close()

	refs, err := repository.GetRefs()
	if nil != err {
		t.Error(err)
	}
	found := false
	for n := range refs {
		if n == refName {
			found = true
			break
		}
	}
	if !found {
		t.Error()
	}

	refs, err = repository.GetRefs()
	if nil != err {
		t.Error(err)
	}
	found = false
	for n := range refs {
		if n == refName {
			found = true
			break
		}
	}
	if !found {
		t.Error()
	}
}

func TestFetchObjects(t *testing.T) {
	repository, err := OpenRepository(remote, token)
	if nil != err {
		t.Error(err)
	}
	defer repository.Close()

	wants := []string{
		hash0,
		hash1,
	}
	found := false
	err = repository.FetchObjects(wants,
		func(hash string, content []byte) error {
			if hash0 == hash {
				found = true
				_, err := DecodeCommit(content)
				if nil != err {
					return err
				}
			}
			if hash1 == hash {
				found = true
				_, err := DecodeTree(content)
				if nil != err {
					return err
				}
			}
			return nil
		})
	if nil != err {
		t.Error(err)
	}
	if !found {
		t.Error()
	}

	wants = []string{
		hash0,
	}
	found = false
	err = repository.FetchObjects(wants,
		func(hash string, content []byte) error {
			if hash0 == hash {
				found = true
				_, err := DecodeCommit(content)
				if nil != err {
					return err
				}
			}
			return nil
		})
	if nil != err {
		t.Error(err)
	}
	if !found {
		t.Error()
	}

	wants = []string{
		hash1,
	}
	found = false
	err = repository.FetchObjects(wants,
		func(hash string, content []byte) error {
			if hash1 == hash {
				found = true
				_, err := DecodeTree(content)
				if nil != err {
					return err
				}
			}
			return nil
		})
	if nil != err {
		t.Error(err)
	}
	if !found {
		t.Error()
	}
}

func TestMain(m *testing.M) {
	var err error
	token, err = keyring.Get("hubfs", "https://github.com")
	if nil != err {
		fmt.Fprintf(os.Stderr, "unable to get auth token: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
