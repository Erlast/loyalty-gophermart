package models

type Model interface {
	Validate() error
}
