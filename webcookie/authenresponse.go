package webcookie

import "encoding/json"

// AuthResponse represents successful authentication response
type AuthResponse struct {
	Success bool       `json:"success"`
	Result  AuthResult `json:"Result"`
	Message string     `json:"Message"`
}

// AuthResult reppresents Result in reponse
type AuthResult struct {
	Version            string          `json:"Version"`
	SessionID          string          `json:"SessionId"`
	AllowLoginMfaCache bool            `json:"AllowLoginMfaCache"`
	Summary            string          `json:"Summary"`
	TenantID           string          `json:"TenantId"`
	Challenges         []AuthChallenge `json:"Challenges"`
}

// AuthChallenge represents list of challenge mchanisims
type AuthChallenge struct {
	Mechanisms []AuthMechanism `json:"Mechanisms"`
}

// AuthMechanism represents authentication mechanism
type AuthMechanism struct {
	AnswerType           string `json:"AnswerType"`           // Text, StartTextOob, StartOob
	Name                 string `json:"Name"`                 // UP, EMAIL, SMS, SQ, PF, OATH
	PartialAddress       string `json:"PartialAddress"`       // For Name = EMAIL
	PartialDeviceAddress string `json:"PartialDeviceAddress"` // For Name = SMS
	Question             string `json:"Question"`             // For Name = SQ
	PartialPhoneNumber   string `json:"PartialPhoneNumber"`   // For Name = PF
	UIPrompt             string `json:"UiPrompt"`             // For AnswerType = StartTextOob

	PromptMechChosen string `json:"PromptMechChosen"`
	PromptSelectMech string `json:"PromptSelectMech"`
	MechanismID      string `json:"MechanismId"`
}

// AdvanceAuthResponse represents successful advance authentication response
type AdvanceAuthResponse struct {
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"Result"`
	Message string                 `json:"Message"`
}

// NewAuthResponse initiates AuthResponse object
func NewAuthResponse(input []byte) (*AuthResponse, error) {
	obj := &AuthResponse{}
	err := json.Unmarshal(input, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// NewAdvanceAuthResponse initiates AdvanceAuthResponse object
func NewAdvanceAuthResponse(input []byte) (*AdvanceAuthResponse, error) {
	obj := &AdvanceAuthResponse{}
	err := json.Unmarshal(input, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
