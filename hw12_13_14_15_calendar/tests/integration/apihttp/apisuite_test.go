package apihttp

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	// pg driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/streadway/amqp"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/config"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository/psql"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/client/calendar_service"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/models"
)

type rabbit struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      [][]byte
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}
}
type apiSuite struct {
	cfg    *config.TestCfg
	db     *sql.DB
	rabbit *rabbit

	modifyResp *models.ModifyEventResponse
	getResp    *models.GetEventsResponse
}

func initNewAPI() (*apiSuite, error) {
	a := new(apiSuite)
	r := new(rabbit)
	a.rabbit = r
	var err error
	a.cfg, err = config.LoadTest()
	if err != nil {
		return nil, fmt.Errorf("test config read error: %v", err)
	}

	if a.cfg.App.RepoType != psql.Type {
		return nil, fmt.Errorf("only psql repo type is allowed for this tests")
	}

	return a, nil
}

func (a *apiSuite) dbConnect(*messages.Pickle) {
	var err error
	a.db, err = sql.Open("pgx", a.cfg.PSQL.GetDSN())
	if err != nil {
		log.Fatalf("open pgx error: %v", err)
	}

	a.db.Stats()
	err = a.db.PingContext(context.Background())
	if err != nil {
		log.Fatalf("open pgx error: %v", err)
	}
}

func (a *apiSuite) dbClose(*messages.Pickle, error) {
	err := a.db.Close()
	if err != nil {
		log.Fatalf("Db conn close error: %v", err)
	}
}

func (a *apiSuite) startConsuming() error {
	a.rabbit.messages = make([][]byte, 0)
	a.rabbit.messagesMutex = new(sync.RWMutex)
	a.rabbit.stopSignal = make(chan struct{})

	var err error

	a.rabbit.conn, err = amqp.Dial(a.cfg.Rabbit.GetDSN())
	if err != nil {
		log.Fatalf("rabbit connect error: %v", err)
	}

	a.rabbit.ch, err = a.rabbit.conn.Channel()
	if err != nil {
		log.Fatalf("rabbit channel error: %v", err)
	}

	queueName := a.cfg.Rabbit.QueueName
	_, err = a.rabbit.ch.QueueDeclare(queueName, true, true, true, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	err = a.rabbit.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("setting QoS error: %v", err)
	}

	// подключаем тестовую очередь к exchange. Т.к. он типа fanout, то никаких проблем с этим нет
	err = a.rabbit.ch.QueueBind(queueName, "", a.cfg.Rabbit.ExchangeName, false, nil)
	if err != nil {
		log.Fatalf("queue bind error: %v", err)
	}
	eventsCh, err := a.rabbit.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume bind error: %v", err)
	}

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-eventsCh:
				a.rabbit.messagesMutex.Lock()
				a.rabbit.messages = append(a.rabbit.messages, event.Body)
				a.rabbit.messagesMutex.Unlock()
			}
		}
	}(a.rabbit.stopSignal)

	return nil
}

func (a *apiSuite) stopConsuming() error {
	a.rabbit.stopSignal <- struct{}{}

	err := a.rabbit.ch.Close()
	if err != nil {
		log.Fatalf("rabbit channel close error: %v", err)
	}

	err = a.rabbit.conn.Close()
	if err != nil {
		log.Fatalf("rabbit connection close error: %v", err)
	}

	a.rabbit.messages = nil
	return nil
}

func (a *apiSuite) clearEvents(*messages.Pickle) {
	_, err := a.db.ExecContext(context.Background(), "DELETE FROM events")
	if err != nil {
		log.Fatalf("delete events error: %v", err)
	}
}

func (a *apiSuite) thereAreEvents(events *godog.Table) error {
	var fields []string
	var marks []string
	head := events.Rows[0].Cells
	i := 1
	for _, cell := range head {
		fields = append(fields, cell.Value)
		if cell.Value == "during" {
			marks = append(marks, fmt.Sprintf("tstzrange($%d, $%d, '[]')", i, (i+1)))
			i++
		} else {
			marks = append(marks, fmt.Sprintf("$%d", i))
		}
		i++
	}

	stmt, err := a.db.Prepare("INSERT INTO events (" + strings.Join(fields, ", ") + ") VALUES(" + strings.Join(marks, ", ") + ")")
	if err != nil {
		return err
	}
	for i := 1; i < len(events.Rows); i++ {
		var vals []interface{}
		for n, cell := range events.Rows[i].Cells {
			switch head[n].Value {
			case "title", "description", "user_id", "id":
				vals = append(vals, cell.Value)
			case "during":
				s := strings.SplitN(cell.Value, " ", 2)
				start, end := s[0], s[1]
				startAt, err := time.Parse(time.RFC3339, start)
				if err != nil {
					return err
				}
				endAt, err := time.Parse(time.RFC3339, end)
				if err != nil {
					return err
				}
				vals = append(vals, startAt)
				vals = append(vals, endAt)
			case "notify_at":
				var notifyAt *time.Time
				if cell.Value != "" {
					notify, err := time.Parse(time.RFC3339, cell.Value)
					if err != nil {
						return err
					}
					notifyAt = &notify
				}
				vals = append(vals, notifyAt)
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		if _, err = stmt.Exec(vals...); err != nil {
			return err
		}
	}
	return nil
}

func (a *apiSuite) getClient() calendar_service.ClientService {
	transport := client.New(a.cfg.Server.Addr, "", nil)

	return calendar_service.New(transport, strfmt.Default)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api, err := initNewAPI()
	if err != nil {
		log.Fatalf("%v", err)
	}

	ctx.BeforeScenario(api.dbConnect)
	ctx.BeforeScenario(api.clearEvents)
	ctx.AfterScenario(api.dbClose)

	ctx.Step(`^I send create request with title "([^"]*)" description "([^"]*)" startAt "([^"]*)" endAt "([^"]*)" notifyAt "([^"]*)"$`, api.iSendCreateRequest)
	ctx.Step(`^I send update request with eventID "(\d+)" title "([^"]*)" description "([^"]*)" startAt "([^"]*)" endAt "([^"]*)" notifyAt "([^"]*)"$`, api.iSendUpdateRequest)
	ctx.Step(`^the response should be without errors$`, api.theRespShouldNotHasError)
	ctx.Step(`^the response has correct event id$`, api.theRespShouldHasCorrectEventID)
	ctx.Step(`^the response should has error text "([^"]*)"$`, api.theRespShouldHasErrorText)
	ctx.Step(`^there are events:$`, api.thereAreEvents)
	ctx.Step(`^I send get events day request with beginAt "([^"]*)"$`, api.iSendGetEventsDayRequest)
	ctx.Step(`^I send get events week request with beginAt "([^"]*)"$`, api.iSendGetEventsWeekRequest)
	ctx.Step(`^I send get events month request with beginAt "([^"]*)"$`, api.iSendGetEventsMonthRequest)
	ctx.Step(`^count event in the response should be (\d+)$`, api.eventCountShouldBe)
	ctx.Step(`^I should receive message with title "([^"]*)"$`, api.iShouldReceiveEventWithTitle)
	ctx.Step(`^the response should has title "([^"]*)"$`, api.theRespShouldHasTitle)
	// не нашел способа как еще запускать консумеры только для нужных сценариев
	ctx.Step(`^I start consuming$`, api.startConsuming)
	ctx.Step(`^I stop consuming$`, api.stopConsuming)
}
