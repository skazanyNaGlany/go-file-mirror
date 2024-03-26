package tests

import (
	"io"
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	fm := gofilemirror.NewFileMirror()

	f, err := os.CreateTemp("/tmp", "testing_file_mirror")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	f2 := fm.NewFile(uintptr(f.Fd()), f.Name())
	assert.NotNil(t, f2)

	writingFiles := fm.GetWritingFiles()
	writingFiles = append(writingFiles, f2)
	fm.SetWritingFiles(writingFiles)

	str := "123abc"
	readed := make([]byte, len(str))

	n, err := f2.WriteString(str)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)

	ret, err := f2.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	assert.Zero(t, ret)

	n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)
	assert.Equal(t, str, string(readed))

	err = f.Close()
	assert.Nil(t, err)

	f2.Close()
	assert.Nil(t, err)
}
