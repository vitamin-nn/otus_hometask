package rmq

import "github.com/streadway/amqp"

// здесь можно, конечно, совсем попробовать отказаться от упоминаний какого-либо типа
// брокеров сообщений (как требуется в ТЗ), но особого смысла в этом не вижу
// перекладывать сообщения из одного канала в другой в данный момент выглядит не очень разумным.
type Consumer interface {
	Handle(fn func(<-chan amqp.Delivery), threads int) error
	Close()
}

type Publisher interface {
	Connect() error
	Publish(body interface{}, routingKey string) error
	Close()
}
