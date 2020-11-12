package golinters

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
)

func TestDecodeFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	b := []byte("linters-settings:\n  bannedfunc:\n    (time).Now: \"xxxx\"\n    (github.com/Missevan/missevan-go/util).TimeNow: \"xxxx\"")
	var setting configSetting
	require.NotPanics(func() { setting = decodeFile(b) })
	require.NotNil(setting.LinterSettings)
	val := setting.LinterSettings.Funcs["(time).Now"]
	assert.NotEmpty(val)
	val = setting.LinterSettings.Funcs["(github.com/Missevan/missevan-go/util).TimeNow"]
	assert.NotEmpty(val)
}

func TestConfigToConfigMap(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := map[string]string{
		"(time).Now": "不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()",
		"(github.com/Missevan/missevan-go/util).TimeNow": "xxxx",
		"().": "(). 情况",
		").":  "). 情况",
	}
	s := configSetting{LinterSettings: BandFunc{Funcs: m}}
	setting := configToConfigMap(s)
	require.Len(setting, 2)
	require.NotNil(setting["time"])
	require.NotNil(setting["time"]["Now"])
	assert.Equal("不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()", setting["time"]["Now"])
	require.NotNil(setting["github.com/Missevan/missevan-go/util"])
	require.NotNil(setting["github.com/Missevan/missevan-go/util"]["TimeNow"])
	assert.Equal("xxxx", setting["github.com/Missevan/missevan-go/util"]["TimeNow"])
	assert.Nil(setting["()."])
	assert.Nil(setting[")."])
}

func TestGetUsedMap(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pkg := types.NewPackage("test", "test")
	importPkg := []*types.Package{types.NewPackage("time", "time"),
		types.NewPackage("github.com/Missevan/missevan-go/util", "util")}
	pkg.SetImports(importPkg)
	pass := analysis.Pass{Pkg: pkg}
	m := map[string]map[string]string{
		"time":                                 {"Now": "xxxx", "Date": "xxxx"},
		"assert":                               {"New": "xxxx"},
		"github.com/Missevan/missevan-go/util": {"TimeNow": "xxxx"},
	}
	usedMap := getUsedMap(&pass, m)
	require.Len(usedMap, 2)
	require.Len(usedMap["time"], 2)
	assert.NotEmpty(usedMap["time"]["Now"])
	assert.NotEmpty(usedMap["time"]["Date"])
	require.Len(usedMap["util"], 1)
	assert.NotEmpty(usedMap["util"]["TimeNow"])
}
