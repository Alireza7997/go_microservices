package file_manager

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type FileManagerInterface interface {
	SetSourcePath(path string)
	CurrentDirectory() string
	CurrentDirectoryPath() string
	Join(path string)
	GoToParent()
	ListDirectories() []string
	ListFiles() []string
	Exists(fullName string) bool
	CreateFile(fullName string) (*os.File, error)
	DeleteFile(fullName string) error
	CopyFile(file, copyFileName string) error
	ReadFile(fullName string) ([]byte, error)
	RenameFileOrDirectory(oldName, newName string) error
	MoveFileOrDirectory(fileName, destPath string) error
	CreateDirectory(name string) error
	DeleteDirectory(name string) error
	GetFileInfo(fullName string) (fs.FileInfo, error)
	SearchFile(fullName string) (fs.DirEntry, error)
	SearchFiles(names ...string) ([]fs.DirEntry, error)
}

type FileManager struct {
	source string
}

var manager = FileManager{}

func New() FileManagerInterface {
	return &manager
}

func (fm *FileManager) SetSourcePath(path string) {
	fm.source = path
}

func (fm FileManager) CurrentDirectory() string {
	return fm.source
}

func (fm FileManager) CurrentDirectoryPath() string {
	abs, _ := filepath.Abs(fm.source)
	return abs
}

func (fm *FileManager) Join(path string) {
	fm.source += "/" + path
}

func (fm *FileManager) GoToParent() {
	fm.source = filepath.Dir(fm.source)
}

func (fm FileManager) ListDirectories() []string {
	direcories, _ := os.ReadDir(fm.source)

	list := make([]string, len(direcories))

	for i, directory := range direcories {
		if directory.IsDir() {
			list[i] = directory.Name()
		}
	}

	return list
}

func (fm FileManager) ListFiles() []string {
	files, _ := os.ReadDir(fm.source)

	list := make([]string, len(files))

	for i, file := range files {
		if file.Type().IsRegular() {
			list[i] = file.Name()
		}
	}

	return list
}

func (fm FileManager) Exists(fullName string) bool {
	_, err := os.Stat(fm.source + "/" + fullName)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (fm FileManager) CreateFile(fullName string) (*os.File, error) {
	if ok := fm.Exists(fullName); !ok {
		file, err := os.Create(fm.source + fullName)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return nil, os.ErrExist
}

func (fm FileManager) DeleteFile(fullName string) error {
	if ok := fm.Exists(fullName); ok {
		err := os.Remove(fm.source + fullName)
		if err != nil {
			return err
		}

		return nil
	}

	return os.ErrNotExist
}

func (fm FileManager) CopyFile(file, copyFileName string) error {
	if ok := fm.Exists(file); ok {
		source, err := os.Open(fm.source + file)
		if err != nil {
			return err
		}

		destination, err := os.Create(fm.source + copyFileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}
	}

	return os.ErrNotExist
}

func (fm FileManager) ReadFile(fullName string) ([]byte, error) {
	if ok := fm.Exists(fullName); ok {
		data, err := os.ReadFile(fm.source + fullName)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	return nil, os.ErrNotExist
}

func (fm FileManager) RenameFileOrDirectory(oldName, newName string) error {
	if ok := fm.Exists(oldName); ok {
		err := os.Rename(fm.source+oldName, fm.source+newName)
		if err != nil {
			return err
		}

		return nil
	}

	return os.ErrNotExist
}

func (fm FileManager) MoveFileOrDirectory(fileName, destPath string) error {
	if ok := fm.Exists(fileName); ok {
		err := os.Rename(fm.source+fileName, fm.source+destPath)
		if err != nil {
			return err
		}

		return nil
	}

	return os.ErrNotExist
}

func (fm FileManager) CreateDirectory(name string) error {
	if ok := fm.Exists(name); !ok {

		err := os.MkdirAll(fm.source+name, os.ModePerm)
		if err != nil {
			return err
		}

		return nil
	}

	return os.ErrExist
}

func (fm FileManager) DeleteDirectory(name string) error {
	if ok := fm.Exists(name); ok {
		entities, err := os.ReadDir(fm.source + name)
		if err != nil {
			return err
		}

		if len(entities) == 0 {
			err = os.Remove(fm.source + name)
			if err != nil {
				return err
			}

			return nil
		} else {

			return errors.New("directory is not empty")
		}
	}

	return os.ErrNotExist
}

func (fm FileManager) GetFileInfo(fullName string) (fs.FileInfo, error) {
	if ok := fm.Exists(fullName); ok {
		info, err := os.Stat(fm.source + fullName)
		if err != nil {
			return nil, err
		}

		return info, nil
	}

	return nil, os.ErrNotExist
}

func (fm FileManager) SearchFile(fullName string) (fs.DirEntry, error) {
	entities, err := os.ReadDir(fm.source)
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		if entity.Type().IsRegular() && entity.Name() == fullName {
			return entity, nil
		}
	}

	return nil, os.ErrNotExist
}

func (fm FileManager) SearchFiles(names ...string) ([]fs.DirEntry, error) {
	files := make([]fs.DirEntry, len(names))
	filesFound := make(map[string]bool)

	entities, err := os.ReadDir(fm.source)
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		if entity.Type().IsRegular() {

			fileName := entity.Name()
			if filesFound[fileName] {
				continue
			}

			for _, name := range names {
				if fileName == name {
					filesFound[fileName] = true
					files = append(files, entity)
				}
			}
		}
	}

	return files, nil
}
