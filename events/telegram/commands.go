package telegram

import (
	"errors"
	"log"
	"net/url"
	"notifier/lib/e"
	"notifier/storage"
	"strings"
)

const (
	RndCmd   = "/rand"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (d *Dispatcher) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	//	add page: http://..
	//	rand page" /rand
	// help: /help
	// start: /start
	if isAddCmd(text) {
		return d.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return d.sendRandom(chatID, username)
	case HelpCmd:
		return d.sendHelp(chatID)
	case StartCmd:
		return d.sendHello(chatID)
	default:
		return d.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (d *Dispatcher) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := d.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return d.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := d.storage.Save(page); err != nil {
		return err
	}

	if err := d.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (d *Dispatcher) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := d.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return d.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := d.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return d.storage.Remove(page)
}

func (d *Dispatcher) sendHelp(chatID int) error {
	return d.tg.SendMessage(chatID, msgHelp)
}

func (d *Dispatcher) sendHello(chatID int) error {
	return d.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
