package tests

import (
	"io"
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
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

	// write string and try to read it at 0 position
	strb := []byte("123abc")
	readed := make([]byte, len(strb))

	_, n, err := f.Write(strb)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)

	err = f.Sync()
	assert.Nil(t, err)

	err = f.Truncate(2)
	assert.Nil(t, err)

	f1i, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), f1i.Size())

	f2i, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), f2i.Size())

	ret, err := f2.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), ret)

	_, n, err = f2.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, 2, n)

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}
