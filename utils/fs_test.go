package utils

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestCreateDirs(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)
	tempDirPath, err := ioutil.TempDir("", "create_dirs_test")
	asserts.Nil(err)

	directories := []string{
		path.Join(tempDirPath, "test_dir_1"),
		path.Join(tempDirPath, "test_dir_2"),
		path.Join(tempDirPath, "test_dir_3"),
		path.Join(tempDirPath, "test_dir_4"),
	}

	t.Run("CreateDirectories", func(t *testing.T) {
		err := CreateDirs(0755, directories...)
		asserts.Nil(err)

		for _, dir := range directories {
			info, err := os.Stat(dir)
			asserts.Nil(err)
			asserts.True(info.IsDir())
			asserts.EqualValues(0755, uint32(info.Mode().Perm()))
		}
	})

	if err := os.RemoveAll(tempDirPath); err != nil {
		log.Printf("Error while removing tempr directory: %s, %v\n", tempDirPath, err)
		t.Fail()
	}
}

func TestDirExists(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)
	tempDirPath, err := ioutil.TempDir("", "create_dirs_test")
	asserts.Nil(err)
	asserts.Nil(os.MkdirAll(path.Join(tempDirPath, "test_dir_1"), 0755))
	_, err = os.Create(path.Join(tempDirPath, "test.txt"))
	asserts.Nil(err)

	data := []struct {
		Dir    string
		Exists bool
	}{
		{
			Dir:    path.Join(tempDirPath, "test_dir_1"),
			Exists: true,
		},
		{
			Dir:    path.Join(tempDirPath, "not_existing"),
			Exists: false,
		},
		{
			Dir:    path.Join(tempDirPath, "test.txt"),
			Exists: false,
		},
	}

	for _, d := range data {
		asserts.Equal(d.Exists, DirExists(d.Dir))
	}

	if err := os.RemoveAll(tempDirPath); err != nil {
		log.Printf("Error while removing tempr directory: %s, %v\n", tempDirPath, err)
		t.Fail()
	}
}
