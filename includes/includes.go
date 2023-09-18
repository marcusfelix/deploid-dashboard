package includes

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/hmac"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"

	"github.com/dchest/uniuri"
	"github.com/google/go-github/v52/github"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/russross/blackfriday"

	"github.com/stripe/stripe-go/v74"
	portal "github.com/stripe/stripe-go/v74/billingportal/session"
	session "github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
	"github.com/stripe/stripe-go/v74/webhook"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
)

func GetSecretsFromRepo(name string) ([]any, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Get Repo Secrets
	opts := &github.ListOptions{}
	secrets, _, err := client.Actions.ListRepoSecrets(ctx, "deploidstudio", name, opts)
	if err != nil {
		return nil, err
	}

	// Return Secrets
	list := []any{}
	for _, secret := range secrets.Secrets {
		data := map[string]interface{}{
			"name":    secret.Name,
			"created": time.Unix(secret.CreatedAt.Unix(), 0),
			"updated": time.Unix(secret.UpdatedAt.Unix(), 0),
		}

		list = append(list, data)
	}

	return list, nil
}

func CreateOrUpdateRepoSecret(c echo.Context) (error, string) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	name := c.PathParam("project")
	jsondata := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return apis.NewBadRequestError("error while updating secrets", err), ""
	}

	// Updating secret
	key := jsondata["key"]
	value := jsondata["value"]

	// Get Repo Public Key
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, "deploidstudio", name)
	if err != nil {
		return err, ""
	}

	// Encrypt Secret
	encryptedSecret, err := EncryptSecretWithPublicKey(publicKey, key, value)
	if err != nil {
		return err, ""
	}

	// Create or Update Repo Secret
	if _, err := client.Actions.CreateOrUpdateRepoSecret(ctx, "deploidstudio", name, encryptedSecret); err != nil {
		return fmt.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err), ""
	}

	// We need nothing from this function
	return nil, key
}

func GetPrivateRepoContent(c echo.Context, app *pocketbase.PocketBase) (map[string]interface{}, error) {

	// Getting params
	id := c.PathParam("project")
	jsondata := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return nil, err
	}
	path := jsondata["path"]

	// Check if user is allowed
	project, err := CheckProjectPermission(app, c, id)
	if err != nil {
		return nil, err
	}

	// Get content
	name := project.GetString("slug")
	data, err := GetRepoContent(name, path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UpdatePrivateRepoContent(c echo.Context) (map[string]interface{}, error) {

	// Getting params
	name := c.PathParam("project")
	jsondata := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return nil, err
	}

	path := jsondata["path"]
	content := jsondata["content"]
	sha := jsondata["sha"]
	message := jsondata["message"]

	// Get Github Auth Token
	ctx, client := GithubAuth()
	decoded, err := b64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	// Update content
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(decoded),
		SHA:     github.String(sha),
	}
	response, _, err := client.Repositories.UpdateFile(ctx, "deploidstudio", name, path, opts)
	if err != nil {
		return nil, err
	}

	// Return response
	data := map[string]interface{}{
		"commit": response.Commit,
	}

	// Return data
	return data, nil
}

func GetRepoContent(repo string, path string) (map[string]interface{}, error) {
	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Get content
	opts := &github.RepositoryContentGetOptions{}
	contents, _, _, err := client.Repositories.GetContents(ctx, "deploidstudio", repo, path, opts)
	if err != nil {
		return nil, err
	}

	// Return content
	data := map[string]interface{}{
		"name":    contents.Name,
		"content": contents.Content,
		"sha":     contents.SHA,
	}

	return data, nil
}

func GetRepoReleases(name string) ([]any, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Get Releases
	opts := &github.ListOptions{}
	releases, _, err := client.Repositories.ListReleases(ctx, "deploidstudio", name, opts)
	if err != nil {
		return nil, err
	}
	list := []any{}
	for _, release := range releases {
		assets := []any{}
		for _, asset := range release.Assets {
			adata := map[string]interface{}{
				"name": asset.Name,
				"id":   asset.ID,
				"url":  asset.BrowserDownloadURL,
			}
			assets = append(assets, adata)
		}
		data := map[string]interface{}{
			"name":    release.Name,
			"body":    release.Body,
			"tag":     release.TagName,
			"created": time.Unix(release.CreatedAt.Unix(), 0),
			"author": map[string]interface{}{
				"avatar": release.Author.AvatarURL,
				"name":   release.Author.Login,
			},
			"assets": assets,
		}
		list = append(list, data)
	}

	// Return list
	return list, nil
}

func GetRepoReleaseAsset(c echo.Context) (string, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	name := c.PathParam("project")
	id := c.PathParam("asset")

	// Parse Id to int64
	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return "", err
	}

	// Get content
	_, redirect, err := client.Repositories.DownloadReleaseAsset(ctx, "deploidstudio", name, aid, nil)
	if err != nil {
		return "", err
	}

	// Return content
	return redirect, nil
}

func GetRepoIssues(c echo.Context) ([]any, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	name := c.PathParam("project")

	// Get Issues
	opts := &github.IssueListByRepoOptions{State: "all"}
	issues, _, err := client.Issues.ListByRepo(ctx, "deploidstudio", name, opts)
	if err != nil {
		return nil, err
	}
	list := []any{}
	for _, issue := range issues {
		labels := []*string{}
		for _, label := range issue.Labels {
			labels = append(labels, label.Name)
		}
		data := map[string]interface{}{
			"issue":    issue.Number,
			"title":    issue.Title,
			"body":     issue.Body,
			"state":    issue.State,
			"labels":   labels,
			"comments": issue.Comments,
			"user": map[string]interface{}{
				"login": issue.User.Login,
				"name":  issue.User.Name,
				"photo": issue.User.AvatarURL,
			},
			"created": time.Unix(issue.CreatedAt.Unix(), 0),
		}

		list = append(list, data)
	}

	// Return data
	return list, nil
}

func CreateRepoIssue(c echo.Context) (*github.Issue, *github.Response, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	slug := c.PathParam("project")
	jsondata := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return nil, nil, err
	}

	title := jsondata["title"]
	body := jsondata["body"]

	// Get user
	user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

	// Including shadow user
	body = body + "\n\n" + "－ " + user.GetString("name") + " <" + user.GetString("email") + ">"

	// Create Issue
	issue := &github.IssueRequest{
		Title: github.String(title),
		Body:  github.String(body),
	}

	// Create Issue
	return client.Issues.Create(ctx, "deploidstudio", slug, issue)
}

func GetCommentsFromIssue(c echo.Context) ([]any, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	name := c.PathParam("project")
	issue := c.PathParam("issue")

	// Convert string to int
	iid, err := strconv.Atoi(issue)
	if err != nil {
		return nil, err
	}

	// Get comments from issue
	opts := &github.IssueListCommentsOptions{}
	comments, _, err := client.Issues.ListComments(ctx, "deploidstudio", name, iid, opts)
	if err != nil {
		return nil, err
	}
	list := []any{}
	for _, comment := range comments {
		data := map[string]interface{}{
			"body": comment.Body,
			"user": map[string]interface{}{
				"login": comment.User.Login,
				"name":  comment.User.Name,
				"photo": comment.User.AvatarURL,
			},

			"created": time.Unix(comment.CreatedAt.Unix(), 0),
		}

		list = append(list, data)
	}

	// Return data
	return list, nil
}

func CreateRepoIssueComment(c echo.Context) (*github.IssueComment, *github.Response, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	slug := c.PathParam("project")
	issue := c.PathParam("issue")

	// Convert string to int
	iid, err := strconv.Atoi(issue)
	if err != nil {
		return nil, nil, err
	}

	// Getting comment data
	jsondata := make(map[string]string)
	err = json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return nil, nil, err
	}
	body := jsondata["body"]

	// Get user
	user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

	// Including shadow user
	body = body + "\n\n" + "－ " + user.GetString("name") + " <" + user.GetString("email") + ">"

	// Create Issue
	comment := &github.IssueComment{
		Body: github.String(body),
	}
	return client.Issues.CreateComment(ctx, "deploidstudio", slug, iid, comment)
}

func EditRepoIssueState(c echo.Context) (*github.Issue, *github.Response, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	slug := c.PathParam("project")
	issue := c.PathParam("issue")

	// Convert string to int
	iid, err := strconv.Atoi(issue)
	if err != nil {
		return nil, nil, err
	}

	// Getting comment data
	jsondata := make(map[string]string)
	err = json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return nil, nil, err
	}
	state := jsondata["state"]

	// Edit Issue state
	data := &github.IssueRequest{
		State: github.String(state),
	}
	return client.Issues.Edit(ctx, "deploidstudio", slug, iid, data)
}

func GetRepoBranches(c echo.Context) ([]any, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	slug := c.PathParam("project")

	// Get branches
	branches, _, err := client.Repositories.ListBranches(ctx, "deploidstudio", slug, nil)
	if err != nil {
		return nil, err
	}
	list := []any{}
	for _, branch := range branches {
		data := map[string]interface{}{
			"name": branch.Name,
			"sha":  branch.Commit.SHA,
		}

		list = append(list, data)
	}

	// Return data
	return list, nil
}

func GetRepoCommits(c echo.Context) ([]any, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	slug := c.PathParam("project")

	// Get commits
	commits, _, err := client.Repositories.ListCommits(ctx, "deploidstudio", slug, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	list := []any{}
	for _, commit := range commits {
		data := map[string]interface{}{
			"sha":     commit.SHA,
			"message": commit.Commit.Message,
			"author": map[string]interface{}{
				"avatar": commit.Author.AvatarURL,
				"name":   commit.Author.Login,
			},
			"url":     "https://" + slug + ".deploid.studio",
			"updated": commit.Commit.Author.Date,
		}

		list = append(list, data)
	}

	// Return data
	return list, nil
}

func GetRepoWorkfows(c echo.Context) ([]any, error) {

	// Get Github Auth Token
	ctx, client := GithubAuth()

	// Getting params
	slug := c.PathParam("project")
	//branch := c.PathParam("branch")

	// Get workflows
	opts := &github.ListWorkflowRunsOptions{
		//Branch: branch,
	}
	workflows, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, "deploidstudio", slug, opts)
	if err != nil {
		return nil, err
	}

	list := []any{}
	for _, workflow := range workflows.WorkflowRuns {
		data := map[string]interface{}{
			"branch":  workflow.HeadBranch,
			"name":    workflow.Name,
			"status":  workflow.Status,
			"created": workflow.CreatedAt,
		}

		list = append(list, data)
	}

	// Return data
	return list, nil
}

func CreateStripeCheckout(c echo.Context, app *pocketbase.PocketBase) (error, string) {

	// Getting params
	jsondata := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return err, ""
	}

	// Setting variables
	name := jsondata["name"]
	slug := slug.Make(name) + "-" + RandStringRunes(6)
	price := jsondata["price"]
	plan := jsondata["plan"]
	template := jsondata["template"]
	region := jsondata["region"]
	email := jsondata["email"]
	user := jsondata["user"]
	cupon := jsondata["cupon"]
	url := os.Getenv("URL")

	projects, err := QueryProjects(app, slug)

	if err != nil {
		return err, ""
	}

	if len(projects) > 0 {
		return errors.New("slug taken"), ""
	}

	checkoutParams := &stripe.CheckoutSessionParams{
		Mode:          stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		CustomerEmail: stripe.String(email),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(price),
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"name":     name,
				"slug":     slug,
				"template": template,
				"region":   region,
				"user":     user,
				"email":    email,
			},
		},
		SuccessURL: stripe.String(url + "/" + slug),
		CancelURL:  stripe.String(url + "/deploy/" + template),
	}

	if plan == "deploy" {
		checkoutParams.SubscriptionData.TrialPeriodDays = stripe.Int64(7)
	}

	if cupon != "" {
		checkoutParams.Discounts = []*stripe.CheckoutSessionDiscountParams{
			{
				Coupon: stripe.String(cupon),
			},
		}
	}

	s, err := session.New(checkoutParams)

	if err != nil {
		return err, ""
	}

	return nil, s.URL
}

func GetStripeCustomerPortal(c echo.Context, app *pocketbase.PocketBase) (string, error) {

	// Check if user is allowed
	project := c.PathParam("project")
	data, err := CheckProjectPermission(app, c, project)
	if err != nil {
		return "", err
	}

	// Getting params
	jsondata := make(map[string]string)
	err = json.NewDecoder(c.Request().Body).Decode(&jsondata)
	if err != nil {
		return "", err
	}

	// Setting variables
	customer := jsondata["customer"]
	slug := data.GetString("slug")
	website := os.Getenv("URL")
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customer),
		ReturnURL: stripe.String(website + "/" + slug),
	}

	// Create portal session
	ps, err := portal.New(params)
	if err != nil {
		return "", err
	}

	// Return portal session
	return ps.URL, nil
}

func HandleStripeEvent(c echo.Context, app *pocketbase.PocketBase) (bool, error) {
	const MaxBodyBytes = int64(65536)

	// Getting raw body
	bodyReader := c.Request().Body
	payload, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		log.Println(err)
		return false, err
	}

	// Checking event
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signatureHeader := c.Request().Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		log.Println(err)
		return false, err
	}

	log.Printf(event.Type)
	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Println(err)
			return false, err
		}

		slug := subscription.Metadata["slug"]

		err = UpdateSubscriptionData(app, slug, "canceled", "", "")
		if err != nil {
			log.Println(err)
			return false, err
		}

		_, err = CreateWorkflowDispatch(slug, "stop.yml", map[string]interface{}{})

		if err != nil {
			log.Println(err)
			return false, err
		}

		log.Println("Workflow dispatch created (stop.yml)")

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Println(err)
			return false, err
		}

		plan, err := GetPlanDetails(subscription.Items.Data[0].Price.Product.ID)

		if err != nil {
			log.Println(err)
			return false, err
		}

		slug := subscription.Metadata["slug"]

		err = UpdateSubscriptionData(app, slug, string(subscription.Status), plan.ID, plan.Name)
		if err != nil {
			log.Println(err)
			return false, err
		}

	case "customer.subscription.created":

		// Get subscription
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Println(err)
			return false, err
		}

		// Create project
		pid, err := CreateProject(app, subscription.Metadata["name"], subscription.Metadata["slug"], subscription.Metadata["user"], subscription.Metadata["template"], subscription.Metadata["region"])
		if err != nil {
			log.Println(err)
			return false, err
		}

		// Get plan details
		plan, err := GetPlanDetails(subscription.Items.Data[0].Price.Product.ID)
		if err != nil {
			log.Println(err)
			return false, err
		}

		// Create subscription
		err = CreateSubscription(app, pid, "trial", subscription.Customer.ID, plan.ID, plan.Name)
		if err != nil {
			log.Println(err)
			return false, err
		}

		// Create repo with template
		repo, _, err := CreateNewRepoFromTemplate(subscription.Metadata["slug"], subscription.Metadata["template"])
		if err != nil {
			log.Println(err)
		}

		// Wait some time
		log.Println("Repository created ", *repo.Name)
		time.Sleep(6 * time.Second)

		// Create machine
		parsed := strings.ReplaceAll(subscription.Metadata["slug"], "-", ".")
		inputs := map[string]interface{}{
			"slug":      subscription.Metadata["slug"],
			"name":      subscription.Metadata["name"],
			"bundle_id": fmt.Sprintf("com.deploid.app.%s", parsed),
			"region":    subscription.Metadata["region"],
			// "disk_size": 1,
		}

		_, err = CreateWorkflowDispatch(subscription.Metadata["slug"], "setup.yml", inputs)

		if err != nil {
			log.Println(err)
		}

		log.Println("Workflow dispatch created")

		// Create Slack channel
		// channel, err := CreateSlackChannel(subscription.Metadata["slug"], subscription.Metadata["email"])
		// if err != nil {
		// 	log.Println(err)
		// }

		// Update project
		// err = UpdateProject(app, pid, channel)
		// if err != nil {
		// 	log.Println(err)
		// }

	case "customer.subscription.trial_will_end":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Println(err)
			return false, err
		}

		// TODO: Send email to user

	default:
		return false, errors.New("unhandled event type")
	}

	return true, nil
}

func HandleGithubUpdate(c echo.Context, app *pocketbase.PocketBase) error {

	// Verifying signature
	event, payload, err := VerifyGithubWebhookSignature(c)
	if err != nil {
		return err
	}

	switch event {
	case "issue_comment":
		{
			if payload["action"] == "created" {
				// Check if body do not contains "－ "
				isFromDash := strings.Contains(payload["comment"].(map[string]interface{})["body"].(string), "－ ")

				if !isFromDash {

					// Get Project
					project, err := GetProjectWithSlug(app, payload["repository"].(map[string]interface{})["name"].(string))
					if err != nil {
						log.Println(err)
					}

					// Get members
					members, err := GetUsersById(app, project.GetStringSlice("members"))
					if err != nil {
						log.Println(err)
					}

					// Build email list
					emails := []string{}
					for _, member := range members {
						emails = append(emails, member.GetString("email"))
					}

					// Build email template
					html, err := EmailTemplate(event, payload)
					if err != nil {
						log.Println(err)
						return err
					}

					// Send email
					err = SendEmail("New comment on task", html, emails, app)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

	}

	return nil
}

func VerifyGithubWebhookSignature(c echo.Context) (string, map[string]any, error) {

	// Checking event
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	signature := c.Request().Header.Get("X-Hub-Signature-256")
	event := c.Request().Header.Get("X-GitHub-Event")
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return event, nil, err
	}

	// Checking signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	// Checking signature
	if !hmac.Equal([]byte(signature), []byte("sha256="+expectedMAC)) {
		return event, nil, errors.New("invalid signature")
	}

	// Load the body data into a Payload struct
	payload := make(map[string]any)
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return event, nil, err
	}

	return event, payload, nil
}

func EmailTemplate(event string, payload map[string]interface{}) (string, error) {
	switch event {
	case "issue_comment":
		markdown := payload["comment"].(map[string]interface{})["body"].(string)
		bytes := blackfriday.MarkdownCommon([]byte(markdown))
		rendered := string(bytes)
		return fmt.Sprintf("<p>Hello, you have a response for your task:</p>\n <p>%s</p> \n <p>－ %s</p> \n <p>Deplōid Studio Team</p>", rendered, payload["comment"].(map[string]interface{})["user"].(map[string]interface{})["login"].(string)), nil

	default:
		return "", errors.New("Unhandled event type")
	}

}

func SendEmail(subject string, html string, emails []string, app *pocketbase.PocketBase) error {

	// Build email list
	list := []mail.Address{}
	for _, email := range emails {
		list = append(list, mail.Address{
			Address: email,
		})
	}

	// Send email
	message := &mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To:      list,
		Subject: subject,
		HTML:    html,
	}
	return app.NewMailClient().Send(message)
}

func CheckProjectPermission(app *pocketbase.PocketBase, c echo.Context, pid string) (*models.Record, error) {

	// Get project
	dao := app.Dao()
	project, err := dao.FindRecordById("projects", pid, nil)
	if err != nil {
		return nil, err
	}

	// Get user
	user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

	// Check if user is allowed
	if !ContainsString(project.GetStringSlice("members"), user.Id) {
		return nil, fmt.Errorf("user not allowed")
	}

	return project, nil
}

func GetProjectWithSlug(app *pocketbase.PocketBase, slug string) (*models.Record, error) {

	// Get project
	dao := app.Dao()
	projects, err := dao.FindRecordsByExpr("projects", dbx.HashExp{"slug": slug})
	if err != nil {
		return nil, err
	}

	// Check if has results
	if len(projects) == 0 {
		return nil, fmt.Errorf("Project not found")
	}

	project := projects[0]

	// Return first result
	return project, nil
}

func GetUsersById(app *pocketbase.PocketBase, uids []string) ([]*models.Record, error) {

	// Get project
	dao := app.Dao()
	users, err := dao.FindRecordsByIds("users", uids)
	if err != nil {
		return nil, err
	}

	// Return results
	return users, nil
}

func GetUserById(app *pocketbase.PocketBase, uid string) (*models.Record, error) {

	// Get project
	dao := app.Dao()
	user, err := dao.FindRecordById("users", uid)
	if err != nil {
		return nil, err
	}

	// Return results
	return user, nil
}

func CreateRelease(name string, tag string, body string) (*github.RepositoryRelease, *github.Response, error) {

	// Get Github client
	ctx, client := GithubAuth()

	// Create release
	opts := github.RepositoryRelease{
		TagName: github.String(tag),
		Name:    github.String(tag),
		Body:    github.String(body),
	}

	// Return release
	return client.Repositories.CreateRelease(ctx, "deploidstudio", name, &opts)
}

func CreateNewRepoFromTemplate(name string, template string) (*github.Repository, *github.Response, error) {

	// Get Github client
	ctx, client := GithubAuth()

	// Create repo
	repo := &github.TemplateRepoRequest{Name: github.String(name), Owner: github.String("deploidstudio"), Private: github.Bool(true), IncludeAllBranches: github.Bool(true)}

	// Return repo
	return client.Repositories.CreateFromTemplate(ctx, "deploidstudio", template, repo)
}

func CreateSubscription(app *pocketbase.PocketBase, pid string, status string, customer string, plan string, name string) error {

	// Creating subscription record
	dao := app.Dao()
	collection, err := dao.FindCollectionByNameOrId("subscriptions")
	if err != nil {
		return err
	}

	// Create record
	record := models.NewRecord(collection)
	record.Set("project", pid)
	record.Set("status", status)
	record.Set("customer", customer)
	record.Set("plan", plan)
	record.Set("name", name)
	if err := dao.SaveRecord(record); err != nil {
		return err
	}

	// We need nothing from here
	return nil
}

func CreateProject(app *pocketbase.PocketBase, name string, slug string, user string, template string, region string) (string, error) {

	// Create project
	dao := app.Dao()
	collection, err := dao.FindCollectionByNameOrId("projects")
	if err != nil {
		return "", err
	}

	// Create record
	record := models.NewRecord(collection)
	record.Set("name", name)
	record.Set("slug", slug)
	record.Set("user", user)
	record.Set("members", user)
	record.Set("template", template)
	record.Set("region", region)
	err = dao.SaveRecord(record)
	if err != nil {
		return "", err
	}

	// Return project id
	return record.Id, nil
}

func GetPlanDetails(id string) (*stripe.Product, error) {
	plan, err := product.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func GetPricesWithProducts() *price.Iter {
	params := &stripe.PriceListParams{
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	params.AddExpand("data.product")

	values := price.List(params)

	return values
}

func CreateWorkflowDispatch(name string, file string, inputs map[string]interface{}) (*github.Response, error) {
	ctx, client := GithubAuth()
	opts := github.CreateWorkflowDispatchEventRequest{
		Ref:    *github.String("main"),
		Inputs: inputs,
	}
	return client.Actions.CreateWorkflowDispatchEventByFileName(ctx, "deploidstudio", name, file, opts)
}

func UpdateProject(app *pocketbase.PocketBase, pid string, channel string) error {
	// Updating project data
	dao := app.Dao()

	// Get project record
	project, err := dao.FindRecordById("projects", pid)
	if err != nil {
		return err
	}

	// Update fields
	if channel != "" {
		project.Set("channel", channel)
	}

	// Update record
	err = dao.SaveRecord(project)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSubscriptionData(app *pocketbase.PocketBase, slug string, status string, plan string, name string) error {
	// Updating project data
	dao := app.Dao()

	// Get project record
	expr := dbx.HashExp{"slug": slug}
	records, err := dao.FindRecordsByExpr("projects", expr)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errors.New("project not found")
	}
	pid := records[0].Id

	// Update subscription
	expr = dbx.HashExp{"project": pid}
	records, err = dao.FindRecordsByExpr("subscriptions", expr)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errors.New("subscription not found")
	}

	// Update fields
	subscription := records[0]
	if status != "" {
		subscription.Set("status", status)
	}
	if plan != "" {
		subscription.Set("plan", plan)
	}
	if name != "" {
		subscription.Set("name", name)
	}

	// Update record
	errr := dao.SaveRecord(subscription)
	if err != nil {
		return errr
	}

	return nil
}

// func CreateSlackChannel(slug string, email string) (string, error) {
// 	token := os.Getenv("SLACK_TOKEN")
// 	api := slack.New(token)

// 	channel, err := api.CreateConversation(slack.CreateConversationParams{
// 		ChannelName: slug,
// 		IsPrivate:   true,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	_, err = api.InviteUsersToConversation(channel.ID, "U059DBHUTJT")
// 	if err != nil {
// 		return "", err
// 	}

// 	url := fmt.Sprintf("https://slack.com/api/conversations.inviteShared?channel=%d&emails=%d", channel.ID, email)
// 	bearer := fmt.Sprintf("Bearer %s", token)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	req.Header.Add("Content-Type", "application/json; charset=utf-8")
// 	req.Header.Add("Authorization", bearer)
// 	client := &http.Client{}
// 	_, err = client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}

// 	return channel.ID, nil
// }

func QueryProjects(app *pocketbase.PocketBase, slug string) ([]*models.Record, error) {

	// Updating project data
	dao := app.Dao()

	// Get project record
	expr := dbx.HashExp{"slug": slug}
	records, err := dao.FindRecordsByExpr("projects", expr)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func GithubAuth() (context.Context, *github.Client) {

	// Get Github Auth Token
	token := os.Getenv("GITHUB_TOKEN")

	// Create Github Client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Return Github Client
	return ctx, client
}

func EncryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {

	// Decode the public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	// Create a 32 byte array from the decoded public key
	var boxKey [32]byte
	copy(boxKey[:], decodedPublicKey)
	secretBytes := []byte(secretValue)
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, crypto_rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("box.SealAnonymous failed with error %w", err)
	}

	// Encode the encrypted bytes to base64
	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}

	// Return the encrypted secret
	return encryptedSecret, nil
}

func ContainsString(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func RandStringRunes(n int) string {
	return uniuri.NewLenChars(6, []byte("abcdefghijklmnopqrstuvwxyz"))
}
