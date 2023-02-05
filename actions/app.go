package actions

import (
	"errors"
	"net/http"
	"sync"

	"paphos/locales"
	"paphos/models"
	"paphos/shared"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n/v2"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
	"github.com/unrolled/secure"

	_ "paphos/docs"

	buffaloSwagger "github.com/swaggo/buffalo-swagger"
	"github.com/swaggo/buffalo-swagger/swaggerFiles"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.

// @title          Paphos API
// @version        1.0
// @description    Base backend API to serve Pygamillion UI
// @termsOfService http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host                     tobedefined.com
// @BasePath                  /api/v1
// @securityDefinitions.bearer BearerAuth
func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_paphos_session",
		})

		// By default, Buffallo returns HTTP 500 with text "EOF" when a client sends
		// a bad payload. That's not _our_ fault though, so we'll replace that with
		// a 400 so our logs actually reflect that there was nothing wrong on our
		// side.
		var originalError500Handler = app.ErrorHandlers[500]
		app.ErrorHandlers[http.StatusInternalServerError] = func(status int, err error, c buffalo.Context) error {
			if err.Error() == "EOF" {
				status = http.StatusBadRequest
				err = c.Error(http.StatusBadRequest, errors.New("malformed request payload"))
			}

			return originalError500Handler(status, err, c)
		}
		// TODO(11b): related to the above ^
		// Marshalling failures result in a 500 even though they're the user's fault
		// most of the time (e.g. sent in an int instead of a string), fix that.

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		apiV1Group := app.Group("/api/v1")
		userGroup := apiV1Group.Group("/users")
		userGroup.POST("/register", UsersRegisterPost)
		userGroup.POST("/login", UsersLoginPost)

		apiV1Group.Use(shared.ExtractDataFromJWTMiddleware)

		apiV1Group.Resource("/characters", CharactersResource{})

		app.GET("/swagger/{doc:.*}", buffaloSwagger.WrapHandler(swaggerFiles.Handler))
		app.GET("/", func(c buffalo.Context) error {
			return c.Redirect(301, "/swagger/index.html") // redirect to swagger route
		})

		// Disabled for now since we don't need this in the front-end yet and it
		// leaks user emails.
		// userGroup.GET("/{user_id}", UsersShowGet)
	})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
