package loopia

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kolo/xmlrpc"
	"github.com/stretchr/testify/assert"
	"github.com/subchen/go-xmldom"
)

var (
	handlers map[string]methodHandler
)

type methodHandler func(t *testing.T, w http.ResponseWriter, params []string)

func init() {
	handlers = make(map[string]methodHandler)
	handlers["getZoneRecords"] = getZoneRecordsHandler
	handlers["getSubdomains"] = getSubdomainsHandler
	handlers["addSubdomain"] = addSubdomainHandler
	handlers["addZoneRecord"] = addZoneRecordHandler
	handlers["updateZoneRecord"] = updateZoneRecordHandler
	handlers["removeZoneRecord"] = returnOkHandler
	handlers["removeSubdomain"] = returnOkHandler
}

type testContext struct {
	mux *http.ServeMux

	rpc    *xmlrpc.Client
	server *httptest.Server
}

func (tc *testContext) getProvider() *Provider {
	p := &Provider{}
	p.rpc = tc.rpc
	return p
}

func setupTest(t *testing.T) *testContext {
	tc := &testContext{}
	tc.mux = http.NewServeMux()
	tc.server = httptest.NewServer(tc.mux)
	tc.rpc, _ = xmlrpc.NewClient(tc.server.URL, nil)
	tc.mux.HandleFunc("/", apiHandler(t))
	return tc
}

func teardownTest(tc *testContext) {
	if tc.server != nil {
		tc.server.Close()
	}
}

func apiHandler(t *testing.T) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "Error reading request body")
		strBody := string(body)
		doc := xmldom.Must(xmldom.ParseXML(strBody))
		root := doc.Root

		method := root.GetChild("methodName").Text
		params := root.GetChild("params")
		values := params.Query("//value")

		strValues := []string{}
		for _, v := range values {
			strValues = append(strValues, v.FirstChild().Text)
		}

		h := handlers[method]
		if h != nil {
			h(t, w, strValues)
			return
		}
		t.Errorf("method %s not implemented", method)
		// byteArray, _ := ioutil.ReadFile("testdata/error.xml")
		// fmt.Fprint(w, string(byteArray[:]))
	}
}

func getSubdomainsHandler(t *testing.T, w http.ResponseWriter, params []string) {
	byteArray, _ := os.ReadFile("testdata/subdomains.xml")
	fmt.Fprint(w, string(byteArray[:]))
}

func getZoneRecordsHandler(t *testing.T, w http.ResponseWriter, params []string) {

	recname := params[len(params)-1] //last parameter
	if recname == "*" {
		recname = ""
	}
	filename := fmt.Sprintf("testdata/zone_records_%s.xml", recname)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		filename = "testdata/empty_list.xml"
	}

	byteArray, _ := os.ReadFile(filename)
	fmt.Fprint(w, string(byteArray[:]))
}

func addSubdomainHandler(t *testing.T, w http.ResponseWriter, params []string) {
	fmt.Printf(" > addSubdomainHandler(%s, %s)\n", params[3], params[4])
	assert.Len(t, params, 5)
	lastp := params[len(params)-1]
	assert.GreaterOrEqual(t, len(lastp), 1)
	byteArray, _ := os.ReadFile("testdata/ok.xml")
	fmt.Fprint(w, string(byteArray[:]))
}

func addZoneRecordHandler(t *testing.T, w http.ResponseWriter, params []string) {
	byteArray, _ := os.ReadFile("testdata/ok.xml")
	fmt.Fprint(w, string(byteArray[:]))
}

func updateZoneRecordHandler(t *testing.T, w http.ResponseWriter, params []string) {
	byteArray, _ := ioutil.ReadFile("testdata/ok.xml")
	fmt.Fprint(w, string(byteArray[:]))
}

func returnOkHandler(t *testing.T, w http.ResponseWriter, params []string) {
	byteArray, _ := os.ReadFile("testdata/ok.xml")
	fmt.Fprint(w, string(byteArray[:]))
}
