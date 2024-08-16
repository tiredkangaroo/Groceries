package main

import (
	"log/slog"
	"time"

	"github.com/nikumar1206/puff"
	"github.com/tiredkangaroo/sculpt"
)

// NewItemInput is used for the input on path POST /api/item.
type NewItemInput struct {
	ItemName string `kind:"body" description:"Name for the new item." example:"laundry"`
}

// SpecificItemInput is used for the input on path GET /api/item and DELETE /api/item
type SpecificItemInput struct {
	ID string `kind:"path" description:"ID of item to retrieve from database." example:"5NQDFNEF099G4997AO3A0GII"`
}

func getAPIRouter(itemModel *sculpt.Model) *puff.Router {
	apiRouter := puff.NewRouter("API", "/api")

	getItemInput := new(SpecificItemInput)
	apiRouter.Get("/item/{id}", "Retrieve a specfic item in the checklist by its ID.", getItemInput, func(c *puff.Context) {
		items, err := sculpt.RunQuery[*Item](itemModel, sculpt.Query{
			Conditions: []sculpt.Condition{
				sculpt.EqualTo("ID", getItemInput.ID),
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
					"error": "no items found with id: " + getItemInput.ID,
				},
			})
		case 1:
			c.SendResponse(puff.JSONResponse{
				Content: items[0],
			})
		default:
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "retrieved multiple items for id: " + getItemInput.ID,
				},
			})
		}
		return
	})

	newItemInput := new(NewItemInput)
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
		newItem := &Item{
			ID:          randomID,
			DateCreated: time.Now().String(),
			Name:        newItemInput.ItemName,
		}
		itemRow, err := itemModel.New(newItem)
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
				"item":  *newItem,
			},
		})
		return
	})

	deleteItemInput := new(SpecificItemInput)
	apiRouter.Delete("/item/{id}", "Deletes item in the checklist that matches id.", deleteItemInput, func(c *puff.Context) {
		err := itemModel.Delete(sculpt.Query{
			Conditions: []sculpt.Condition{
				sculpt.EqualTo("ID", deleteItemInput.ID),
			},
		})
		if err != nil {
			slog.Error("unable to run delete for item", slog.String("error", err.Error()))
			c.SendResponse(puff.JSONResponse{
				Content: map[string]any{
					"error": "unable to delete item " + deleteItemInput.ID,
				},
			})
			return
		}
		c.SendResponse(puff.JSONResponse{
			Content: map[string]any{
				"error": nil,
			},
		})
	})

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
			Content: map[string]any{
				"error": nil,
				"items": items,
			},
		})
	})

	return apiRouter
}
