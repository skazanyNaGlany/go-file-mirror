package tests

import (
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestReadAt(t *testing.T) {
	fm := gofilemirror.NewFileMirror()

	f, err := fm.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := fm.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	// read at 0 position
	strb := []byte("123456abc")
	readed := make([]byte, 6)

	n, err := f.ReadAt(readed, 0)
	assert.NotNil(t, err)
	assert.Zero(t, n)

	n, err = f2.ReadAt(readed, 0)
	assert.NotNil(t, err)
	assert.Zero(t, n)

	n, err = f2.Write(strb)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)

	err = f2.Sync()
	assert.Nil(t, err)

	n, err = f.ReadAt(readed, 2)
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	n, err = f2.ReadAt(readed, 3)
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}
