// Copyright 2017-2019 Lei Ni (nilei81@gmail.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileutil

import (
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
)

const (
	// DefaultFileMode is the default file mode for files generated by
	// Dragonboat.
	DefaultFileMode    = 0640
	defaultDirFileMode = 0750
	deleteFilename     = "DELETED.dragonboat"
)

// Exist returns whether the specified filesystem entry exists.
func Exist(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// MkdirAll creates the specified dir along with any necessary parents.
func MkdirAll(dir string) error {
	exist, err := Exist(dir)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	parent := filepath.Dir(dir)
	exist, err = Exist(parent)
	if err != nil {
		return err
	}
	if !exist {
		if err := MkdirAll(parent); err != nil {
			return err
		}
	}
	return Mkdir(dir)
}

// Mkdir creates the specified dir.
func Mkdir(dir string) error {
	if err := os.Mkdir(dir, defaultDirFileMode); err != nil {
		return err
	}
	return SyncDir(filepath.Dir(dir))
}

// SyncDir calls fsync on the specified directory.
func SyncDir(dir string) (err error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		panic("not a dir")
	}
	df, err := os.Open(filepath.Clean(dir))
	if err != nil {
		return err
	}
	defer func() {
		if cerr := df.Close(); err == nil {
			err = cerr
		}
	}()
	return df.Sync()
}

// MarkDirAsDeleted marks the specified directory as deleted.
func MarkDirAsDeleted(dir string, msg proto.Message) error {
	return CreateFlagFile(dir, deleteFilename, msg)
}

// IsDirMarkedAsDeleted returns a boolean flag indicating whether the specified
// directory has been marked as deleted.
func IsDirMarkedAsDeleted(dir string) (bool, error) {
	fp := filepath.Join(dir, deleteFilename)
	return Exist(fp)
}
