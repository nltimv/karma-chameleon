package db

import (
	"errors"

	"gorm.io/gorm"
)

func UpdateUserKarma(userID string, teamID string, increment int) (*User, error) {
	var user *User
	var err error
	if user, err = GetUserKarma(userID, teamID); err != nil {
		return user, err
	}

	user.Karma += increment

	if err = db.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetUserKarma(userID string, teamID string) (*User, error) {
	var err error
	var user User
	if err = db.Where("user_id = ? AND team_id = ?", userID, teamID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = User{
				UserId: userID,
				TeamId: teamID,
				Karma:  0,
			}
		} else {
			return &user, err
		}
	}

	return &user, nil
}

func UpdateGroupKarma(groupID string, teamID string, increment int) (*Group, error) {
	var group *Group
	var err error
	if group, err = GetGroupKarma(groupID, teamID); err != nil {
		return group, err
	}

	group.Karma += increment

	if err = db.Save(&group).Error; err != nil {
		return group, err
	}

	return group, nil
}

func GetGroupKarma(groupID string, teamID string) (*Group, error) {
	var err error
	var group Group
	if err = db.Where("user_id = ? AND team_id = ?", groupID, teamID).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			group = Group{
				GroupId: groupID,
				TeamId:  teamID,
				Karma:   0,
			}
		} else {
			return &group, err
		}
	}

	return &group, nil
}
