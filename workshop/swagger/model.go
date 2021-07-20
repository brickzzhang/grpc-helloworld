package swagger

const (
	defaultSwaggerAssetsUIPath = "./assets/swagger-ui/"
	// sourceIndexFile used to generate index.html file under assets directory
	sourceIndexFile = "workshop/swagger/index.html"
	// targetIndexFile index file name
	targetIndexFile   = "index.html"
	swFilePath        = "api"
	swSuffix          = ".swagger.json"
	swaggerConfigFile = "swagger-config.json"

	// SwJSONRoute route to json config file
	SwJSONRoute = "/swagger/"
	// SwaggerGatewayRoute route to business api
	SwaggerGatewayRoute = "/"
)

// SwWebRoute will be assigned when start swagger server
var SwWebRoute string

type swaggerURLs struct {
	URLs []*swaggerURL `json:"urls"`
}

type swaggerURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
