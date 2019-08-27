package krakend

import (
    "github.com/devopsfaith/krakend/config"
    "github.com/devopsfaith/krakend/logging"
    "github.com/devopsfaith/krakend/proxy"

    krakendjose "github.com/devopsfaith/krakend-jose"
    ginjose "github.com/devopsfaith/krakend-jose/gin"
    ginkrakend "github.com/devopsfaith/krakend/router/gin"

    "net/http"
    "github.com/gin-gonic/gin"
    auth0 "github.com/auth0-community/go-auth0"
)

func NewJoseHandlerFactory(hf ginkrakend.HandlerFactory, logger logging.Logger, rf krakendjose.RejecterFactory) ginkrakend.HandlerFactory {
    hf = planValidatorFactory(hf, logger)
    // inject the plan validator just after the jose validator
    return ginjose.HandlerFactory(hf, logger, rf)
}

func planValidatorFactory(hf ginkrakend.HandlerFactory, logger logging.Logger) ginkrakend.HandlerFactory {
    return func(ecfg *config.EndpointConfig, prxy proxy.Proxy) gin.HandlerFunc {
        handler := hf(ecfg, prxy)

        cfg, ok := ecfg.ExtraConfig["github.com/openrm/iam/plans"]
        if !ok {
            return handler
        }

        var grade int = 0

        if cfg, ok := cfg.(map[string]interface{}); ok {
            if v, ok := cfg["grade"]; ok {
                switch v := v.(type) {
                case int:
                    grade = v
                case int64:
                    grade = int(v)
                case float64:
                    grade = int(v)
                }
            }
        }

        if grade == 0 {
            return handler
        }

        scfg, err := krakendjose.GetSignatureConfig(ecfg)
        if err != nil {
            return handler
        }

        extractor := auth0.FromMultiple(
            auth0.RequestTokenExtractorFunc(auth0.FromHeader),
            auth0.RequestTokenExtractorFunc(ginjose.FromCookie(scfg.CookieKey)),
        )

        logger.Info("IAM: plan validator enabled for the endpoint", ecfg.Endpoint)

        return func(c *gin.Context) {
            token, err := extractor.Extract(c.Request)
            if err != nil {
                c.AbortWithError(http.StatusUnauthorized, err)
                return
            }

            claims := map[string]interface{}{}
            // NOTICE it requires the validation is done before coming here
            if err := token.UnsafeClaimsWithoutVerification(&claims); err != nil {
                c.AbortWithError(http.StatusUnauthorized, err)
                return
            }

            if v, ok := claims["grade"]; ok {
                if claimGrade, ok := v.(int); ok {
                    if claimGrade >= grade {
                        handler(c)
                    }
                }
            }

            c.AbortWithStatus(http.StatusPaymentRequired)
        }
    }
}
