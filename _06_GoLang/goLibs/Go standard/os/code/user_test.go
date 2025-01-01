package gostd

import (
	"os/user"
	"testing"
)

/*
! LookupGroup, LookupGroupId 按 name 或 groupid 查找组
! user.Group 表示一组用户。在 POSIX 系统上，g.Gid 包含一个表示组 ID 的十进制数
! Lookup, LookupId 按 name 或 userid 查找用户
! user.User 表示用户帐户
	GroupIds 返回用户所属的组 ID 列表
! Current 返回当前用户；第一次调用将缓存当前用户信息
*/

func TestUID(t *testing.T) {
	cur, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	logfln("current user(%s) in %s: \nuid=%s\ngid=%s\nhomeDir=%s\n",
		cur.Username, cur.Name, cur.Uid, cur.Gid, cur.HomeDir)

	ids, err := cur.GroupIds()
	if err != nil {
		t.Fatal(err)
	}
	for _, gid := range ids {
		log(gid)
	}
}
