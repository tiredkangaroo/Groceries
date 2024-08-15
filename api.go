package main

import (
	"log/slog"

	"github.com/nikumar1206/puff"
	"github.com/tiredkangaroo/sculpt"
)

// NewItemInput is used for the input on path POST /api/item.
type NewItemInput struct {
	ItemName string `kind:"body" description:"Name for the new item."`
}

// GetSpecificItemInput is used for the input on path GET /api/item
type GetSpecificItemInput struct {
	ID string `kind:"path" description:"ID of item to retrieve from database."`
}

func startAPIServer(itemModel *sculpt.Model) {
	app := puff.DefaultApp("Checklist")
	apiRouter := puff.NewRouter("API", "/api")

	newItemInput := new(NewItemInput)

	apiRouter.Get("/items", "Retrieves all items in the checklist.", nil, func(c *puff.Context) {
		items, err := sculpt.RunQuery[*Item](itemModel, sculpt.Query{})
		if err != nil {
			slog.Error("unable to run query for items", slog.String("error", err.Error()))
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "unable to get items",
				},
			})
		}
		c.SendResponse(puff.JSONResponse{
			Content: items,
		})
	})
	specificItemInput := new(GetSpecificItemInput)
	apiRouter.Get("/item/{id}", "Retrieve a specfic item in the checklist by its ID.", specificItemInput, func(c *puff.Context) {
		items, err := sculpt.RunQuery[*Item](itemModel, sculpt.Query{
			Conditions: []sculpt.Condition{
				sculpt.EqualTo("ID", specificItemInput.ID),
			},
		})
		if err != nil {
			slog.Error("unable to run query for item", slog.String("error", err.Error()))
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "unable to get item",
				},
			})
			return
		}
		switch len(items) {
		case 0:
			c.SendResponse(puff.JSONResponse{
				StatusCode: 404,
				Content: map[string]any{
					"error": "no items found with id: " + specificItemInput.ID,
				},
			})
		case 1:
			c.SendResponse(puff.JSONResponse{
				Content: items[0],
			})
		default:
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "retrieved multiple items for id: " + specificItemInput.ID,
				},
			})
		}
		return
	})

	apiRouter.Post("/item", "Add a new item to the checklist.", newItemInput, func(c *puff.Context) {
		randomID, err := generateNanoID()
		if err != nil {
			slog.Error("unable to generate nano id", slog.String("error", err.Error()))
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "unable to generate random nano id",
				},
			})
			return
		}
		itemRow, err := itemModel.New(&Item{
			ID:   randomID,
			Name: newItemInput.ItemName,
		})
		if err != nil {
			slog.Error("unable to create item row", slog.String("error", err.Error()))
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "unable to create new row",
				},
			})
			return
		}
		err = itemRow.Save()
		if err != nil {
			slog.Error("unable to save item row", slog.String("error", err.Error()))
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "unable to save new row",
				},
			})
			return
		}
		c.SendResponse(puff.JSONResponse{
			Content: map[string]any{
				"error": nil,
			},
		})
		return
	})

	app.IncludeRouter(apiRouter)
	app.ListenAndServe(":8080")
}
