package apidocs

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func Register(c *restful.Container, port string, tls bool) {
	var scheme string
	if tls {
		scheme = "https"
	} else {
		scheme = "http"
	}
	config := swagger.Config{
		ApiVersion:      "1",
		WebServices:     c.RegisteredWebServices(),
		WebServicesUrl:  scheme + "://localhost:" + port,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "apidocs/swagger-ui/dist"}

	swagger.RegisterSwaggerService(config, c)
}
