package main

import (
	"./internal/service"
	"fmt"
)

/*
	题目：我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
	解答：dao 层遇到 sql.ErrNoRows 后返回一个自定义的 error ErrRecordNotFound，直接抛给上层。因为底层的 error 可能会变化比如当换一个 ORM 实现时，所以这里选择返回一个自定义的 Error。
*/

func main() {
	svc := service.New()
	user, err := svc.GetUser(1)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
	} else {
		fmt.Printf("user: %+v\n", user)
	}
}
