// +build integration

package test

import (
	"github.com/qlik-oss/corectl/test/toolkit"
	"os"
	"testing"
)

func TestBasicAnalyzing(t *testing.T) {

	os.Setenv("CONN_TYPE", "folder")
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()

	p.ExpectGolden().Run("build")

	p.ExpectGolden().Run("reload")
	p.ExpectGolden().Run("reload", "--silent")
	p.ExpectGolden().Run("reload", "--silent", "--no-save")

	p.ExpectGolden().Run("status")
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
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectOK().Run("build")
	p.ExpectGolden().Run("reload", "--silent")
	p.ExpectGolden().Run("reload", "--silent", "--no-save")
}

func TestConnectionManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/using-entities/corectl.yml", Engine: *toolkit.EngineStdIP, App: t.Name()}
	defer p.Reset()
	p.Run("build")
	p.ExpectIncludes(`myconnection | testconnector`).Run("connection", "ls")
	p.ExpectIncludes(`"qConnectionString": "CUSTOM CONNECT TO \"provider=testconnector;host=corectl-test-connector;\""`).Run("connection", "ls", "--json")
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
	p.ExpectError("Invalid handle: Invalid Params (-32602)").Run("object", "data", "nosuchobject")

	p.ExpectJsonArray("qId", "object-using-dims-and-measures", "object-using-inline").Run("object", "ls", "--json")

	// Remove one object and check
	p.ExpectOK().Run("object", "rm", "object-using-inline")
	p.ExpectJsonArray("qId", "object-using-dims-and-measures").Run("object", "ls", "--json")

	// Re-add the object and check
	p.ExpectOK().Run("object", "set", "test/projects/using-entities/object-using-inline.json")
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
	p.ExpectOK().Run("dimension", "set", "test/projects/using-entities/dimension-abcs.json")
	p.ExpectJsonArray("qId", "dimension-xyz", "dimension-abcs").Run("dimension", "ls", "--json")
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
	p.ExpectOK().Run("measure", "set", "test/projects/using-entities/measure-count-numbers.json")
	p.ExpectJsonArray("qId", "measure-sum-numbers", "measure-count-numbers").Run("measure", "ls", "--json")

}

func TestOpeningWithoutData(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, Config: "test/projects/using-entities/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectOK().Run("build")

	// Open app without data and verify key printouts
	p.ExpectIncludes(`{"qSessionState":"SESSION_CREATED"}`, "without data").Run("connection", "ls", "--no-data", "--verbose")

	// Save objects in app opened without data
	p.ExpectIncludes("Saving objects in app... Done").Run("build", "--no-data")
}

func TestScriptManagementCommands(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name(), Ttl: "0"}
	defer p.Reset()

	p.ExpectOK().Run("build")
	// Set the script with zero TTL and --no-save This shouldn't persist the script1 qvs file
	p.ExpectGolden().Run("script", "set", "test/projects/using-script/script1.qvs", "--no-save")
	p.ExpectGolden().Run("script", "get")
	// Set the script without the --no-save-flag. This should persist the script1 qvs file
	p.ExpectGolden().Run("script", "set", "test/projects/using-script/script1.qvs")
	p.ExpectGolden().Run("script", "get")

	// Change it to see that that works
	p.ExpectGolden().Run("script", "set", "test/projects/using-script/script2.qvs")
	p.ExpectGolden().Run("script", "get")
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
	p2.ExpectErrorIncludes("Please check the --engine parameter or your config file.").Run("status")
	p3.ExpectOK().ExpectIncludes("Connected without app to").Run("status")
	p4.ExpectOK().ExpectIncludes("Connected without app to").Run("status")
}

func TestHelp(t *testing.T) {
	p := toolkit.Params{T: t}
	p.ExpectGolden().Run("")
	p.ExpectGolden().Run("help")
	p.ExpectGolden().Run("help", "build")
}

func TestAppMissing(t *testing.T) {
	p := toolkit.Params{T: t}
	p.ExpectError("Error: No app specified").Run("connection", "ls")
}

func TestCatwalkUrl(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name()}
	p.ExpectIncludes("Please provide a valid URL starting with 'https://', 'http://' or 'www'").Run("catwalk", "--catwalk-url=not-a-valid-url")
}

func TestEvalOnUnknownAppl(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineStdIP, App: t.Name()}
	p.ExpectIncludes("Could not find app: App not found (1003)").Run("eval", "count(numbers)", "by", "xyz")
}

func TestEvalOnUnknownAppEngine(t *testing.T) {
	p := toolkit.Params{T: t, Engine: "localhost:9999", App: t.Name()}
	p.ExpectIncludes("Please check the --engine parameter or your config file.").Run("eval", "count(numbers)", "by", "xyz")
}

func TestLicenseServiceDown(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineBadLicenseServerIP, App: t.Name()}
	p.ExpectIncludes("Failed to connect to engine with error message:  SESSION_ERROR_NO_LICENSE").Run("app", "ls")
}

func TestAppsInABAC(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.EngineAbacIP, Config: "test/projects/abac/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectGolden().Run("build")
	p.ExpectGolden().Run("status")
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
	p.ExpectIncludes("Saving app... Done").Run("object", "rm", "a699ee97-152d-4470-9655-ae7c82d71491")

	// Verify that all three object are gone
	p.ExpectEqual("[]").Run("object", "ls", "--json")

}

func TestConnectionDefinitionVariations(t *testing.T) {

	os.Setenv("CONN_TYPE", "folder")
	pNoConnections := toolkit.Params{T: t, Config: "test/projects/connections/corectl-no-connections.yml", App: t.Name() + "-1"}
	pCommandLine := toolkit.Params{T: t, Config: "test/projects/connections/corectl-no-connections.yml", App: t.Name() + "-2"}
	pWithConnections := toolkit.Params{T: t, Config: "test/projects/connections/corectl-with-connections.yml", App: t.Name() + "-3"}
	pConnectionsFile := toolkit.Params{T: t, Config: "test/projects/connections/corectl-connectionsref.yml", App: t.Name() + "-4"}
	defer pNoConnections.Reset() //This resets all apps since last reset

	//Build the apps
	pNoConnections.ExpectOK().Run("build")
	pCommandLine.ExpectOK().Run("build", "--connections=test/projects/connections/connections.yml")
	pWithConnections.ExpectOK().Run("build")
	pConnectionsFile.ExpectOK().Run("build")

	pNoConnections.ExpectEmptyJsonArray().Run("connection", "ls", "--json")
	pCommandLine.ExpectJsonArray("qName", "testdata-separate-file").Run("connection", "ls", "--json")
	pWithConnections.ExpectJsonArray("qName", "testdata-inline").Run("connection", "ls", "--json")
	pConnectionsFile.ExpectJsonArray("qName", "testdata-separate-file").Run("connection", "ls", "--json")
}

// TestPrecedence checks that command line flags overrides config props
func TestCommandLineOverridingConfigFile(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/projects/presedence/corectl.yml"}
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
