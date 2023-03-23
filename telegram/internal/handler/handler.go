package handler

import (
	"context"
	"log"

	"github.com/SakoDroid/telego"

	"github.com/escalopa/gopray/telegram/internal/application"
)

type Handler struct {
	c context.Context
	b *telego.Bot
	u *application.UseCase

	botOwner int                 // Bot owner's ID.
	userCtx  map[int]userContext // userID => latest user context
}

func New(ctx context.Context, b *telego.Bot, ownerID int, u *application.UseCase) *Handler {
	return &Handler{
		b:        b,
		u:        u,
		c:        ctx,
		botOwner: ownerID,
		userCtx:  make(map[int]userContext),
	}
}

func (h *Handler) Run() error {
	err := h.register()
	if err != nil {
		return err
	}
	h.setupBundler()
	go h.notifySubscribers() // Notify subscriber about the prayer times.
	return nil
}

func (h *Handler) register() error {
	var err error
	err = h.b.AddHandler("/start", h.contextWrapper(h.Start), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/help", h.contextWrapper(h.Help), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/subscribe", h.contextWrapper(h.Subscribe), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/unsubscribe", h.contextWrapper(h.Unsubscribe), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/today", h.contextWrapper(h.GetPrayers), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/date", h.contextWrapper(h.GetPrayersByDate), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/lang", h.contextWrapper(h.SetLang), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/feedback", h.contextWrapper(h.Feedback), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/bug", h.contextWrapper(h.Bug), "all")
	if err != nil {
		return err
	}

	//////////////////////////
	///// Admin Commands /////
	//////////////////////////

	err = h.b.AddHandler("/respond", h.admin(h.contextWrapper(h.Respond)), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/subs", h.admin(h.contextWrapper(h.GetSubscribers)), "all")
	if err != nil {
		return err
	}
	err = h.b.AddHandler("/sall", h.admin(h.contextWrapper(h.SendAll)), "all")
	if err != nil {
		return err
	}
	return nil
}

// TODO: Implement bundler for multi language support.
func (h *Handler) setupBundler() {}

// simpleSend sends a simple message to the chat with the given chatID & text and replyTo.
func (h *Handler) simpleSend(chatID int, text string, replyTo int) (messageID int) {
	r, err := h.b.SendMessage(chatID, text, "", replyTo, false, false)
	if err != nil {
		log.Printf("failed to send message on simpleSend: %s", err)
		return 0
	}
	return r.Result.MessageId
}

// cancelOperation checks if the message is /cancel and sends a response.
// Returns true if the message is /cancel.
func (h *Handler) cancelOperation(message, response string, chatID int) bool {
	if message == "/cancel" {
		h.simpleSend(chatID, response, 0)
		return true
	}
	return false
}

// deleteMessage deletes the message with the given chatID & messageID.
// If error occurs, it will be logged.
func (h *Handler) deleteMessage(chatID, messageID int) {
	editor := h.b.GetMsgEditor(chatID)
	_, err := editor.DeleteMessage(messageID)
	if err != nil {
		log.Printf("failed to delete message: %s", err)
		return
	}
}
