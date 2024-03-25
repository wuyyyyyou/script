package service

import (
	"strconv"

	"gorm.io/gorm"

	"github.com/wuyyyyou/script/asset_mapping/models"
	"github.com/wuyyyyou/script/asset_mapping/models/nmap"
)

func Upsert[T models.BaseModeler](db *gorm.DB, oldModel T, newModel T) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var oldModels []T
		err := tx.Where(oldModel).Find(&oldModels).Error
		if err != nil {
			return err
		}

		if len(oldModels) == 0 {
			err = tx.Create(newModel).Error
			if err != nil {
				return err
			}
		} else {
			err = tx.Model(oldModels[0]).Updates(newModel).Error
			if err != nil {
				return err
			}
			newModel.SetID(oldModels[0].GetID())
		}
		return nil
	})
}

func (s *Service) UpsertSeed(seed *models.Seed) error {
	oldSeed := &models.Seed{
		SeedName: seed.SeedName,
	}
	return Upsert(s.DB, oldSeed, seed)
}

func (s *Service) GetAllSeeds() ([]*models.Seed, error) {
	seeds := make([]*models.Seed, 0)
	err := s.DB.Find(&seeds).Error
	return seeds, err
}

func (s *Service) GetAllSeedsWithAll() ([]*models.Seed, error) {
	seeds := make([]*models.Seed, 0)
	err := s.DB.Preload("Domains").Preload("Domains.IPs").Preload("Domains.IPs.Ports").
		Find(&seeds).Error
	return seeds, err
}

func (s *Service) FindSubDomains(seed *models.Seed) ([]*models.Domain, error) {
	domains := make([]*models.Domain, 0)
	err := s.DB.Where("domain like ?", "%"+seed.SeedName+"%").Find(&domains).Error
	return domains, err
}

func (s *Service) AddSeedDomain(seed *models.Seed, domains []*models.Domain) error {
	return s.DB.Model(seed).Association("Domains").Append(domains)
}

func (s *Service) GetAllIPsWithPorts() ([]*models.IP, error) {
	ips := make([]*models.IP, 0)
	err := s.DB.Preload("Ports").Find(&ips).Error
	return ips, err
}

func (s *Service) UpsertIP(ip *models.IP) error {
	oldIP := &models.IP{
		IP: ip.IP,
	}
	return Upsert(s.DB, oldIP, ip)
}

func (s *Service) UpsertIPAndDomain(ip *models.IP, domain *models.Domain) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		err := s.UpsertIP(ip)
		if err != nil {
			return err
		}

		err = s.DB.Model(domain).Association("IPs").Append(ip)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) UpsertIPAndPorts(ip *nmap.Host) error {
	oldIP := &models.IP{
		IP: ip.Address.Addr,
	}
	newIP := &models.IP{
		IP: ip.Address.Addr,
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		err := Upsert(tx, oldIP, newIP)
		if err != nil {
			return err
		}

		for _, port := range ip.Ports.Port {
			portID, err := strconv.Atoi(port.Portid)
			if err != nil {
				return err
			}
			var finger string
			if port.Service.Product != "" {
				finger = port.Service.Product
			} else {
				finger = port.Service.Name
			}

			newPort := &models.Port{
				Port:      portID,
				Protocol:  port.Protocol,
				Finger:    finger,
				BelongsIP: newIP.ID,
			}
			oldPort := &models.Port{
				Port:      portID,
				BelongsIP: newIP.ID,
			}

			err = Upsert(tx, oldPort, newPort)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Service) UpsertSubDomain(domain *models.Domain) error {
	oldDomain := &models.Domain{
		Domain: domain.Domain,
	}
	return Upsert(s.DB, oldDomain, domain)
}

func (s *Service) GetAllDomains() ([]*models.Domain, error) {
	domains := make([]*models.Domain, 0)
	err := s.DB.Find(&domains).Error
	return domains, err
}
