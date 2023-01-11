package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Input struct {
	ID                                   uint `gorm:"primaryKey"`
	CreatedAt                            time.Time
	UpdatedAt                            time.Time
	DeletedAt                            gorm.DeletedAt `gorm:"index"`
	Tag                                  string         `json:"tag"`
	Name                                 string         `json:"name"`
	NameColor                            string         `json:"nameColor"`
	Trophies                             uint           `json:"trophies"`
	HighestTrophies                      uint           `json:"highestTrophies"`
	HighestPowerPlayPoints               uint           `json:"highestPowerPlayPoints"`
	ExpLevel                             uint           `json:"expLevel"`
	ExpPoints                            uint           `json:"expPoints"`
	IsQualifiedFromChampionshipChallenge bool           `json:"isQualifiedFromChampionshipChallenge"`
	TeamVictories                        uint           `json:"3vs3Victories"`
	SoloVictories                        uint           `json:"soloVictories"`
	DuoVictories                         uint           `json:"duoVictories"`
	BestRoboRumbleTime                   uint           `json:"bestRoboRumbleTime"`
	BestTimeAsBigBrawler                 uint           `json:"bestTimeAsBigBrawler"`
}

type User struct {
	ID                                   uint           `gorm:"primaryKey" json:"id"`
	CreatedAt                            time.Time      `json:"created_at"`
	UpdatedAt                            time.Time      `json:"updated_at"`
	DeletedAt                            gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Tag                                  string         `json:"tag"`
	Name                                 string         `json:"name"`
	NameColor                            string         `json:"name_color"`
	Trophies                             uint           `json:"trophies"`
	HighestTrophies                      uint           `json:"highest_trophies"`
	HighestPowerPlayPoints               uint           `json:"highest_power_playPoints"`
	ExpLevel                             uint           `json:"exp_level"`
	ExpPoints                            uint           `json:"exp_points"`
	IsQualifiedFromChampionshipChallenge bool           `json:"is_qualified_from_championship_challenge"`
	TeamVictories                        uint           `json:"team_victories"`
	SoloVictories                        uint           `json:"solo_victories"`
	DuoVictories                         uint           `json:"duo_victories"`
	BestRoboRumbleTime                   uint           `json:"best_robo_rumble_time"`
	BestTimeAsBigBrawler                 uint           `json:"best_time_as_big_brawler"`
}

type UserService struct {
	DB *gorm.DB
}

func (us *UserService) Migrate() error {
	us.DB.AutoMigrate(&User{})
	return nil
}

func (us *UserService) Create(newU *User) error {
	result := us.DB.Create(newU)
	if result.Error != nil {
		fmt.Println(result.Error)
		return fmt.Errorf("an error occurred while creating")
	}
	return nil
}

func (us *UserService) FindOrCreate(newU *User) error {
	result := us.DB.Where(User{Tag: newU.Tag}).FirstOrCreate(newU)
	if result.Error != nil {
		fmt.Println(result.Error)
		return fmt.Errorf("an error occurred while creating")
	}
	return nil
}

func (us *UserService) CreateOrUpdate(newU *User) error {
	if us.DB.Model(newU).Clauses(clause.Returning{}).Where(User{Tag: newU.Tag}).Updates(newU).RowsAffected == 1 {
		fmt.Println("Updated successfully")
		return nil
	}

	result := us.DB.Create(newU)
	if result.Error != nil {
		fmt.Println(result.Error)
		return fmt.Errorf("an error occurred while creating")
	}
	return nil
}
