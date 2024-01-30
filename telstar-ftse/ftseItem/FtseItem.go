package ftseItem

import (
	"errors"
	"strconv"
	"strings"
)

type FtseItem struct { // code, name, currency, market cap, price, change, change %

	Code          string
	Name         string
//	CurrentPrice string
//	MarketCap    float64
	Price         float64
	Change        float64
	ChangePerCent float64
}

func (f *FtseItem) Load(data string) error {

	var (
		err error
	)

	if len(data)==0{
		return nil
	}

	item := strings.Split(data, "|") // Code, Name, Current Price, Day Change, Percentage Change
	if len(item) != 6 {
		return errors.New("data is invalid")
	}

	f.Code = item[0]
	f.Name = item[1]

	if f.Price, err = strconv.ParseFloat(strings.Replace(item[2],",","", -1), 64); err != nil {
		return err
	}
	if f.Change, err = strconv.ParseFloat(strings.Replace(item[3],",","", -1), 64); err != nil {
		return err
	}
	item[4] = strings.Replace(item[4],"%","", -1)
	if f.ChangePerCent, err = strconv.ParseFloat(strings.Replace(item[4],",","", -1), 64); err != nil {
		return err
	}

	return nil
}
