package block

type block struct {
	Reward        int `json:"rewardamount"`
	Confirmations int `json:"confirmations"`
}

/*
	TODO:
		- which one has less mutation, with successful HVEC encoding
*/
