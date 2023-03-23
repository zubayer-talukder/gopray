package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	objs "github.com/SakoDroid/telego/objects"
	"github.com/escalopa/gopray/pkg/core"
)

func (h *Handler) SetLang(u *objs.Update) {
	if true {
		h.simpleSend(u.Message.Chat.Id, "This feature is not available yet.", 0)
		return
	}
	var messageID int
	chatID := u.Message.Chat.Id
	kb := h.b.CreateInlineKeyboard()

	ctx, cancel := context.WithTimeout(h.userCtx[u.Message.Chat.Id].ctx, 1*time.Minute)
	// Deletes the message after the button is pressed or after 1 hour.
	go func() {
		defer cancel()
		<-ctx.Done()
		h.deleteMessage(chatID, messageID)
	}()

	for i, language := range core.AvaliableLanguages() {
		//Adds a callback button with handler.
		row := i/2 + 1 // 2 buttons per row.
		kb.AddCallbackButtonHandler(language, language, row, func(u *objs.Update) {
			defer cancel()
			// Sets the language.
			err := h.u.SetLang(h.c, chatID, u.CallbackQuery.Data)
			if err != nil {
				log.Printf("failed to set language to %s: %v", u.CallbackQuery.Data, err)
				_, err = h.b.AdvancedMode().AAnswerCallbackQuery(u.CallbackQuery.Id,
					fmt.Sprintf("Failed to set language to %s, Please try again later", u.CallbackQuery.Data),
					true, "", 0)
				if err != nil {
					log.Printf("failed to send callback query on /lang: %s", err)
				}
				return
			}
			h.simpleSend(chatID, fmt.Sprintf("Successfully set language to %s", u.CallbackQuery.Data), 0)
		})
	}

	// Sends the message along with the keyboard.
	r, err := h.b.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Choose language", "", u.Message.MessageId, false, false, nil, false, false, kb)
	if err != nil {
		log.Printf("failed to send message on /lang: %s", err)
	}
	messageID = r.Result.MessageId
}
