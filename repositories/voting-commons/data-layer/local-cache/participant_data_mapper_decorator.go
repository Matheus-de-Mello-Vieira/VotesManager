package localdatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"sync"
	"time"
)

type cacheManager struct {
	baseLock         sync.Mutex
	participantsById map[int]domain.Participant
	participants     []domain.Participant
	expiryDate       time.Time
	ttl              time.Duration
}
type ParticipantDataMapperLocalCacheDecorator struct {
	base domain.ParticipantRepository
	cm   *cacheManager
}

func DecorateParticipantRepository(base domain.ParticipantRepository, ttl time.Duration) *ParticipantDataMapperLocalCacheDecorator {
	return &ParticipantDataMapperLocalCacheDecorator{base: base, cm: &cacheManager{ttl: ttl}}
}

func (mapper ParticipantDataMapperLocalCacheDecorator) FindAll(ctx context.Context) ([]domain.Participant, error) {
	mapper.loadCacheIfHaveNotLoaded(ctx)
	return mapper.cm.participants, nil
}
func (mapper ParticipantDataMapperLocalCacheDecorator) FindByID(ctx context.Context, id int) (*domain.Participant, error) {
	mapper.loadCacheIfHaveNotLoaded(ctx)
	participant := mapper.cm.participantsById[id]
	return &participant, nil
}
func (mapper *ParticipantDataMapperLocalCacheDecorator) loadCacheIfHaveNotLoaded(ctx context.Context) error {
	mapper.cm.baseLock.Lock()
	defer mapper.cm.baseLock.Unlock()

	if !mapper.isCacheValid() {
		return mapper.loadCache(ctx)
	}

	return nil
}

func (mapper *ParticipantDataMapperLocalCacheDecorator) isCacheValid() bool {
	return mapper.cm.expiryDate.After(time.Now())
}

func (mapper *ParticipantDataMapperLocalCacheDecorator) loadCache(ctx context.Context) error {
	var err error
	mapper.cm.participants, err = mapper.base.FindAll(ctx)
	if err != nil {
		return err
	}
	mapper.cm.participantsById = assemblyParticipantById(mapper.cm.participants)
	mapper.cm.expiryDate = time.Now().Add(mapper.cm.ttl)

	return nil
}

func assemblyParticipantById(participants []domain.Participant) map[int]domain.Participant {
	result := map[int]domain.Participant{}
	for _, participant := range participants {
		result[participant.ParticipantID] = participant
	}

	return result
}
