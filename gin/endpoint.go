package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	auth "github.com/gosha20777/krakend-cooke-auth"
)

// HandlerFactory decorates a krakendgin.HandlerFactory with the auth layer
func HandlerFactory(hf krakendgin.HandlerFactory, logger logging.Logger) krakendgin.HandlerFactory {
	return func(configuration *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
		next := hf(configuration, proxy)

		v := auth.ConfigGetter(configuration.ExtraConfig)
		if v == nil {
			return next
		}
		credentials, ok := v.(auth.Credentials)
		if !ok {
			return next
		}

		validator := auth.NewCredentialsValidator(credentials)

		return func(c *gin.Context) {
			logger.Info("COOKE: ", "a")
			if !validator.IsValid("a") {
				c.String(http.StatusForbidden, "wrong auth header")
				return
			}
			next(c)
		}
	}
}
