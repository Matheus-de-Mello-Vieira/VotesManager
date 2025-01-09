package domain

type VoteRepository interface {
	SaveOne(vote *Vote) error
	SaveMany(votes []Vote) error
}