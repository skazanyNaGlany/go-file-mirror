package tests

import (
	"io"
	"os"
	"sync"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestWriteString(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close()

	f, err := os.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := os.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	fm.SetReadingFile(f)
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	fm.SetFileMutex(f, &sync.Mutex{})
	fm.SetFileMutex(f2, &sync.Mutex{})

	// write string and try to read it at 0 position
	str := "123abc"
	readed := make([]byte, len(str))

	ops, n, err := fm.WriteString(str)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)
	assert.Empty(t, ops)

	ops, err = fm.Sync()
	assert.Nil(t, err)
	assert.Empty(t, ops)

	ops, ret, err := fm.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	assert.Zero(t, ret)
	assert.Empty(t, ops)

	ops, n, err = fm.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)
	assert.Empty(t, ops)

	// write other string at 2 position
	str2 := "defghi"
	readed = make([]byte, len(str2))

	ops, ret, err = fm.Seek(2, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), ret)
	assert.Empty(t, ops)

	ops, n, err = fm.WriteString(str2)
	assert.Nil(t, err)
	assert.Equal(t, len(str), n)
	assert.Empty(t, ops)

	ops, err = fm.Sync()
	assert.Nil(t, err)
	assert.Empty(t, ops)

	ops, ret, err = fm.Seek(2, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), ret)
	assert.Empty(t, ops)

	ops, n, err = fm.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, len(str2), n)
	assert.Equal(t, string(readed), str2)
	assert.Empty(t, ops)

	f1i, err := fm.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(8), f1i.Size())
}
