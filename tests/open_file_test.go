package tests

import (
	"io"
	"os"
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestOpenFile(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close()

	f, err := gofilemirror.Create("/tmp/testing_file_mirror")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	assert.True(t, fm.AddReadingFile(f))

	err = f.Close()
	assert.Nil(t, err)
	assert.True(t, fm.RemoveReadingFile(f))

	f, err = gofilemirror.OpenFile(
		"/tmp/testing_file_mirror",
		os.O_APPEND|os.O_RDWR,
		0777,
	)
	assert.Nil(t, err)
	assert.NotNil(t, f)

	assert.True(t, fm.AddReadingFile(f))
	assert.True(t, fm.AddWritingFile(f))

	strb := []byte("123abc")
	readed := make([]byte, len(strb))

	ops, n, err := f.Write(strb)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)

	ops, err = f.Sync()
	assert.Nil(t, err)
	assert.Empty(t, ops)

	ops, ret, err := f.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	assert.Zero(t, ret)
	assert.Empty(t, ops)

	ops, n, err = f.Read(readed)
	assert.Nil(t, err)
	assert.Equal(t, len(strb), n)
	assert.Empty(t, ops)

	err = f.Close()
	assert.Nil(t, err)
	assert.True(t, fm.RemoveReadingFile(f))
	assert.True(t, fm.RemoveWritingFile(f))
}
