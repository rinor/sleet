package test

import (
	"fmt"
	"github.com/BoltApp/sleet/gateways/authorize_net"
	"math/rand"
	"testing"
	"time"

	"github.com/BoltApp/sleet"
)

func TestAuthNet(t *testing.T) {
	client := authorize_net.NewClient(getEnv("AUTH_NET_LOGIN_ID"), getEnv("AUTH_NET_TXN_KEY"))
	rand.Seed(time.Now().Unix())
	randAmount := rand.Int63n(1000000)
	amount := sleet.Amount{
		Amount:   randAmount,
		Currency: "USD",
	}
	postalCode := "94103"
	address := sleet.BillingAddress{PostalCode: &postalCode}
	card := sleet.CreditCard{
		FirstName:       "Bolt",
		LastName:        "Checkout",
		Number:          "4111111111111111",
		ExpirationMonth: 8,
		ExpirationYear:  2024,
		CVV:             "111",
	}
	resp, err := client.Authorize(&sleet.AuthorizationRequest{Amount: &amount, CreditCard: &card, BillingAddress: &address})
	fmt.Printf("resp: [%+v] err [%s]\n", resp, err)

	capResp, err := client.Capture(&sleet.CaptureRequest{
		Amount:               &amount,
		TransactionReference: resp.TransactionReference,
	})
	fmt.Printf("capResp: [%+v] err [%s]\n", capResp, err)

	lastFour := card.Number[len(card.Number)-4:]
	options := make(map[string]interface{})
	options["credit_card"] = lastFour
	refundResp, err := client.Refund(&sleet.RefundRequest{
		Amount:               &amount,
		TransactionReference: resp.TransactionReference,
		Options:              options,
	})
	fmt.Printf("refundResp: [%+v] err [%s]\n", refundResp, err)
}