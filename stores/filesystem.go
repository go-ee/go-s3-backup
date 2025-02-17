/*
Copyright 2018 codestation

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

package stores

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	log "unknwon.dev/clog/v2"
)

// FilesystemConfig has the config options for the FilesystemConfig service
type FilesystemConfig struct {
	SaveDir string
}

// Store moves/copies a file to another directory
func (f *FilesystemConfig) Store(src, prefix, filename string) error {
	dest := path.Clean(path.Join(f.SaveDir, prefix, filename))

	if src == dest {
		log.Trace("Using the same path as source and destination, do nothing")
		return nil
	}

	err := os.Rename(dest, src)
	if err != nil {
		log.Warn("Cannot rename %s to %s, trying to copy instead", src, dest)
	} else {
		return nil
	}

	var removeSourceFile = false

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open source file %s, %v", src, err)
	}

	defer func() {
		srcFile.Close()

		// if there aren't any errors on the file copy then delete the source file
		if removeSourceFile {
			if err = os.Remove(src); err != nil {
				log.Warn("Cannot remove source file %s", src)
			}
		}
	}()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("cannot create destination file %s, %v", dest, err)
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return fmt.Errorf("error while copying file, %v", err)
	}

	if err = destFile.Sync(); err != nil {
		return fmt.Errorf("cannot flush file contents, %v", err)
	}

	removeSourceFile = true

	return nil
}

func (f *FilesystemConfig) getFileListing(basedir, namePrefix string) ([]string, error) {
	fullBasedir := path.Join(f.SaveDir, basedir)
	files, err := ioutil.ReadDir(fullBasedir)
	if err != nil {
		return nil, fmt.Errorf("cannot list contents of directory %s, %v", f.SaveDir, err)
	}
	re := generatePattern(namePrefix)

	var filenames []string
	for _, f := range files {
		if !f.IsDir() {
			// ignore files not created by this program
			if re.MatchString(f.Name()) {
				filenames = append(filenames, path.Join(fullBasedir, f.Name()))
			}
		}
	}

	return filenames, nil
}

// RemoveOlderBackups keeps the most recent backups of a directory and deletes the old ones
func (f *FilesystemConfig) RemoveOlderBackups(basedir, namePrefix string, keep int) error {
	filePaths, err := f.getFileListing(basedir, namePrefix)
	if err != nil {
		return err
	}

	if len(filePaths) == 0 {
		return nil
	}

	count := len(filePaths) - keep
	deleted := 0

	if count > 0 {
		for _, file := range filePaths[:count] {
			fullpath := path.Clean(file)
			err = os.Remove(fullpath)
			if err != nil {
				log.Error("Failed to remove file %s", fullpath)
			} else {
				deleted++
			}
		}

		log.Trace("Deleted %d objects from %s", deleted, path.Join(f.SaveDir, basedir))
	}

	return nil
}

// FindLatestBackup returns the most recent backup of the specified directory
func (f *FilesystemConfig) FindLatestBackup(basedir, namePrefix string) (string, error) {
	files, err := f.getFileListing(basedir, namePrefix)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("cannot find a recent backup on %s", f.SaveDir)
	}

	return files[len(files)-1], nil
}

// Retrieve returns the path of the requested file
func (f *FilesystemConfig) Retrieve(filename string) (string, error) {
	return path.Clean(path.Join(f.SaveDir, filename)), nil
}

// Close deinitializes the store (no dothing)
func (f *FilesystemConfig) Close() {
}
