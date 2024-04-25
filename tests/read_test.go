package tests

import (
	"io"
	"os"
	"sync"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close()

	f, err := os.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := os.CreateTemp("/tmp", "testing_file_mirror2")
	if err != nil {
		panic(err)
	}

	fm.SetReadingFile(f)
	assert.True(t, fm.AddWritingFile(f))
	assert.True(t, fm.AddWritingFile(f2))

	fm.SetFileMutex(f, &sync.Mutex{})
	fm.SetFileMutex(f2, &sync.Mutex{})

	// read at 0 position
	strb := []byte("123456abc")
	readed := make([]byte, 6)

	ops, n, err := fm.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)
	assert.Empty(t, ops)

	ops, n, err = fm.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)
	assert.Empty(t, ops)

	ops, n, err = fm.Write(strb)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)

	ops, err = fm.Sync()
	assert.Nil(t, err)
	assert.Empty(t, ops)

	// cannot read at EOF
	ops, n, err = fm.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)
	assert.Empty(t, ops)

	ops, ret, err := fm.Seek(4, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(4), ret)
	assert.Empty(t, ops)

	ops, n, err = fm.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, n, 5)
	assert.Empty(t, ops)

	ops, ret, err = fm.Seek(5, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(5), ret)
	assert.Empty(t, ops)

	ops, n, err = fm.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, n, 4)
	assert.Empty(t, ops)
}
