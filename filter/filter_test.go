// Copyright 2015 Eleme Inc. All rights reserved.

package filter

import (
	"fmt"
	"testing"
	"time"

	"github.com/eleme/banshee/models"
)

func TestFilter(t *testing.T) {
	//	filter:=NewFilter()
	//	words:=[]string{"*"}
	//	for i:=1;i<1000;i++{
	//		words=append(words,string(i))
	//	}
	//	now:=time.Now()
	//	x:=4+rand.Intn(3)
	//	filter.AddRule <-
	filter := NewFilter()

	filter.AddRule <- &models.Rule{Pattern: "5.93.355.*.22"}
	filter.AddRule <- &models.Rule{Pattern: "6.123.45.*"}
	now := time.Now()
	for i := 0; i < 1000000; i++ {
		filter.MatchedRules(&models.Metric{Name: "5.23.71.35"})
	}
	fmt.Println(time.Since(now).Nanoseconds() / 1000000)
	for i := 0; i < 1000000; i++ {
		filter.MatchedRules(&models.Metric{Name: "6.123.45.*"})
	}
	fmt.Println(time.Since(now).Nanoseconds() / 1000000)
}
