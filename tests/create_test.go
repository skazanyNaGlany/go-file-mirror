package tests

import (
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	fm := gofilemirror.NewFileMirror(FILE_MIRROR_QUEUE_SIZE)
	defer fm.Close()

	f, err := gofilemirror.Create("/tmp/testing_file_mirror")
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.True(t, fm.AddReadingFile(f))

	err = f.Close()
	assert.Nil(t, err)
	assert.True(t, fm.RemoveReadingFile(f))
}
