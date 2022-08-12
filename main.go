package main

import (
	"flag"
	"log"
	event_consumer "notifier/consumer/event-consumer"

	tgClient "notifier/clients/telegram"

	"notifier/events/telegram"
	"notifier/storage/files"
)

const (
	tgHost = "api.telegram.org"
	//zulip       = "zulip.etecs.ru"
	storagePath = "files_storage"
	batchSize   = 100
)

// 5503985166:AAERqm49nLmZUmK3QgPh_x43P-ZS3KWDoFk

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

	// zulipClient = zulip.New(token)
}

func mustToken() string {
	token := flag.String("tg", "", "telegram bot token")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not valid")
	}
	return *token
}
