package tests

import (
	"testing"

	gofilemirror "github.com/skazanyNaGlany/go-file-mirror"
	"github.com/stretchr/testify/assert"
)

func TestCreateTempSingle(t *testing.T) {
	fm := gofilemirror.NewFileMirror()

	f, err := fm.CreateTemp("/tmp", "testing_file_mirror")

	if err != nil {
		panic(err)
	}

	err = f.Close()

	assert.Nil(t, err)
}
