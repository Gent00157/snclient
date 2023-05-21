package snclient

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckUnknown(t *testing.T) {
	snc := Agent{}
	res := snc.RunCheck("not_there", []string{})
	assert.Equalf(t, CheckExitUnknown, res.State, "state Unknown")
	assert.Regexpf(t,
		regexp.MustCompile(`^UNKNOWN - No such check: not_there`),
		string(res.BuildPluginOutput()),
		"output matches",
	)
}

func TestCheckSNClientVersion(t *testing.T) {
	snc := Agent{}
	res := snc.RunCheck("check_snclient_version", []string{})
	assert.Equalf(t, CheckExitOK, res.State, "state OK")
	assert.Regexpf(t,
		regexp.MustCompile(`^SNClient\+ v\d+`),
		string(res.BuildPluginOutput()),
		"output matches",
	)
}

func TestCheckCPU(t *testing.T) {
	snc := StartTestAgent(t, "")

	res := snc.RunCheck("check_cpu", []string{"warn=load = 101", "crit=load = 102"})
	assert.Equalf(t, CheckExitOK, res.State, "state OK")
	assert.Regexpf(t,
		regexp.MustCompile(`^OK: CPU load is ok. \|'total 5m'=\d+%;101;102 'total 1m'=\d+%;101;102 'total 5s'=\d+%;101;102$`),
		string(res.BuildPluginOutput()),
		"output matches",
	)

	StopTestAgent(t, snc)
}

func TestCheckAlias(t *testing.T) {
	config := `
[/modules]
CheckExternalScripts = enabled

[/settings/external scripts/alias]
alias_cpu = check_cpu warn=load=101 crit=load=102
`
	snc := StartTestAgent(t, config)

	res := snc.RunCheck("alias_cpu", []string{})
	assert.Equalf(t, CheckExitOK, res.State, "state OK")
	assert.Regexpf(t,
		regexp.MustCompile(`^OK: CPU load is ok. \|'total 5m'=\d+%;101;102 'total 1m'=\d+%;101;102 'total 5s'=\d+%;101;102$`),
		string(res.BuildPluginOutput()),
		"output matches",
	)

	StopTestAgent(t, snc)
}

func TestCheckUptime(t *testing.T) {
	snc := Agent{}
	res := snc.RunCheck("check_uptime", []string{})
	assert.Regexpf(t,
		regexp.MustCompile(`^\w+: uptime:.*?(\d+w \d+d|\d+:\d+h), boot: \d+\-\d+\-\d+ \d+:\d+:\d+ \(UTC\) \|'uptime'=\d+s;172800;86400`),
		string(res.BuildPluginOutput()),
		"output matches",
	)
}