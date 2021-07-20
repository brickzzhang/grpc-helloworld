package swagger

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

func init() {
	// swagger error doesn't matter
	_ = CreateSwaggerIndex()
	_ = CreateSwaggerConfigJSONFile()
}

// WebJSONHandler used by
// http://localhost:8082/swagger/swagger/application/v1/application_service.swagger.json route
func WebJSONHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, swSuffix) {
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, SwJSONRoute)
	p = path.Join(swFilePath, p)

	http.ServeFile(w, r, p)
}

// WebHandler used by http://localhost:8082/sw/ route
func WebHandler() http.Handler {
	return http.FileServer(http.Dir(defaultSwaggerAssetsUIPath))
}

//Forward2GprcGatewayHandler used by http://localhost:8082/ route
func Forward2GprcGatewayHandler(grpcGwPort string) http.Handler {
	var rp *httputil.ReverseProxy
	u, _ := url.Parse(fmt.Sprintf("http://localhost:%s", grpcGwPort))
	rp = httputil.NewSingleHostReverseProxy(u)

	return rp
}
