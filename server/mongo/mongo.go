package mongo

type PlayerInfo struct {
	PlayerID string
	Level    int32
	Age      int32
}

func QueryPlayerInfo(playerID string) (*PlayerInfo, error) {
	if playerID == "" {
		return nil, nil
	}

	if playerID == "leo666" {
		return &PlayerInfo{
			PlayerID: playerID,
			Level:    2,
			Age:      23,
		}, nil
	}

	return &PlayerInfo{
		PlayerID: playerID,
		Level:    1,
		Age:      18,
	}, nil
}
