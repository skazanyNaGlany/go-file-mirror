package tests

import (
	"io"
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestWriteString(t *testing.T) {
	fm := gofilemirror.NewFileMirror()

	f, err := gofilemirror.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := gofilemirror.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddReadingFile(f2))
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	// write string and try to read it at 0 position
	str := "123abc"
	readed := make([]byte, len(str))

	n, err := f.WriteString(str)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)

	err = f.Sync()
	assert.Nil(t, err)

	ret, err := f.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	assert.Zero(t, ret)

	n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)

	// write other string at 2 position
	str2 := "defghi"
	readed = make([]byte, len(str2))

	ret, err = f2.Seek(2, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), ret)

	n, err = f2.WriteString(str2)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)

	err = f.Sync()
	assert.Nil(t, err)

	ret, err = f2.Seek(2, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), ret)

	n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, len(str2), n)
	assert.Equal(t, string(readed), str2)

	f1i, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(8), f1i.Size())

	f2i, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, f1i.Size(), f2i.Size())

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}
