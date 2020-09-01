package sender

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

func GetSenderFunc() func(<-chan amqp.Delivery) {
	return func(msgCh <-chan amqp.Delivery) {
		for msg := range msgCh {
			n := new(repository.Notification)
			err := json.Unmarshal(msg.Body, n)
			if err != nil {
				log.Errorf("unmarshal error: %v", err)
				err = msg.Nack(false, false)
				if err != nil {
					log.Errorf("Nack error: %v", err)
				}

				continue
			}

			log.Infof("received message: %v", n)
			err = msg.Ack(false)
			if err != nil {
				log.Errorf("Ack error: %v", err)
			}
		}
	}
}
