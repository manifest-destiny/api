package apidocs

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

// Register sets up Swagger API documentation.
func Register(c *restful.Container, url string) {
	config := swagger.Config{
		ApiVersion:      "1",
		WebServices:     c.RegisteredWebServices(),
		WebServicesUrl:  url,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "../apidocs/swagger-ui/dist"}

	swagger.RegisterSwaggerService(config, c)
}
