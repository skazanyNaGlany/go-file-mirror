package tests

import (
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestRemoveFile(t *testing.T) {
	fm := gofilemirror.NewFileMirror()

	f, err := fm.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	f2, err := fm.CreateTemp("/tmp", "testing_file_mirror")
	if err != nil {
		panic(err)
	}

	assert.True(t, fm.HasFile(f))
	assert.True(t, fm.HasFile(f2))

	assert.True(t, fm.RemoveFile(f))
	assert.True(t, fm.RemoveFile(f2))

	assert.False(t, fm.HasFile(f))
	assert.False(t, fm.HasFile(f2))

	// after removal from FileMirror we must close
	// files manually, since they are not managed
	// by FileMirror anymore
	err = f.GetUnderlyingFile().Close()
	assert.Nil(t, err)

	err = f2.GetUnderlyingFile().Close()
	assert.Nil(t, err)
}
