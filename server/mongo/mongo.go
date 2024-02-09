package mongo

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlayerInfo struct {
	PlayerID string
	Level    int32
	Age      int32
}

func QueryPlayerInfo(playerID string) (*PlayerInfo, error) {
	if playerID == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"请输入用户PlayerID参数",
		)
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
