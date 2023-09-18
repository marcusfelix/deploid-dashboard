package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"main/includes"
	"net/http"
	"os"

	_ "main/migrations"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/otiai10/opengraph"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v74"
	"mvdan.cc/xurls/v2"
)

//go:embed app/dist/*
var static embed.FS

func main() {

	// Getting env vars
	err := godotenv.Load(`.env.local`)
	if err != nil {
		fmt.Print("Local .env not found, using environment variables.")
	}

	// Initialize Stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Initialize Pocketbase
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		// Serving embed static files
		subFs := echo.MustSubFS(static, "app/dist")
		e.Router.GET("/*", apis.StaticDirectoryHandler(subFs, true))

		// Getting secrets
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/secrets",
			Handler: func(c echo.Context) error {

				// Getting secrets from Github
				secrets, err := includes.GetSecretsFromRepo(c.PathParam("project"))
				if err != nil {
					return apis.NewBadRequestError("error while getting secrets", err)
				}

				// Returning secrets
				return c.JSON(200, secrets)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Updating secret
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/projects/:project/secrets",
			Handler: func(c echo.Context) error {

				// Creating or updating secret
				er, key := includes.CreateOrUpdateRepoSecret(c)
				if er != nil {
					return apis.NewBadRequestError("error while updating secrets", err)
				}

				return c.JSON(200, map[string]interface{}{key: "***"})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Getting contents
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/projects/:project/read",
			Handler: func(c echo.Context) error {

				// Getting contents
				data, err := includes.GetPrivateRepoContent(c, app)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Getting releases
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/releases",
			Handler: func(c echo.Context) error {
				releases, err := includes.GetRepoReleases(c.PathParam("project"))
				if err != nil {
					return apis.NewBadRequestError("error while getting releases", err)
				}
				return c.JSON(200, releases)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
			},
		})

		// Getting release asset
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/releases/:asset",
			Handler: func(c echo.Context) error {
				releases, err := includes.GetRepoReleaseAsset(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting releases", err)
				}
				return c.JSON(200, releases)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
			},
		})

		// Getting repo file
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/projects/:project/read",
			Handler: func(c echo.Context) error {

				// Getting contents
				data, err := includes.GetPrivateRepoContent(c, app)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Writing repo file
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/projects/:project/write",
			Handler: func(c echo.Context) error {

				// Getting contents
				data, err := includes.UpdatePrivateRepoContent(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Getting template readme file
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/templates/:slug",
			Handler: func(c echo.Context) error {

				// Getting contents
				data, err := includes.GetRepoContent(c.PathParam("slug"), "readme.md")
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Getting repo issues
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/issues",
			Handler: func(c echo.Context) error {

				// Getting contents
				data, err := includes.GetRepoIssues(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Create repo issue
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/projects/:project/issues",
			Handler: func(c echo.Context) error {

				// Getting contents
				issue, _, err := includes.CreateRepoIssue(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, issue)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Getting repo comments
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/issues/:issue/comments",
			Handler: func(c echo.Context) error {

				// Getting contents
				data, err := includes.GetCommentsFromIssue(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Create repo issue comment
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/projects/:project/issues/:issue/comments",
			Handler: func(c echo.Context) error {

				// Getting contents
				issue, _, err := includes.CreateRepoIssueComment(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, issue)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Change issue state
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/projects/:project/issues/:issue/state",
			Handler: func(c echo.Context) error {

				// Getting contents
				issue, _, err := includes.EditRepoIssueState(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, issue)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Get repo branches
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/branches",
			Handler: func(c echo.Context) error {

				// Getting contents
				branches, err := includes.GetRepoBranches(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, branches)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Get repo commits
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/commits",
			Handler: func(c echo.Context) error {

				// Getting contents
				commits, err := includes.GetRepoCommits(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, commits)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Get repo workflows
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/projects/:project/workflow/:branch",
			Handler: func(c echo.Context) error {

				// Getting contents
				commits, err := includes.GetRepoWorkfows(c)
				if err != nil {
					return apis.NewBadRequestError("error while getting content", err)
				}

				return c.JSON(200, commits)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Getting prices with products
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/prices",
			Handler: func(c echo.Context) error {

				// Getting prices and products
				prices := includes.GetPricesWithProducts()
				return c.JSON(200, prices.PriceList())
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Checkout
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/checkout",
			Handler: func(c echo.Context) error {

				// Creating checkout session
				err, url := includes.CreateStripeCheckout(c, app)
				if err != nil {
					return apis.NewBadRequestError("error while creating checkout session", err)
				}

				// Returning url
				return c.JSON(200, map[string]string{
					"url": url,
				})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		// Getting customer portal
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/management/:project",
			Handler: func(c echo.Context) error {

				// Getting customer portal
				url, err := includes.GetStripeCustomerPortal(c, app)
				if err != nil {
					return apis.NewBadRequestError("error while creating customer portal", err)
				}

				// Returning url
				return c.JSON(200, map[string]string{
					"url": url,
				})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
			},
		})

		// Stripe webhook
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/events",
			Handler: func(c echo.Context) error {
				includes.HandleStripeEvent(c, app)

				return c.JSON(200, map[string]interface{}{})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		// Github webhook
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/updates",
			Handler: func(c echo.Context) error {

				err := includes.HandleGithubUpdate(c, app)
				if err != nil {
					return apis.NewBadRequestError("error while handling update", err)
				}

				return c.JSON(200, map[string]interface{}{})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		return nil
	})

	// Fetching Open Graph data from message body urls
	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "messages" {
			output := xurls.Strict().FindAllString(e.Record.GetString("body"), -1)
			var metadata []any

			for _, url := range output {
				graph, err := opengraph.Fetch(url)
				if err != nil {
					log.Fatal(err)
					return nil
				}
				graph.ToAbsURL()

				image := ""

				if len(graph.Image) > 0 {
					image = graph.Image[0].URL
				}

				json := map[string]string{
					"url":         url,
					"image":       image,
					"title":       graph.Title,
					"description": graph.Description,
				}
				metadata = append(metadata, json)
			}
			out, err := json.Marshal(metadata)
			if err != nil {
				log.Fatal(err)
				return nil
			}

			write := string(out)
			e.Record.Set("metadata", write)
		}
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
