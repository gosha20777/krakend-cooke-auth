package gin

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	auth "github.com/gosha20777/krakend-cookie-auth"
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
		logger.Info("COOKIE: enabled for the endpoint", configuration.Endpoint)

		return func(c *gin.Context) {
			cookie, err := c.Request.Cookie(credentials.Cookie)
			if err != nil {
				logger.Warning("COOKIE: unable to get cookie", credentials.Cookie)
				c.String(http.StatusForbidden, "no auth header")
				return
			}

			info, err := validator.IsValid(cookie.Value)
			if err != nil {
				c.String(http.StatusUnauthorized, "wrong auth header")
				return
			}

			c.Request.Header.Set("X-Session-Id", strconv.Itoa(info.SessionId))
			c.Request.Header.Set("X-User-Id", info.UserId)
			next(c)
		}
	}
}
