package models

import "errors"

type Goods struct {
	Match      string `json:"match"`
	Reward     int64  `json:"reward"`
	RewardType string `json:"reward_type"`
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
