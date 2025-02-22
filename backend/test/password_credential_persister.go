package test

import (
	"github.com/gofrs/uuid"
	"github.com/teamhanko/hanko/persistence"
	"github.com/teamhanko/hanko/persistence/models"
)

func NewPasswordCredentialPersister(init []models.PasswordCredential) persistence.PasswordCredentialPersister {
	return &passwordCredentialPersister{append([]models.PasswordCredential{}, init...)}
}

type passwordCredentialPersister struct {
	passwords []models.PasswordCredential
}

func (p passwordCredentialPersister) Create(password models.PasswordCredential) error {
	p.passwords = append(p.passwords, password)
	return nil
}

func (p passwordCredentialPersister) GetByUserID(userId uuid.UUID) (*models.PasswordCredential, error) {
	var found *models.PasswordCredential
	for _, data := range p.passwords {
		if data.UserId == userId {
			d := data
			found = &d
		}
	}
	return found, nil
}

func (p passwordCredentialPersister) Update(password models.PasswordCredential) error {
	for i, data := range p.passwords {
		if data.ID == password.ID {
			p.passwords[i] = password
		}
	}
	return nil
}
