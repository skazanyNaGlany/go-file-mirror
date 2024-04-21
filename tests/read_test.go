package tests

import (
	"io"
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close()

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

	// read at 0 position
	strb := []byte("123456abc")
	readed := make([]byte, 6)

	ops, n, err := f.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)
	assert.Empty(t, ops)

	ops, n, err = f2.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)
	assert.Empty(t, ops)

	ops, n, err = f2.Write(strb)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)

	ops, err = f2.Sync()
	assert.Nil(t, err)
	assert.Empty(t, ops)

	// cannot read at EOF
	ops, n, err = f.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)
	assert.Empty(t, ops)

	ops, ret, err := f.Seek(4, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(4), ret)
	assert.Empty(t, ops)

	ops, n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, n, 5)
	assert.Empty(t, ops)

	ops, ret, err = f2.Seek(5, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(5), ret)
	assert.Empty(t, ops)

	ops, n, err = f2.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, n, 4)
	assert.Empty(t, ops)

	err = f.Close()
	assert.Nil(t, err)

	assert.True(t, fm.RemoveReadingFile(f))
	assert.True(t, fm.RemoveWritingFile(f))

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)

	assert.True(t, fm.RemoveReadingFile(f2))
	assert.True(t, fm.RemoveWritingFile(f2))
}
