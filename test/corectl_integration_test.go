// +build integration

package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/qlik-oss/corectl/test/toolkit"
)

func TestBasicAnalyzing(t *testing.T) {

	os.Setenv("CONN_TYPE", "folder")
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()

	p.ExpectOK().Run("build")

	p.ExpectGolden().Run("reload")
	p.ExpectGolden().Run("reload", "--silent")
	p.ExpectGolden().Run("reload", "--silent", "--no-save")

	p.ExpectIncludes("Connected to ", "The data model has 2 tables.").Run("status")
	p.ExpectGolden().Run("tables")
	p.ExpectGolden().Run("assoc")
	p.ExpectGolden().Run("fields")
	p.ExpectGolden().Run("values", "numbers")
	p.ExpectGolden().Run("meta")

	p.ExpectGolden().Run("eval", "count(numbers)", "by", "xyz")
	p.ExpectGolden().Run("eval", "count(numbers)")
	p.ExpectGolden().Run("eval", "=1+1")
	p.ExpectGolden().Run("eval", "1+1")
	p.ExpectGolden().Run("eval", "by", "numbers")
}

func TestReload(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()
	p.ExpectIncludes("<<  5 Lines fetched").Run("build")
	p.ExpectIncludes("<<  3 Lines fetched").Run("build", "--limit", "3")
	p.ExpectIncludes("<<  3 Lines fetched").Run("reload", "--limit", "3")
	p.ExpectGolden().Run("reload", "--silent")
	p.ExpectGolden().Run("reload", "--silent", "--no-save")
}

func TestContextManagement(t *testing.T) {
	jwt := toolkit.Params{T: t, Config: "test/projects/using-jwts/corectl.yml", Engine: *toolkit.EngineJwtIP}
	abac := toolkit.Params{T: t, Config: "test/projects/abac/corectl.yml", Engine: *toolkit.EngineAbacIP}
	params := []toolkit.Params{jwt, abac}
	contexts := []string{t.Name() + "_JWT", t.Name() + "_ABAC"}
	for i, p := range params {
		// Create a context using the config and engine url
		p.ExpectOK().Run("context", "set", contexts[i])
		// Remove config and engine from p to see if context stored them
		p.Config, p.Engine = "", ""
		p.ExpectOK().Run("status")
	}

	// Empty params, should default to localhost:9076 when there is no context
	p := toolkit.Params{T: t}
	// Contexts stored locally
	p.ExpectOK().Run("context", "ls")
	// Current context should be abac, connecting to jwt shouldn't work
	p.ExpectIncludes(contexts[1]).Run("context", "get")
	p.ExpectError().Run("status", "--engine", *toolkit.EngineJwtIP)
	// Check if we can update contexts
	p.ExpectOK().Run("context", "set", contexts[1], "--engine", *toolkit.EngineStdIP)
	p.ExpectIncludes(*toolkit.EngineStdIP).Run("status")
	abac.ExpectOK().Run("context", "set", contexts[1])
	p.ExpectIncludes(*toolkit.EngineAbacIP).Run("context", "get")
	// See if all context commands work
	for i, ctx := range contexts {
		p.ExpectOK().Run("context", "get", ctx)
		// With context status should work
		p.ExpectOK().Run("context", "use", ctx)
		p.ExpectIncludes(params[i].Engine).Run("status")
		// Without context status should default to localhost:9076
		p.ExpectOK().Run("context", "clear")
		p.ExpectIncludes("No current context").Run("context", "get")
		p.ExpectOK().Run("context", "rm", ctx)
	}
	// No context here, expecting default
	p.ExpectIncludes("localhost:9076").Run("status")
}

func TestQuietCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/quiet/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	cmds := []string{
		"connection", "dimension", "measure",
		"object", "state", "variable",
	}
	defer p.Reset()
	p.ExpectOK().Run("build")
	for _, cmd := range cmds {
		out := p.ExpectOK().Run(cmd, "ls", "-q")
		ids := bytes.Split(out, []byte("\n"))
		for _, id := range ids {
			if len(id) > 0 {
				p.ExpectOK().Run(cmd, "rm", string(id))
			}
		}
	}
	// As the quiet flag trumps the traffic flag these two commands
	// should be equal.
	out1 := p.ExpectIncludes(t.Name()).Run("app", "ls", "-q")
	out2 := p.ExpectOK().Run("app", "ls", "-q", "-t")
	if string(out1) != string(out2) {
		t.Error("Expected 'corectl app ls -q -t' to be equal to 'corectl app ls -q'")
	}
}

func TestLogBuffer(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/quiet/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	p.ExpectOK().Run("context", "set", t.Name())
	// The quiest flag should mute the warnings
	p.ExpectEmptyOK().Run("app", "ls", "-q")
	// We should have a warning saying something about context here
	p.ExpectIncludes("context").Run("app", "ls")
	p.ExpectOK().Run("context", "rm", t.Name())
}

func TestConnectionManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()
	p.Run("build")
	p.ExpectIncludes(`myconnection | testconnector`).Run("connection", "ls")
	p.ExpectOK().Run("connection", "ls", "--bash")
}

func TestObjectManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()

	// Build with both objects and check
	p.ExpectOK().Run("build")
	p.ExpectGolden().Run("object", "ls")
	p.ExpectGolden().Run("object", "ls", "--json")
	p.ExpectGolden().Run("object", "ls", "--bash")
	p.ExpectGolden().Run("object", "ls", "--json")
	p.ExpectGolden().Run("object", "properties", "object-using-inline")
	p.ExpectGolden().Run("object", "properties", "--minimum", "object-using-inline")
	p.ExpectGolden().Run("object", "layout", "object-using-inline")
	p.ExpectGolden().Run("object", "data", "object-using-inline")
	p.ExpectGolden().Run("object", "data", "object-using-dims-and-measures")
	p.ExpectErrorIncludes("Invalid handle: Invalid Params (-32602").Run("object", "data", "nosuchobject")

	p.ExpectJsonArray("qId", "object-using-dims-and-measures", "object-using-inline").Run("object", "ls", "--json")

	// Remove one object and check
	p.ExpectOK().Run("object", "rm", "object-using-inline")
	p.ExpectJsonArray("qId", "object-using-dims-and-measures").Run("object", "ls", "--json")

	// Re-add the object and check
	p.ExpectOK().Run("object", "set", "test/projects/using-entities/objects.json")
	p.ExpectJsonArray("qId", "object-using-dims-and-measures", "object-using-inline").Run("object", "ls", "--json")
}

func TestDimensionManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()

	// Build with both dimensions and check
	p.ExpectOK().Run("build")
	p.ExpectGolden().Run("dimension", "ls")
	p.ExpectGolden().Run("dimension", "ls", "--json")
	p.ExpectGolden().Run("dimension", "ls", "--bash")
	p.ExpectGolden().Run("dimension", "properties", "dimension-abcs")
	p.ExpectGolden().Run("dimension", "layout", "dimension-abcs")
	p.ExpectJsonArray("qId", "dimension-abcs", "dimension-xyz").Run("dimension", "ls", "--json")

	// Remove one dimension an check
	p.ExpectOK().Run("dimension", "rm", "dimension-abcs")
	p.ExpectJsonArray("qId", "dimension-xyz").Run("dimension", "ls", "--json")

	// Re-add the measure and check
	p.ExpectOK().Run("dimension", "set", "test/projects/using-entities/dimensions.json")
	p.ExpectJsonArray("qId", "dimension-abcs", "dimension-xyz").Run("dimension", "ls", "--json")
}

func TestMeasureManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()

	// Build with both measures and check
	p.ExpectOK().Run("build")
	p.ExpectGolden().Run("measure", "ls")
	p.ExpectGolden().Run("measure", "ls", "--json")
	p.ExpectGolden().Run("measure", "ls", "--bash")
	p.ExpectGolden().Run("measure", "properties", "measure-sum-numbers")
	p.ExpectGolden().Run("measure", "layout", "measure-sum-numbers")
	p.ExpectJsonArray("qId", "measure-count-numbers", "measure-sum-numbers").Run("measure", "ls", "--json")

	// Remove one measure and check
	p.ExpectOK().Run("measure", "rm", "measure-count-numbers")
	p.ExpectJsonArray("qId", "measure-sum-numbers").Run("measure", "ls", "--json")

	// Re-add the measure and check
	p.ExpectOK().Run("measure", "set", "test/projects/using-entities/measures.json")
	p.ExpectJsonArray("qId", "measure-count-numbers", "measure-sum-numbers").Run("measure", "ls", "--json")
}

func TestVariableManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()

	// Build with both variables and check
	p.ExpectOK().Run("build")
	p.ExpectGolden().Run("variable", "ls")
	p.ExpectGolden().Run("variable", "ls", "--json")
	p.ExpectGolden().Run("variable", "ls", "--bash")
	p.ExpectGolden().Run("variable", "properties", "variable-abc")
	p.ExpectGolden().Run("variable", "layout", "variable-xyz")
	p.ExpectJsonArray("qId", "variable-abc", "variable-xyz").Run("variable", "ls", "--json")

	//Remove one variable and check
	p.ExpectOK().Run("variable", "rm", "variable-abc")
	p.ExpectJsonArray("qId", "variable-xyz").Run("variable", "ls", "--json")

	//Re-add the variable and check
	p.ExpectOK().Run("variable", "set", "test/projects/using-entities/variables.json")
	p.ExpectJsonArray("qId", "variable-xyz", "variable-abc").Run("variable", "ls", "--json")
}

func TestBookmarkManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}

	// Build with two bookmarks
	p.ExpectOK().Run("build")
	p.ExpectOK().Run("bookmark", "set", "test/projects/using-entities/bookmarks.json")
	p.ExpectOK().Run("bookmark", "ls") // Cannot ensure order of bookmarks so can't use golden.
	p.ExpectOK().Run("bookmark", "ls", "--json")
	p.ExpectOK().Run("bookmark", "ls", "--bash")
	p.ExpectOK().Run("bookmark", "properties", "alpha-bookmark-1")
	p.ExpectOK().Run("bookmark", "layout", "zeta-bookmark-2")
	p.ExpectIncludes("qId", "alpha-bookmark-1", "zeta-bookmark-2").Run("bookmark", "ls", "--json")

	// Reomve one bookmark and see that only the other one remains
	p.ExpectOK().Run("bookmark", "rm", "alpha-bookmark-1")
	p.ExpectJsonArray("qId", "zeta-bookmark-2").Run("bookmark", "ls", "--json")

	// Re-add the bookmarks and check
	p.ExpectOK().Run("bookmark", "set", "test/projects/using-entities/bookmarks.json")
	p.ExpectJsonArray("qId", "zeta-bookmark-2", "alpha-bookmark-1").Run("bookmark", "ls", "--json")
}

func TestOpeningWithoutData(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, Config: "test/projects/using-entities/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectOK().Run("build")

	// Open app without data and verify key printouts
	p.ExpectIncludes(`{"qSessionState":"SESSION_CREATED"}`, "without data").Run("connection", "ls", "--no-data", "--verbose")

	// Save objects in app opened without data
	p.ExpectIncludes("Saving objects in app...", "App successfully saved").Run("build", "--no-data")
}

func TestScriptManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Ttl: "0"}
	defer p.Reset()

	p.ExpectOK().Run("build", "--script", "test/projects/using-script/script1.qvs")

	// Set the script with zero TTL and --no-save This shouldn't persist the script1 qvs file
	p.ExpectGolden().Run("script", "set", "test/projects/using-script/script2.qvs", "--no-save")
	p.ExpectGolden().Run("script", "get")

	// Set the script without the --no-save-flag. This should persist the script1 qvs file
	p.ExpectGolden().Run("script", "set", "test/projects/using-script/script2.qvs")
	p.ExpectGolden().Run("script", "get")
}

func TestScriptVariables(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()
	// Build with script that creates two variables and check
	p.ExpectOK().Run("build", "--script", "test/projects/using-script/script3.qvs")
	p.ExpectJsonArray("title", "a", "b").Run("variable", "ls", "--json")
}

func TestTrafficFlag(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Ttl: "0"}
	defer p.Reset()

	p.ExpectOK().Run("build")
	p.ExpectGolden().Run("script", "get", "--traffic")

}

func TestUsingJwt(t *testing.T) {
	p := toolkit.Params{T: t, Ttl: "0"}
	p1 := p.WithParams(toolkit.Params{Engine: *toolkit.EngineStdIP})
	p2 := p.WithParams(toolkit.Params{Engine: *toolkit.EngineJwtIP})
	p3 := p.WithParams(toolkit.Params{Engine: *toolkit.EngineJwtIP, Headers: `authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0`})
	p4 := p.WithParams(toolkit.Params{Engine: *toolkit.EngineJwtIP, Config: "test/projects/using-jwts/corectl.yml"})

	p1.ExpectOK().ExpectIncludes("Connected without app to").Run("status")
	p2.ExpectErrorIncludes("headers", "authorization").Run("status")
	p3.ExpectOK().ExpectIncludes("Connected without app to").Run("status")
	p4.ExpectOK().ExpectIncludes("Connected without app to").Run("status")
}

func TestHelp(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP}
	p.ExpectGolden().Run("")
	p.ExpectGolden().Run("help")
	p.ExpectGolden().Run("help", "build")
}

func TestAppMissing(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP}
	p.ExpectErrorIncludes("no app specified").Run("connection", "ls")
}

func TestCatwalkUrl(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP}
	p.ExpectIncludes("Please provide a valid URL starting with 'https://', 'http://' or 'www'").Run("catwalk", "--catwalk-url=not-a-valid-url")
}

func TestWithoutApp(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP}

	p.ExpectOK().Run("status")
	p.ExpectOK().Run("version")
	// Can't run catwalk because it tries to open a browser when successful

	p.ExpectErrorIncludes("no app specified").Run("build")
	p.ExpectErrorIncludes("no app specified").Run("reload")

	p.ExpectErrorIncludes("no app specified").Run("assoc")
	p.ExpectErrorIncludes("no app specified").Run("eval", "count(a)")
	p.ExpectErrorIncludes("no app specified").Run("fields")
	p.ExpectErrorIncludes("no app specified").Run("keys")
	p.ExpectErrorIncludes("no app specified").Run("meta")
	p.ExpectErrorIncludes("no app specified").Run("tables")
	p.ExpectErrorIncludes("no app specified").Run("values", "foo")

	p.ExpectErrorIncludes("no app specified").Run("connection", "ls")
	p.ExpectErrorIncludes("no app specified").Run("dimension", "ls")
	p.ExpectErrorIncludes("no app specified").Run("measure", "ls")
	p.ExpectErrorIncludes("no app specified").Run("object", "ls")
	p.ExpectErrorIncludes("no app specified").Run("script", "get")
}

func TestWithUnknownApp(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name()}
	p.ExpectError().Run("reload")

	p.ExpectErrorIncludes("Could not find app").Run("assoc")
	p.ExpectErrorIncludes("Could not find app").Run("eval", "count(a)")
	p.ExpectErrorIncludes("Could not find app").Run("fields")
	p.ExpectErrorIncludes("Could not find app").Run("keys")
	p.ExpectErrorIncludes("Could not find app").Run("meta")
	p.ExpectErrorIncludes("Could not find app").Run("tables")
	p.ExpectErrorIncludes("Could not find app").Run("values", "foo")

	p.ExpectErrorIncludes("Could not find app").Run("connection", "ls")
	p.ExpectErrorIncludes("Could not find app").Run("dimension", "ls")
	p.ExpectErrorIncludes("Could not find app").Run("measure", "ls")
	p.ExpectErrorIncludes("Could not find app").Run("object", "ls")
	p.ExpectErrorIncludes("Could not find app").Run("script", "get")
}

func TestEvalOnUnknownAppl(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name()}
	p.ExpectIncludes("Could not find app: App not found (1003").Run("eval", "count(numbers)", "by", "xyz")
}

func TestEvalOnUnknownAppEngine(t *testing.T) {
	p := toolkit.Params{T: t, Engine: "localhost:9999", App: t.Name()}
	p.ExpectErrorIncludes("engine", "url").Run("eval", "count(numbers)", "by", "xyz")
}

func TestLicenseServiceDown(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineBadLicenseServerIP, App: t.Name()}
	p.ExpectIncludes("SESSION_ERROR_NO_LICENSE").Run("app", "ls")
}

func TestAppsInABAC(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineAbacIP, Config: "test/projects/abac/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectGolden().Run("build")
	p.ExpectIncludes("Connected to", "The data model has 1 table.").Run("status")
	p.ExpectJsonArray("name", t.Name()).Run("app", "ls", "--json")
	p.ExpectEqual(t.Name()).Run("app", "ls", "--bash")
	p.ExpectGolden().Run("meta")
}

func TestInvalidConfigs(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name()}
	pi1 := p.WithParams(toolkit.Params{Config: "test/projects/invalid-config/corectl-invalid.yml"})
	pi2 := p.WithParams(toolkit.Params{Config: "test/projects/invalid-config/corectl-invalid2.yml"})
	pi1.ExpectIncludes("'engin': did you mean 'engine'?",
		"'dimension': did you mean 'dimensions'?",
		"'verbos': did you mean 'verbose'?",
		"'apps': did you mean 'app'?",
		"'measure': did you mean 'measures'?",
		"'scrip': did you mean 'script'?",
		"'connection': did you mean 'connections'?",
		"'header': did you mean 'headers'?",
		"'trafic': did you mean 'traffic'?",
		"'object': did you mean 'objects'?",
	).Run("build")
	pi2.ExpectIncludes("'header': did you mean 'headers'?").Run("build")
}

func TestChildObjectsAndFullPropertyTree(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, Ttl: "0", App: t.Name()}
	defer p.Reset()

	// Build the app with an object with two children
	p.ExpectOK().Run("build", "--objects=test/projects/nested-objects/sheet.json")

	// List the objects and verify
	p.ExpectGolden().Run("object", "ls", "--json")

	// Remove the main object
	p.ExpectIncludes("Saving app...", "App successfully saved").Run("object", "rm", "a699ee97-152d-4470-9655-ae7c82d71491")

	// Verify that all three object are gone
	p.ExpectEqual("[]").Run("object", "ls", "--json")

}

func TestGetFullPropertyTree(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, Ttl: "0", App: t.Name()}
	defer p.Reset()

	// Build the app with an object with two children
	p.ExpectOK().Run("build", "--objects=test/projects/nested-objects/sheet.json")

	// List the objects and verify
	p.ExpectGolden().Run("object", "properties", "a699ee97-152d-4470-9655-ae7c82d71491", "--full")
}

func TestGetFullPropertyTreeMinimum(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, Ttl: "0", App: t.Name()}
	defer p.Reset()

	// Build the app with an object with two children
	p.ExpectOK().Run("build", "--objects=test/projects/nested-objects/sheet.json")

	// List the objects and verify
	p.ExpectGolden().Run("object", "properties", "a699ee97-152d-4470-9655-ae7c82d71491", "--full", "--minimum")
}

func TestConnectionDefinitionVariations(t *testing.T) {

	os.Setenv("CONN_TYPE", "folder")
	pNoConnections := toolkit.Params{T: t, Config: "test/projects/connections/corectl-no-connections.yml", Engine: *toolkit.EngineStdIP, App: t.Name() + "-1"}
	pCommandLine := toolkit.Params{T: t, Config: "test/projects/connections/corectl-no-connections.yml", Engine: *toolkit.EngineStdIP, App: t.Name() + "-2"}
	pWithConnections := toolkit.Params{T: t, Config: "test/projects/connections/corectl-with-connections.yml", Engine: *toolkit.EngineStdIP, App: t.Name() + "-3"}
	pConnectionsFile := toolkit.Params{T: t, Config: "test/projects/connections/corectl-connectionsref.yml", Engine: *toolkit.EngineStdIP, App: t.Name() + "-4"}
	pConnectionsFileEmpty := toolkit.Params{T: t, Config: "test/projects/connections/corectl-connectionsref-empty.yml", Engine: *toolkit.EngineStdIP, App: t.Name() + "-5"}
	defer pNoConnections.Reset() //This resets all apps since last reset

	//Build the apps
	pNoConnections.ExpectOK().Run("build")
	pCommandLine.ExpectOK().Run("build", "--connections=test/projects/connections/connections.yml")
	pWithConnections.ExpectOK().Run("build")
	pConnectionsFile.ExpectOK().Run("build")
	pConnectionsFileEmpty.ExpectOK().Run("build")

	pNoConnections.ExpectEmptyJsonArray().Run("connection", "ls", "--json")
	pCommandLine.ExpectJsonArray("qName", "testdata-separate-file").Run("connection", "ls", "--json")
	pWithConnections.ExpectJsonArray("qName", "testdata-inline").Run("connection", "ls", "--json")
	pConnectionsFile.ExpectJsonArray("qName", "testdata-separate-file").Run("connection", "ls", "--json")
}

// TestPrecedence checks that command line flags overrides config props
func TestCommandLineOverridingConfigFile(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/presedence/corectl.yml", Engine: *toolkit.EngineStdIP}
	defer p.Reset()
	p1 := p.WithParams(toolkit.Params{App: t.Name() + "-1"})
	p2 := p.WithParams(toolkit.Params{App: t.Name() + "-2"})
	p1.ExpectOK().Run("build")
	p1.ExpectJsonArray("qId", "my-hypercube").Run("object", "ls", "--json")
	p1.ExpectJsonArray("qId", "numbers-dimension").Run("dimension", "ls", "--json")
	p1.ExpectJsonArray("qId", "measure-1", "measure-2").Run("measure", "ls", "--json")
	p1.ExpectJsonArray("qName", "testdata").Run("connection", "ls", "--json")

	// Set objects, dimensions, measures and connection explicitly.
	// The information in the config should therefore be overriden.
	p2.ExpectOK().Run("build",
		"--config=test/projects/presedence/corectl.yml",
		"--objects=test/projects/presedence/o/*",
		"--dimensions=test/projects/presedence/d/*",
		"--measures=test/projects/presedence/m/*",
		"--connections=test/projects/presedence/connections.yml")
	p2.ExpectJsonArray("qId", "my-hypercube2").Run("object", "ls", "--json")
	p2.ExpectJsonArray("qId", "swedish-dimension").Run("dimension", "ls", "--json")
	p2.ExpectJsonArray("qId", "measure-x").Run("measure", "ls", "--json")
	p2.ExpectJsonArray("qName", "bogusname").Run("connection", "ls", "--json")
}

func TestImportApp(t *testing.T) {
	// Create tests for the standard, abac and jwt case.
	pStd := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP}
	pAbac := toolkit.Params{T: t, Engine: *toolkit.EngineAbacIP}
	pJwt := toolkit.Params{T: t, Engine: *toolkit.EngineJwtIP, Config: "test/projects/using-jwts/corectl.yml"}
	params := []toolkit.Params{pStd, pAbac, pJwt}
	for _, p := range params {
		// See if we can import the app test.qvf
		output := p.ExpectOK().Run("app", "import", "test/projects/import/test.qvf", "-q")
		// If it was created, we can remove it
		p.ExpectOK().Run("app", "rm", string(output), "--suppress")
		p.Reset()
	}
}

func TestUnbuild(t *testing.T) {
	os.Setenv("CONN_TYPE", "folder")
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	p2 := toolkit.Params{T: t, Config: "test/golden/unbuild/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name() + "-rebuild"}
	defer p.Reset()

	p.ExpectOK().Run("build")
	p.ExpectOK().Run("unbuild", "--dir", "test/golden/unbuild")
	p2.ExpectOK().Run("build")
	p2.ExpectGolden().Run("object", "ls")
	p2.ExpectGolden().Run("measure", "ls")
	p2.ExpectGolden().Run("dimension", "ls")
	p2.ExpectGolden().Run("meta")
	p2.ExpectIncludes("\"qTitle\":\"Test Unbuild App\",\"qThumbnail\":{\"qUrl\":\"/appcontent/qUrl.jpeg\"}").Run("app", "ls", "--traffic")
	os.RemoveAll("test/golden/unbuild")
}

func TestAddState(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()
	p.ExpectOK().Run("build")
	p.ExpectIncludes("Saving app...", "success").Run("state", "add", "MyTestState")
	p.ExpectGolden().Run("state", "ls")
	p.ExpectOK().Run("state", "rm", "MyTestState")
	p.ExpectError().Run("state", "rm", "MyTestState")
}

func TestCertificatesPath(t *testing.T) {
	relativePath := "test/projects/certificates/"
	pFlag := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Certificates: relativePath}
	pConfig := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Config: relativePath + "corectl-certificates.yml"}
	absolutePath, _ := filepath.Abs(relativePath)
	contextName := "cert-test"

	params := []toolkit.Params{pFlag, pConfig}
	for _, p := range params {
		defer p.Reset()
		p.ExpectOK().Run("context", "set", contextName)
		p.ExpectIncludes(absolutePath).Run("context", "get", contextName)
		p.ExpectOK().Run("context", "rm", contextName)
		params := []toolkit.Params{pFlag, pConfig}
		for _, p := range params {
			defer p.Reset()
			p.ExpectOK().Run("context", "set", contextName)
			p.ExpectIncludes(absolutePath).Run("context", "get", contextName)
			p.ExpectOK().Run("context", "rm", contextName)
		}
	}
}

func TestCertificatesPathNegative(t *testing.T) {
	pFlagNoCerts := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Certificates: "test/projects/"}
	pConfigNoCerts := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Config: "test/projects/certificates/corectl-certificates-invalid-path.yml"}
	pInvalidPath := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Certificates: "test/projects/non-existing-folder"}

	params := []toolkit.Params{pFlagNoCerts, pConfigNoCerts, pInvalidPath}
	for _, p := range params {
		defer p.Reset()
		p.ExpectErrorIncludes("could not load client certificate").Run("context", "set", "cert-test")
	}
}
