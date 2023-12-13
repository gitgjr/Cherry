package nc

import (
	"fmt"

	"github.com/eugene-eeo/vivaldi-go"
)

func NcExample() {
	local := vivaldi.NewContext()
	remote := vivaldi.NewHVector(
		23.0, // x
		45.0, // y
		10.0, // height
	)
	rtt := 5.0
	local.Update(rtt, vivaldi.NewContextFromValues(
		remote,
		5.0, // error estimate
	))
	predicts, err := local.EstimateRTT(remote)
	if err != nil {
		panic(err)
	}
	fmt.Println(predicts)
}
