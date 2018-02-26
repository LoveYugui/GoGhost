package user

import (
	"testing"
)

func TestNewUser(t *testing.T)  {

	user := NewUser("weixin", "LoveYugui1", "yugui", "")

	if user == nil {
		t.Errorf("error TestNewUser")
	} else {
		t.Errorf("user : %v", *user)
	}

}
