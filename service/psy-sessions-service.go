package service

import (
	"fmt"
	"log"

	"github.com/mashingan/smapping"
	"github.com/putukrisna6/golang-api/dto"
	"github.com/putukrisna6/golang-api/entity"
	"github.com/putukrisna6/golang-api/repository"
)

type PsySessionService interface {
	Insert(b dto.PsySessionCreateDTO) entity.PsySession
	Update(b dto.PsySessionUpdateDTO) entity.PsySession
	Get(bookID uint64) entity.PsySession
	Delete(b entity.PsySession)
	All() []entity.PsySession
	IsAllowedToEdit(userID string, sessID uint64) bool
}

type psySessionService struct {
	psySessionRepository repository.PsySessionRepository
}

func NewPsySessionService(psySessionRepository repository.PsySessionRepository) PsySessionService {
	return &psySessionService{
		psySessionRepository: psySessionRepository,
	}
}

func (service *psySessionService) Insert(dto dto.PsySessionCreateDTO) entity.PsySession {
	newSession := entity.PsySession{}
	err := smapping.FillStruct(&newSession, smapping.MapFields(&dto))
	if err != nil {
		log.Fatalf("failed to map: %v", err.Error())
	}

	res := service.psySessionRepository.InsertSession(newSession)
	return res
}

func (service *psySessionService) Update(dto dto.PsySessionUpdateDTO) entity.PsySession {
	PsySession := entity.PsySession{}
	err := smapping.FillStruct(&PsySession, smapping.MapFields(&dto))
	if err != nil {
		log.Fatalf("failed to map: %v", err.Error())
	}

	res := service.psySessionRepository.UpdateSession(PsySession)
	return res
}

func (service *psySessionService) Get(sessID uint64) entity.PsySession {
	return service.psySessionRepository.GetSession(sessID)
}

func (service *psySessionService) All() []entity.PsySession {
	return service.psySessionRepository.AllSessions()
}

func (service *psySessionService) Delete(session entity.PsySession) {
	service.psySessionRepository.DeleteSession(session)
}

func (service *psySessionService) IsAllowedToEdit(userID string, sessID uint64) bool {
	session := service.psySessionRepository.GetSession(sessID)
	return userID == fmt.Sprintf("%v", session.SpecialistID)
}
