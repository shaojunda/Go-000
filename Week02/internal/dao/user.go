package dao

import "../model"

func (d *Dao) Get(id int) (model.User, error) {
	user := model.User{ID: id}
	return user.Get(d.engine)
}
