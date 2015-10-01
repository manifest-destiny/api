package apidocs

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func Register(c *restful.Container, port string) {
	config := swagger.Config{
		ApiVersion:      "1",
		WebServices:     c.RegisteredWebServices(),
		WebServicesUrl:  "http://localhost:" + port,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "apidocs/swagger-ui/dist"}

	swagger.RegisterSwaggerService(config, c)
}
