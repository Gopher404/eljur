package httpHandler

import (
	"context"
	"eljur/internal/service/users"
	"eljur/pkg/tr"
)

type headerTmpData struct {
	TmpData map[string]any
}

func (h *headerTmpData) setData(data ...any) {
	for i := 0; i < len(data)-1; i += 2 {
		h.TmpData[data[i].(string)] = data[i+1]
	}
}

func (h *headerTmpData) initData() {
	h.TmpData = make(map[string]any)
}

type HeaderSetter interface {
	setData(data ...any)
	initData()
}

func (h *Handler) SetHeaderData(ctx context.Context, headerData HeaderSetter, activePage string, login string) error {
	headerData.initData()
	headerData.setData("active-page", activePage)

	user, err := h.usersService.GetUser(ctx, login)
	if err != nil {
		return tr.Trace(err)
	}
	headerData.setData(
		"user-name", user.Name,
		"is-admin", user.Perm >= users.PermAdmin,
	)

	return nil
}
