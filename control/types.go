package control

type Person struct {
	Name       string
	Department *Department
	Score      int
}

type Department struct {
	Name         string
	NumberPeople int
}

type Team struct {
	Members []*Person
	Score   int
}

type Match struct {
	Match [2]*Person
	Score int
}

type Leader struct {
	Person     *Person
	TotalScore int
	Index      int
}

type Departments []*Department

func (slice Departments) Len() int {
	return len(slice)
}

func (slice Departments) Less(i, j int) bool {
	return slice[i].NumberPeople < slice[j].NumberPeople
}

func (slice Departments) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type People []*Person

func (slice People) Len() int {
	return len(slice)
}

func (slice People) Less(i, j int) bool {
	return slice[i].Score < slice[j].Score
}

func (slice People) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Matches []*Match

func (slice Matches) Len() int {
	return len(slice)
}

func (slice Matches) Less(i, j int) bool {
	return slice[i].Score < slice[j].Score
}

func (slice Matches) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Leaderboard []*Leader

func (slice Leaderboard) Len() int {
	return len(slice)
}

func (slice Leaderboard) Less(i, j int) bool {
	return slice[i].TotalScore < slice[j].TotalScore
}

func (slice Leaderboard) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
