package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2023-04-24 17:59:43.925Z",
				"updated": "2023-06-17 17:48:41.406Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"protected": false
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": false
				}
			},
			{
				"id": "zmk76cdfosvkvh5",
				"created": "2023-04-24 18:08:25.741Z",
				"updated": "2023-06-17 17:48:41.406Z",
				"name": "projects",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "bouzmhdi",
						"name": "slug",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "8r1vnwim",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "c15srptx",
						"name": "template",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ltc0wgqd",
						"name": "region",
						"type": "select",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"ams",
								"cdg",
								"den",
								"dfw",
								"ewr",
								"gru",
								"hkg",
								"iad",
								"jnb",
								"lax",
								"lhr",
								"mad",
								"mia",
								"nrt",
								"ord",
								"scl",
								"sea",
								"sin",
								"sjc",
								"syd",
								"yul",
								"yyz"
							]
						}
					},
					{
						"system": false,
						"id": "nae2wreh",
						"name": "user",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "surie0li",
						"name": "members",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 50,
							"displayFields": []
						}
					}
				],
				"indexes": [],
				"listRule": "members.id ?= @request.auth.id",
				"viewRule": "members.id ?= @request.auth.id",
				"createRule": "@request.auth.id != null",
				"updateRule": "user.id = @request.auth.id",
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "qbowm1gt1ne68j4",
				"created": "2023-04-24 18:19:34.508Z",
				"updated": "2023-06-17 17:48:41.407Z",
				"name": "subscriptions",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "zrjsx4ih",
						"name": "status",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "p2eankki",
						"name": "name",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "y5hwbedc",
						"name": "plan",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "7jvtpxwd",
						"name": "customer",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ze2jjks7",
						"name": "project",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "zmk76cdfosvkvh5",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					}
				],
				"indexes": [],
				"listRule": "project.members.id ?= @request.auth.id",
				"viewRule": "project.members.id ?= @request.auth.id",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
