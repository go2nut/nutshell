package rel

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type friendDb struct {
	Friends map[string]struct{}
}

var db = &friendDb{make(map[string]struct{}, 0)}

func init() {


	ii := [][]int64{
		{100, 101, 102},
		{103, 104},
		{105, 106, 107, 108, 109, 110, 111},
	}
	for _, i := range ii {
		for _, j := range i {
			for _, j2 := range i {
				k := fmt.Sprintf("%d:%d", j, j2)
				if j > j2 {
					k = fmt.Sprintf("%d:%d", j2, j)
				}
				if j != j2 {
					log.Infof("add friend:%s", k)
					db.Friends[k] = struct{}{}
				}
		    }
		}
	}

}


