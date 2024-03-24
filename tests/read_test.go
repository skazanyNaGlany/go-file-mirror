package tests

import (
	"io"
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
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

	n, err := f.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)

	n, err = f2.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)

	n, err = f2.Write(strb)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)

	err = f2.Sync()
	assert.Nil(t, err)

	// cannot read at EOF
	n, err = f.Read(readed)
	assert.NotNil(t, err)
	assert.Zero(t, n)

	ret, err := f.Seek(4, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(4), ret)

	n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, n, 5)

	ret, err = f2.Seek(5, io.SeekStart)
	assert.Nil(t, err)
	assert.Equal(t, int64(5), ret)

	n, err = f2.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, n, 4)

	err = f.Close()
	assert.Nil(t, err)

	// all files within that FileMirror instance
	// have been closed, calling Close() again
	// should return an error
	err = f2.Close()
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &os.ErrClosed)
}
