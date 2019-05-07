// build integration

package test

import (
	"github.com/qlik-oss/corectl/test/toolkit"
	"os"
	"testing"
)

func TestBasicAnalyzing(t *testing.T) {

	os.Setenv("CONN_TYPE", "folder")
	p := toolkit.Params{T: t, Config: "test/project1/corectl.yml", Engine: *toolkit.Engine0IP, App: t.Name() + "-ga.qvf"}
	defer p.Reset()

	p.ExpectGolden().Run("build")
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
	p.ExpectGolden().Run("object", "ls", "--json")
	p.ExpectGolden().Run("object", "data", "my-hypercube")
	p.ExpectGolden().Run("object", "properties", "my-hypercube")
	p.ExpectGolden().Run("object", "layout", "my-hypercube")
	p.ExpectGolden().Run("object", "data", "nosuchobject")
	p.ExpectGolden().Run("measure", "ls", "--json")
	p.ExpectGolden().Run("dimension", "ls", "--json")
}

func TestReload(t *testing.T) {
	p := toolkit.Params{T: t, Config: "test/project1/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectGolden().Run("build")
	p.ExpectGolden().Run("reload", "--silent")
	p.ExpectGolden().Run("reload", "--silent", "--nosave")
}

func TestMeasures(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project1/corectl.yml", App: t.Name()}
	defer p.Reset()
	envStd := p
	envAlt := p.WithParams(toolkit.Params{Config: "test/project1/corectl-alt.yml"})

	envStd.ExpectOK().Run("build")

	envStd.ExpectGolden().Run("measure", "set", "test/project1/not-following-glob-pattern-measure.json", "--no-save")
	envAlt.ExpectGolden().Run("measure", "ls", "--json")
	envAlt.ExpectGolden().Run("measure", "rm", "measure-3", "--no-save")
	envStd.ExpectGolden().Run("measure", "ls", "--json")
}

func TestScript(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project1/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectOK().Run("build")

	p.ExpectGolden().Run("script", "set", "test/project1/dummy-script.qvs", "--no-save")
	p.ExpectGolden().Run("script", "set", "test/project1/dummy-script.qvs")
	p.ExpectGolden().Run("script", "set", "test/project1/dummy-script.qvs", "--traffic")

}

func TestOpeningWithoutData(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project1/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.Run("build")

	// Open app without data and verify key printouts
	p.ExpectIncludes(`{"qSessionState":"SESSION_CREATED"}`, "without data").Run("connection", "ls", "--no-data", "--verbose")

	// Save objects in app opened without data
	p.ExpectIncludes("Saving objects in app... Done").Run("build", "--no-data")
}

func TestHelp(t *testing.T) {
	p := toolkit.Params{T: t}
	p.ExpectGolden().Run("")
	p.ExpectGolden().Run("help")
	p.ExpectGolden().Run("help", "build")
}

func TestNoData(t *testing.T) {
	p := toolkit.Params{T: t}
	p.ExpectGolden().Run("")
	p.ExpectGolden().Run("connection ls")
	p.ExpectGolden().Run("help", "build")
}

func TestCatwalkUrl(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine0IP, App: t.Name()}
	p.ExpectIncludes("Please provide a valid URL starting with 'https://', 'http://' or 'www'").Run("catwalk", "--catwalk-url=not-a-valid-url")
}

func TestEvalOnUnknownAppl(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine0IP, App: t.Name()}
	p.ExpectIncludes("Could not find app: App not found (1003)").Run("eval", "count(numbers)", "by", "xyz")
}

func TestEvalOnUnknownAppEngine(t *testing.T) {
	p := toolkit.Params{T: t, Engine: "localhost:9999", App: t.Name()}
	p.ExpectIncludes("Please check the --engine parameter or your config file.").Run("eval", "count(numbers)", "by", "xyz")
}

func TestMissingJWTHeader(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, App: t.Name()}
	p.ExpectIncludes("Please check the --engine parameter or your config file.").Run("app", "ls")
}

func TestLicenseServiceDown(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine3IP, App: t.Name()}
	p.ExpectIncludes("Failed to connect to engine with error message:  SESSION_ERROR_NO_LICENSE").Run("app", "ls")
}

func TestABAC(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine2IP, Config: "test/projects/abac/corectl.yml", App: t.Name()}
	defer p.Reset()
	p.ExpectGolden().Run("build")
	p.ExpectGolden().Run("status")
	p.ExpectJsonArray("name", "TestABAC").Run("app", "ls", "--json")
	p.ExpectGolden().Run("meta")
}

func TestInvalidConfig(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine0IP, App: t.Name()}
	pi1 := p.WithParams(toolkit.Params{Config: "test/projects/invalid-config/corectl-invalid.yml"})
	pi2 := p.WithParams(toolkit.Params{Config: "test/projects/invalid-config/corectl-invalid2.yml"})
	pi1.ExpectGolden().Run("build")
	pi2.ExpectGolden().Run("build")
}

func TestNestedObjects(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/projects/nested-objects/corectl-alt.yml", Ttl: "0", App: t.Name()}
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

func TestConnections(t *testing.T) {

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
func TestPrecedence2(t *testing.T) {
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
