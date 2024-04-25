package db

import "gorm.io/gorm"

func updateUserKarma(tx *gorm.DB, userID string, teamID string, increment int) (*User, error) {
	var user *User
	var err error
	if user, err = getUserKarma(tx, userID, teamID); err != nil {
		return user, err
	}

	user.Karma += increment

	if err = tx.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func getUserKarma(tx *gorm.DB, userID string, teamID string) (*User, error) {
	var user *User
	//err := tx.Attrs(User{Karma: 0}).FirstOrInit(user, User{UserId: userID, TeamId: teamID}).Error
	err := tx.First(&user, User{UserId: userID, TeamId: teamID}).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		user = &User{UserId: userID, TeamId: teamID, Karma: 0}
		return user, nil
	}

	return user, err
}

func updateGroupKarma(tx *gorm.DB, groupID string, teamID string, increment int) (*Group, error) {
	var group *Group
	var err error
	if group, err = getGroupKarma(tx, groupID, teamID); err != nil {
		return group, err
	}

	group.Karma += increment

	if err = tx.Save(&group).Error; err != nil {
		return group, err
	}

	return group, nil
}

func getGroupKarma(tx *gorm.DB, groupID string, teamID string) (*Group, error) {
	var group *Group
	err := tx.First(group, Group{GroupId: groupID, TeamId: teamID}).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		group = &Group{GroupId: groupID, TeamId: teamID, Karma: 0}
		return group, nil
	}

	return group, err
}
