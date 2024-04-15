package httpServer

import (
	"context"
	"eljur/pkg/tr"
)

type headerTmpData struct {
	UserName   string
	ActivePage string
}

func (h *headerTmpData) SetUserName(username string) {
	h.UserName = username
}

func (h *headerTmpData) SetActivePage(page string) {
	h.ActivePage = page
}

type HeaderSetter interface {
	SetUserName(username string)
	SetActivePage(username string)
}

func (h *Handler) SetHeaderData(ctx context.Context, headerData HeaderSetter, activePage string, login string) error {
	headerData.SetActivePage(activePage)
	userName, err := h.usersService.GetUserName(ctx, login)
	if err != nil {
		return tr.Trace(err)
	}
	headerData.SetUserName(userName)
	return nil
}
