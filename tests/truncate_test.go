package tests

import (
	"io"
	"os"
	"sync"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
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
	strb := []byte("123abc")
	readed := make([]byte, len(strb))

	ops, n, err := fm.Write(strb, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)

	ops, err = fm.Sync(nil)
	assert.Nil(t, err)
	assert.Empty(t, ops)

	ops, err = fm.Truncate(2, nil)
	assert.Nil(t, err)
	assert.Empty(t, ops)

	f1i, err := fm.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), f1i.Size())

	f2i, err := fm.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), f2i.Size())

	ops, ret, err := fm.Seek(0, io.SeekStart, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), ret)
	assert.Empty(t, ops)

	ops, n, err = fm.Read(readed, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, n)
	assert.Empty(t, ops)
}
