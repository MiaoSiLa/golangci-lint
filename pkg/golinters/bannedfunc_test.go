package golinters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	var b = []byte("linters-settings:\n  bannedfunc:\n    (time).Now: \"不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()\"\n    (github.com/Missevan/missevan-go/util/time).TimeNow: \"xxxx\"")
	require.NotPanics(func() {
		setting := decodeFile(b)
		require.NotNil(setting.LinterSettings)
		val, ok := setting.LinterSettings.Funcs["(time).Now"]
		assert.True(ok)
		assert.Equal("不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()", val)
		val, ok = setting.LinterSettings.Funcs["(github.com/Missevan/missevan-go/util/time).TimeNow"]
		assert.True(ok)
		assert.Equal("xxxx", val)
	})
}

func TestConfigToConfigMap(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	var m = make(map[string]string)
	m["(time).Now"] = "不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()"
	m["(github.com/Missevan/missevan-go/util/time).TimeNow"] = "xxxx"
	m["()."] = "().情况"
	m[")."] = ").情况"
	s := configSetting{LinterSettings: BandFunc{Funcs: m}}
	setting := configToConfigMap(s)
	require.NotNil(setting["time"])
	require.NotNil(setting["time"]["Now"])
	assert.Equal("不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()", setting["time"]["Now"])
	require.NotNil(setting["github.com/Missevan/missevan-go/util/time"])
	require.NotNil(setting["github.com/Missevan/missevan-go/util/time"]["TimeNow"])
	assert.Equal("xxxx", setting["github.com/Missevan/missevan-go/util/time"]["TimeNow"])
	assert.Len(setting, 2)
	assert.Nil(setting["()."])
	assert.Nil(setting[")."])
}
