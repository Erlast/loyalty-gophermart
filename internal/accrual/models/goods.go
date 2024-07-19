package models

import "errors"

type Goods struct {
	RewardType string  `json:"reward_type"`
	Match      string  `json:"match"`
	Reward     float32 `json:"reward"`
}

func (o *Goods) Validate() error {
	if o.Match == "" {
		return errors.New("match is required")
	}

	if o.Reward <= 0 {
		return errors.New("reward is required")
	}

	if o.RewardType != "%" && o.RewardType != "pt" {
		return errors.New("reward type is invalid")
	}
	return nil
}
