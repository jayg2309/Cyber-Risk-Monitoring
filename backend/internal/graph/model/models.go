package model

type AuthPayload struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type CreateAssetInput struct {
	Name      string `json:"name"`
	Target    string `json:"target"`
	AssetType string `json:"assetType"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Asset struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Target        string  `json:"target"`
	AssetType     string  `json:"assetType"`
	CreatedAt     string  `json:"createdAt"`
	LastScannedAt *string `json:"lastScannedAt"`
	Scans         []*Scan `json:"scans"`
}

type Scan struct {
	ID           string        `json:"id"`
	Asset        *Asset        `json:"asset"`
	Status       string        `json:"status"`
	StartedAt    string        `json:"startedAt"`
	CompletedAt  *string       `json:"completedAt"`
	ErrorMessage *string       `json:"errorMessage"`
	Results      []*ScanResult `json:"results"`
}

type ScanResult struct {
	ID       string  `json:"id"`
	Port     int     `json:"port"`
	Protocol string  `json:"protocol"`
	State    string  `json:"state"`
	Service  *string `json:"service"`
	Version  *string `json:"version"`
	Banner   *string `json:"banner"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}
