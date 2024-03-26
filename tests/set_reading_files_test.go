package tests

import (
	"io"
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestSetReadingFiles(t *testing.T) {
	fm := gofilemirror.NewFileMirror()

	f, err := fm.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := fm.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	fm.SetReadingFiles([]gofilemirror.IFile{})

	readed := make([]byte, 5)

	n, err := f2.Read(readed)
	assert.ErrorAs(t, gofilemirror.ErrNoFilesToRead, &err)
	assert.Zero(t, n)

	fm.SetReadingFiles([]gofilemirror.IFile{f, f2})

	// no data to read so it will return just EOF error
	n, err = f2.Read(readed)
	assert.ErrorAs(t, io.EOF, &err)
	assert.Zero(t, n)

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}
