package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qlik-oss/enigma-go"
)

type Bookmark struct {
	Info *enigma.NxInfo `json:"qInfo,omitempty"`
}

func (b Bookmark) validate() error {
	if b.Info == nil {
		return errors.New("missing qInfo attribute")
	}
	if b.Info.Id == "" {
		return errors.New("missing qInfo qId attribute")
	}
	if b.Info.Type == "" {
		return errors.New("missing qInfo qType attribute")
	}
	if b.Info.Type != "bookmark" {
		return errors.New("bookmarks must have qType: bookmark")
	}
	return nil
}

// ListBookmarks lists all dimenions in an app
func ListBookmarks(ctx context.Context, doc *enigma.Doc) []NamedItem {
	props := &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "corectl_entity_list",
		},
		BookmarkListDef: &enigma.BookmarkListDef{
			Type: "bookmark",
			Data: json.RawMessage(`{
				"id":"/qInfo/qId",
				"title":"/qMetaDef/title"
			}`),
		},
	}
	sessionObject, _ := doc.CreateSessionObject(ctx, props)
	defer doc.DestroySessionObject(ctx, sessionObject.GenericId)
	layout, _ := sessionObject.GetLayout(ctx)
	result := []NamedItem{}
	for _, item := range layout.BookmarkList.Items {
		parsedRawData := &ParsedEntityListData{}
		json.Unmarshal(item.Data, parsedRawData)
		result = append(result, NamedItem{Title: parsedRawData.Title, Id: item.Info.Id})
	}
	return result
}

// SetBookmarks adds all bookmarks that match the specified glob pattern
func SetBookmarks(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "bookmarks")
	if err != nil {
		FatalError("could not interpret glob pattern: ", err)
	}
	for _, path := range paths {
		rawEntities, err := parseEntityFile(path)
		if err != nil {
			FatalErrorf("could not parse file %s: %s", path, err)
		}
		for _, raw := range rawEntities {
			var bm Bookmark
			err := json.Unmarshal(raw, &bm)
			if err != nil {
				FatalErrorf("could not parse data in file %s: %s", path, err)
			}
			err = bm.validate()
			if err != nil {
				FatalErrorf("validation error in file %s: %s", path, err)
			}
			err = setBookmark(ctx, doc, bm.Info.Id, raw)
			if err != nil {
				FatalError(err)
			}
		}
	}
}

func setBookmark(ctx context.Context, doc *enigma.Doc, bookmarkID string, raw json.RawMessage) error {
	bookmark, err := doc.GetBookmark(ctx, bookmarkID)
	if err != nil {
		return err
	}
	if bookmark.Handle != 0 {
		LogVerbose("Updating bookmark " + bookmarkID)
		err = bookmark.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not update %s with %s: %s", "bookmark", bookmarkID, err)
		}
	} else {
		LogVerbose("Creating bookmark " + bookmarkID)
		_, err = doc.CreateBookmarkRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not create %s with %s: %s", "bookmark", bookmarkID, err)
		}
	}
	return nil
}
