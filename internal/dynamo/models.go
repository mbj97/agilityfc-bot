package dynamo

type User struct {
	UserID   string `json:"UserID"`
	LastSeen string `json:"LastSeen"`
	World    string `json:"World"`
}

type Snapshot struct {
	ID          string      `json:"ID"`
	Timestamp   string      `json:"Timestamp"`
	Snapshot    []UserEntry `json:"Snapshot"`
	NumberOnline int        `json:"NumberOnline"`
}

type UserEntry struct {
    UserID string `json:"UserID"`
    World  string `json:"World"`
}
