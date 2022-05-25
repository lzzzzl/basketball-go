package stats

type Stats struct {
	Assits    int
	Points    int
	Blocks    int
	Steals    int
	Turnovers int
	Rebounds  int
}

type PlayerStats struct {
	Name                    string  `json:"name"`
	Position                string  `json:"position"`
	Status                  string  `json:"status"`
	Minutes                 string  `json:"minutes"`
	JerseyNum               string  `json:"jeseynum"`
	IsStarter               bool    `json:"isstarter"`
	Assits                  int     `json:"assits"`
	Points                  int     `json:"points"`
	Blocks                  int     `json:"blocks"`
	Steals                  int     `json:"steals"`
	Turnovers               int     `json:"turnovers"`
	Rebounds                int     `json:"rebounds"`
	PlusMinusPoints         string  `json:"plusminuspoints"`
	FieldGoalsAttempted     int     `json:"fieldgoalsattempted"`
	FieldGoalsMade          int     `json:"fieldgoalsmade"`
	FieldGoalsPercentage    float64 `json:"fieldgoalspercentage"`
	ThreePointersAttempted  int     `json:"threepointsattempted"`
	ThreePointersMade       int     `json:"threepointsmade"`
	ThreePointersPercentage float64 `json:"threepointspercentage"`
	FreeThrowsAttempted     int     `json:"freethrowsattempted"`
	FreeThrowsMade          int     `json:"freethrowsmade"`
	FreeThrowsPercentage    float64 `json:"freethrowspercentage"`
}
