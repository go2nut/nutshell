package rel



type friendDb struct {
	Users map[int64]string
	Friends map[int64][]int64
}

var db = &friendDb{make(map[int64]string, 0), make(map[int64][]int64, 0)}

func init() {
	db.Users[100] = "John"
	db.Users[101] = "Mik"
	db.Users[102] = "Trump"
	db.Users[103] = "Bidden"
	db.Users[104] = "Jobs"
	db.Users[105] = "Steven"
	db.Users[105] = "Lisa"
	db.Users[106] = "Nasa"
	db.Users[107] = "Lucy"
	db.Users[108] = "Joe"
	db.Users[109] = "Zoe"
	db.Users[110] = "Dave"

	for i:= int64(102); i<= 108; i++ {
		newFrs := []int64{i-2, i-1, i+1, i+2}
		frs, exist := db.Friends[i]
		if exist {
			frs = append(newFrs, []int64{i-2, i-1, i+1, i+2}...)
		} else {
			frs = newFrs
		}
		db.Friends[i] = frs

		for _, newFrId := range newFrs {
			oppoNewFrs, oppoExist := db.Friends[newFrId]
			if oppoExist {
				oppoNewFrs = append(oppoNewFrs, newFrId)
			} else {
				oppoNewFrs = []int64{newFrId}
			}
			db.Friends[newFrId] = oppoNewFrs
		}
	}

}


