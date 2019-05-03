// +build integration

package test

import (
	"encoding/json"
	"github.com/qlik-oss/corectl/test/toolkit"
	"github.com/qlik-oss/enigma-go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestActualTestCase(t *testing.T) {

	p := toolkit.Params{T: t, Config: "test/project1/corectl.yml", App: t.Name()}
	defer p.Reset()

	p.ExpectIncludes("Connected",
		"TableA <<  5 Lines fetched",
		"TableB <<  5 Lines fetched",
		"Reload finished successfully",
		"Saving app... Done",
	).Run("build")

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
	p.ExpectGolden().Run("object", "properties", "my-hypercube", "--json")
	p.ExpectGolden().Run("measure", "ls", "--json")
	p.ExpectGolden().Run("dimension", "ls", "--json")
}

func TestActualTestCase2(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project1/corectl.yml"}
	defer p.Reset()
	p.ExpectOK().Run("build")
	p.ExpectGolden().Run("reload", "--silent")
	p.ExpectGolden().Run("reload", "--silent", "--nosave")
}

func TestActualTestCase3(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project1/corectl.yml"}
	defer p.Reset()
	envStd := p
	envAlt := p.WithParams(toolkit.Params{Config: "test/project1/corectl-alt.yml"})

	envStd.ExpectOK().Run("build")

	envStd.ExpectGolden().Run("measure", "set", "test/project1/not-following-glob-pattern-measure.json", "--no-save")
	envAlt.ExpectGolden().Run("measure", "ls", "--json")
	envAlt.ExpectGolden().Run("measure", "rm", "measure-3", "--no-save")
	envStd.ExpectGolden().Run("measure", "ls", "--json")

	envStd.ExpectGolden().Run("script", "set", "test/project1/dummy-script.qvs", "--no-save")
	envAlt.ExpectGolden().Run("script", "set", "test/project1/dummy-script.qvs")
	envAlt.ExpectGolden().Run("script", "set", "test/project1/dummy-script.qvs", "--traffic")

}

func TestOpeningWithoutData(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project1/corectl.yml", Ttl: "0"}
	defer p.Reset()
	envAlt := p.WithParams(toolkit.Params{Config: "test/project1/corectl-alt.yml"})

	// Open app without data
	envAlt.ExpectIncludes("without data").Run("connection", "ls", "--no-data", "--verbose")

	// Save objects in app opened without data
	p.ExpectIncludes("Saving objects in app... Done").Run("build", "--no-data")
}

func TestHelp(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project1/corectl.yml", Ttl: "0"}
	p.ExpectGolden().Run("")
	p.ExpectGolden().Run("help")
	p.ExpectGolden().Run("help", "build")
}

func TestNoData(t *testing.T) {
	//		// Verify behaviour when opening an app without data
	//		{"project 1 - open app without data", []string{"--config=test/project1/corectl-alt.yml", "--ttl", "0", connectToEngine}, []string{"connection", "ls", "--no-data", "--verbose"}, []string{"without data"}, initTest{true, true}},
	//		{"project 1 - save objects in app opened without data", []string{"--config=test/project1/corectl.yml", "--ttl", "0", connectToEngine, "--traffic=false", "--verbose=false"}, []string{"build", "--no-data"}, []string{"Saving objects in app... Done"}, initTest{false, true}},
}

func TestSeparateConnectionsFile(t *testing.T) {
	//		// Project 2 has separate connections file
	//		{"project 2 - build with connections", []string{connectToEngine, "-a=" + testAppName, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0"}, []string{"build", "--script=test/project2/script.qvs", "--connections=test/project2/connections.yml", "--objects=test/project2/object-*.json"}, []string{"datacsv << data 1 Lines fetched", "Reload finished successfully", "Saving app... Done"}, initTest{false, true}},
	//		{"project 2 - build with connections 2", []string{connectToEngine, "--config=test/project2/corectl-connectionsref.yml"}, []string{"build"}, []string{"datacsv << data 1 Lines fetched", "Reload finished successfully", "Saving app... Done"}, initTest{false, true}},
	//		{"project 2 - get fields ", []string{"--config=test/project2/corectl-alt.yml ", connectToEngine}, []string{"fields"}, []string{"golden", "project2-fields.golden"}, initTest{true, true}},
	//		{"project 2 - get data", []string{"--config=test/project2/corectl-alt.yml ", connectToEngine}, []string{"object", "data", "my-hypercube-on-commandline"}, []string{"golden", "project2-data.golden"}, initTest{true, true}},
	//		{"project 2 - keys", []string{"--config=test/project2/corectl-alt2.yml", connectToEngine}, []string{"keys"}, []string{"animal"}, initTest{true, true}},
}

func TestCatwalkUrl(t *testing.T) {
	//		{"err project 1 - invalid-catwalk-url", defaultConnectString1, []string{"catwalk", "--catwalk-url=not-a-valid-url"}, []string{"golden", "project1-catwalk-error.golden"}, initTest{false, false}},
	//		{"err 2", []string{connectToEngine, "--app=nosuchapp.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0"}, []string{"eval", "count(numbers)", "by", "xyz"}, []string{"golden", "err-2.golden"}, initTest{false, false}},
	//		{"err 3", []string{connectToEngine, "--app=" + testAppName, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0"}, []string{"object", "data", "nosuchobject"}, []string{"golden", "err-3.golden"}, initTest{true, true}},
}

func TestGetStatus(t *testing.T) {
	//		{"project 1 - get status", defaultConnectString1, []string{"status"}, []string{"Connected to " + testAppName + " @ ", "The data model has 2 tables."}, initTest{true, true}},
	//		{"list apps json", defaultConnectString1, []string{"app", "ls", "--json"}, []string{"\"id\": \"/apps/" + testAppName + "\","}, initTest{true, true}},
	//		{"err 1", []string{"--app=bogus", "--engine=localhost:9999"}, []string{"fields"}, []string{"Please check the --engine parameter or your config file", "Error details:  dial tcp"}, initTest{false, false}},
}

func TestMissingJWTHeader(t *testing.T) {
	//		// trying to connect to an engine that has JWT authorization activated without a JWT Header
	//		{"err jwt", []string{connectToEngine}, []string{"app", "ls"}, []string{"Error details:  401 from ws server: websocket: bad handshake"}, initTest{false, false}},
	//		{"err no license", []string{connectToEngineWithInccorectLicenseService}, []string{"app", "ls"}, []string{"Failed to connect to engine with error message:  SESSION_ERROR_NO_LICENSE"}, initTest{false, false}},
}

func TestABAC(t *testing.T) {
	//		// Verifying corectl against an engine running with ABAC enabled
	//		{"project 4 - get status", []string{"--config=test/project4/corectl.yml ", connectToEngineABAC}, []string{"status"}, []string{"Connected to " + testAppName + " @ ", "The data model has 1 table."}, initTest{true, true}},
	//		{"project 4 - list apps", []string{"--config=test/project4/corectl.yml ", connectToEngineABAC}, []string{"app", "ls", "--json"}, []string{"\"title\": \"" + testAppName + "\","}, initTest{true, true}},
	//		{"project 4 - get meta", []string{"--config=test/project4/corectl.yml ", connectToEngineABAC}, []string{"meta"}, []string{"golden", "project4-meta.golden"}, initTest{true, true}},
}

func TestInvalidConfig(t *testing.T) {
	//		// Verifying config validation
	//		{"err invalid 1", []string{"--config=test/project2/corectl-invalid.yml ", connectToEngine}, []string{"build"}, []string{"apps", "header", "object", "measure", "verbos", "trafic", "connection", "dimension"}, initTest{false, false}},
	//		{"err invalid 2", []string{"--config=test/project2/corectl-invalid2.yml ", connectToEngine}, []string{"build"}, []string{"'header': did you mean 'headers'?", "test/project2/corectl-invalid2.yml"}, initTest{false, false}},

}

func TestNestedObjectSupport(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project2/corectl-alt.yml", Ttl: "0", App: t.Name()}
	defer p.Reset()

	// Build the app with an object with two children
	p.ExpectOK().Run("build", "--objects=test/project2/sheet.json")

	// List the objects and verify
	p.ExpectGolden().Run("object", "ls", "--json")

	// Remove the main object
	p.ExpectIncludes("Saving app... Done").Run("object", "rm", "a699ee97-152d-4470-9655-ae7c82d71491")

	// Verify that all three object are gone
	p.ExpectEqual("[]").Run("object", "ls", "--json")

}

func TestConnections(t *testing.T) {
	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project2/corectl.yml", Ttl: "0", App: t.Name()}
	defer p.Reset()
	//setup env var for project 2
	os.Setenv("CONN_TYPE", "folder")

	//Build the app with a connection
	p.ExpectOK().Run("build", "--connections=test/project2/connections.yml")

	//Verify that the connecgtion is there
	output := p.ExpectOK().Run("connection", "ls", "--json")
	var connections []*enigma.Connection
	err := json.Unmarshal(output, &connections)
	assert.NoError(t, err)
	assert.NotNil(t, connections[0])
	assert.NotNil(t, connections[0].Id)

	//verify that removing the connection works
	p.ExpectOK().Run("connection", "rm", connections[0].Id)
	//verify that there is no connections in the app anymore.
	p.ExpectEqual("[]").Run("connection", "ls", "--json")

}

//
// TestPrecedence checks that command line flags overrides config props
func TestPrecedence(t *testing.T) {
	// Set objects, dimensions, measures and connection explicitly.
	// The information in the config should therefore be overriden.

	p := toolkit.Params{T: t, Engine: *toolkit.Engine1IP, Config: "test/project5/corectl.yml", Ttl: "0", App: t.Name()}
	defer p.Reset()

	output := p.ExpectOK().Run("build",
		"--config=test/project5/corectl.yml",
		"--objects=test/project5/o/*",
		"--dimensions=test/project5/d/*",
		"--measures=test/project5/m/*",
		"--connections=test/project5/connections.yml")

	var data []map[string]string
	entities := []string{"object", "dimension", "measure"}
	expected := []string{"my-hypercube2", "swedish-dimension", "measure-x"}
	for i, entity := range entities {
		output = p.ExpectOK().Run(entity, "ls", "--json")
		json.Unmarshal(output, &data)
		assert.Len(t, data, 1)
		assert.Equal(t, expected[i], data[0]["qId"])
	}

	var connections []*enigma.Connection
	output = p.ExpectOK().Run("connection", "ls", "--json")
	json.Unmarshal(output, &connections)
	assert.Len(t, connections, 1)
	assert.Equal(t, "bogusname", connections[0].Name)
}
