
package main

import (
	//	"fmt"
	"github.com/Maginawork/shortService/src/shortlib"
)

func CountThread(count_chan_in chan shortlib.CountChannl) {

	var count int64
	count = 1000
	for {
		select {
		case ok := <-count_chan_in:
			count = count + 1
			ok.CountOutChan <- count

		}

	}
}
