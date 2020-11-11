package golinters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	b := []byte("linters-settings:\n  bannedfunc:\n    (time).Now: \"xxxx\"\n    (github.com/Missevan/missevan-go/util/time).TimeNow: \"xxxx\"")
	var setting configSetting
	require.NotPanics(func() { setting = decodeFile(b) })
	require.NotNil(setting.LinterSettings)
	val := setting.LinterSettings.Funcs["(time).Now"]
	assert.NotEmpty(val)
	val = setting.LinterSettings.Funcs["(github.com/Missevan/missevan-go/util/time).TimeNow"]
	assert.NotEmpty(val)
}

func TestConfigToConfigMap(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	m := map[string]string{
		"(time).Now": "不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()",
		"(github.com/Missevan/missevan-go/util/time).TimeNow": "xxxx",
		"().": "(). 情况",
		").":  "). 情况",
	}
	s := configSetting{LinterSettings: BandFunc{Funcs: m}}
	setting := configToConfigMap(s)
	require.Len(setting, 2)
	require.NotNil(setting["time"])
	require.NotNil(setting["time"]["Now"])
	assert.Equal("不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()", setting["time"]["Now"])
	require.NotNil(setting["github.com/Missevan/missevan-go/util/time"])
	require.NotNil(setting["github.com/Missevan/missevan-go/util/time"]["TimeNow"])
	assert.Equal("xxxx", setting["github.com/Missevan/missevan-go/util/time"]["TimeNow"])
	assert.Nil(setting["()."])
	assert.Nil(setting[")."])
}
