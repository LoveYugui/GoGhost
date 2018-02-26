package Role

import "github.com/GoGhost/MyLord/user"
import "github.com/GoGhost/persistence/redisConn"

const (
	RoleFarmer = iota
	RoleLord
)

type UserRole struct {
	GameUser user.User
	RoleType uint8
	Cards []CardRole

}

func NewUserRole() *UserRole  {
	return &UserRole{
		Cards:make([]CardRole, 17),
	}
}

type OnlineStatus struct {
	UserId string
	Status int8
}
func UserOnlineStatus(uids []string) []OnlineStatus {
	redisConn.RedisGet()
}
