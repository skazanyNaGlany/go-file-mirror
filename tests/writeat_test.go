package tests

import (
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestWriteAt(t *testing.T) {
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
	strb2 := []byte("456def")
	readed := make([]byte, len(strb))

	ops, n, err := f.WriteAt(strb, 0)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)

	ops, n, err = f2.WriteAt(strb2, int64(len(strb)))
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)

	ops, n, err = f.ReadAt(readed, int64(len(strb)))
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)
	assert.Equal(t, strb2, readed)

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
