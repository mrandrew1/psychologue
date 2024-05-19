package repository

import (
	"github.com/mrandrew1/psychologue/entity"
	"gorm.io/gorm"
)

type PsySessionRepository interface {
	InsertSession(PsySession entity.PsySession) entity.PsySession
	UpdateSession(PsySession entity.PsySession) entity.PsySession
	GetSession(SessionID uint64) entity.PsySession
	DeleteSession(PsySession entity.PsySession)
	AllSessions() []entity.PsySession
}

type PsySessionConnection struct {
	connection *gorm.DB
}

func NewPsySessionRepository(db *gorm.DB) PsySessionRepository {
	return &PsySessionConnection{
		connection: db,
	}
}

func (db *PsySessionConnection) InsertSession(psySession entity.PsySession) entity.PsySession {
	db.connection.Save(&psySession)
	db.connection.Preload("User").Find(&psySession)
	return psySession
}

func (db *PsySessionConnection) UpdateSession(psySession entity.PsySession) entity.PsySession {
	db.connection.Save(&psySession)
	db.connection.Preload("User").Find(&psySession)
	return psySession
}

func (db *PsySessionConnection) GetSession(sessID uint64) entity.PsySession {
	var psySession entity.PsySession
	db.connection.Preload("User").Find(&psySession, sessID)
	return psySession
}

func (db *PsySessionConnection) DeleteSession(psySession entity.PsySession) {
	db.connection.Delete(&psySession)
}

func (db *PsySessionConnection) AllSessions() []entity.PsySession {
	var psySessions []entity.PsySession
	db.connection.Preload("User").Find(&psySessions)
	return psySessions
}
