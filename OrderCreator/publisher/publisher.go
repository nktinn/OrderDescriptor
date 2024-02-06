package publisher

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderCreator/model"
)

func randomString(randGenerator *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randGenerator.Intn(len(charset))]
	}
	return string(b)
}

func randomPhone(randGenerator *rand.Rand) string {
	return fmt.Sprintf("+972%09d", randGenerator.Intn(1000000000))
}

func randomInt(randGenerator *rand.Rand, min, max int) int {
	return min + randGenerator.Intn(max-min+1)
}

func randomOrder(randGenerator *rand.Rand) model.Order {
	orderUID := randomString(randGenerator, 20)
	trackNumber := randomString(randGenerator, 15)
	customerID := randomString(randGenerator, 10)
	email := randomString(randGenerator, 8) + "@gmail.com"
	transaction := randomString(randGenerator, 20)

	return model.Order{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        customerID,
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              randomInt(randGenerator, 10, 200),
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: &model.Delivery{
			OrderID: orderUID,
			Name:    randomString(randGenerator, 10) + " " + randomString(randGenerator, 10),
			Phone:   randomPhone(randGenerator),
			Zip:     strconv.Itoa(randomInt(randGenerator, 100000, 999999)),
			City:    "New York",
			Address: fmt.Sprintf("%d %s %d", randomInt(randGenerator, 10, 100), randomString(randGenerator, 8), randomInt(randGenerator, 10, 100)),
			Region:  "NY",
			Email:   email,
		},
		Payment: &model.Payment{
			OrderID:       orderUID,
			TransactionID: transaction,
			RequestID:     "",
			Currency:      "USD",
			Provider:      "wbpay",
			Amount:        randomInt(randGenerator, 1000, 2000),
			PaymentDT:     randomInt(randGenerator, 10000, 20000),
			Bank:          "alpha",
			DeliveryCost:  randomInt(randGenerator, 1000, 2000),
			GoodsTotal:    randomInt(randGenerator, 100, 500),
			CustomFee:     randomInt(randGenerator, 0, 50),
		},
		Items: []model.Item{
			{
				OrderID:     orderUID,
				ChrtID:      randomInt(randGenerator, 400000, 500000),
				TrackNumber: trackNumber,
				Price:       randomInt(randGenerator, 100, 500),
				Rid:         randomString(randGenerator, 20),
				Name:        "T-Shirt",
				Sale:        randomInt(randGenerator, 10, 50),
				Size:        "M",
				TotalPrice:  randomInt(randGenerator, 100, 500),
				NmID:        randomInt(randGenerator, 100, 200000),
				Brand:       "WB Tech",
				Status:      202,
			},
		},
	}
}

func Publish(natsConn stan.Conn, subject string) error {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGenerator := rand.New(randSource)

	order := randomOrder(randGenerator)
	log.Infoln("Order generated. UID:", order.OrderUID)

	msg, err := json.Marshal(order)
	if err != nil {
		log.Errorf("error marshalling order: %v", err)
		return err
	}
	err = natsConn.Publish(subject, msg)
	if err != nil {
		log.Errorf("error publishing order: %v", err)
		return err
	}
	return nil
}
