package db

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
	var user *User
	err := db.Attrs(User{Karma: 0}).FirstOrInit(user, User{UserId: userID, TeamId: teamID}).Error

	return user, err
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
	var group *Group
	err := db.Attrs(Group{Karma: 0}).FirstOrInit(group, Group{GroupId: groupID, TeamId: teamID}).Error

	return group, err
}
