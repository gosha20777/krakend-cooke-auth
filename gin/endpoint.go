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
			cookie, err := c.Request.Cookie("mojo")
			if err != nil {
				logger.Warning("COOKE: ", "unable to get cookie")
				c.String(http.StatusForbidden, "wrong auth header")
				return
			}

			logger.Info("COOKE: get cookie", cookie.Value)
			info, err := validator.IsValid(cookie.Value)
			if err != nil {
				c.String(http.StatusForbidden, "wrong auth header")
				return
			}

			logger.Info("COOKE: session_id", info.SessionId)
			logger.Info("COOKE: user_id", info.UserId)
			c.Request.Header.Set("X-SessionId", info.SessionId)

			next(c)
		}
	}
}
