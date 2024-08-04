package apps

import "github.com/johnnewcombe/telstar/config"

func ShopAddPurchase(sessionId string, settings config.Config, args []string) (bool, error) {

	// this represents the shopping basket
	// use session and DB lookup to get list of products and prices as tabbed columns

	return true, nil
}

func ShopGetPurchases(sessionId string, settings config.Config) string {

	// this represents the shopping basket
	// use session and DB lookup to get list of products and prices as tabbed columns

	return ""
}
