/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-09 20:14
 * description :
 * history :
 */

package dps

import (
	"errors"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/infrastructure/log"
	"go2o/src/core/query"
	"time"
)

type memberService struct {
	_memberRep member.IMemberRep
	_query     *query.MemberQuery
}

func NewMemberService(rep member.IMemberRep, q *query.MemberQuery) *memberService {
	return &memberService{
		_memberRep: rep,
		_query:     q,
	}
}

func (this *memberService) GetMember(id int) *member.ValueMember {
	v, err := this._memberRep.GetMember(id)
	if err == nil {
		nv := v.GetValue()
		return &nv
	}
	return nil
}

func (this *memberService) SaveMember(v *member.ValueMember) (int, error) {
	if v.Id > 0 {
		return this.updateMember(v)
	}
	return this.createMember(v)
}

func (this *memberService) updateMember(v *member.ValueMember) (int, error) {
	m, err := this._memberRep.GetMember(v.Id)
	if err != nil {
		log.PrintErr(err)
		return -1, errors.New("no such member")
	}
	mv := m.GetValue()

	//更新
	mv.Name = v.Name
	mv.Address = v.Address
	mv.Birthday = v.Birthday
	mv.Email = v.Email
	mv.Phone = v.Phone
	mv.Sex = v.Sex
	mv.Qq = v.Qq

	if v.Avatar != "" {
		mv.Avatar = v.Avatar
	}

	m.SetValue(&mv)
	return m.Save()
}

func (this *memberService) createMember(v *member.ValueMember) (int, error) {
	m := this._memberRep.CreateMember(v)
	return m.Save()
}

func (this *memberService) SaveRelation(memberId int, cardId string, tgId, partnerId int) error {
	m, err := this._memberRep.GetMember(memberId)
	if err != nil {
		log.PrintErr(err)
		return errors.New("no such member")
	}

	rl := m.GetRelation()
	rl.TgId = tgId
	rl.Reg_PtId = partnerId
	rl.CardId = cardId

	return m.SaveRelation(rl)
}

func (this *memberService) GetLevel(memberId int) *member.MemberLevel {
	//todo:
	return nil
}

func (this *memberService) GetLevelById(levelValue int) member.MemberLevel {
	return *this._memberRep.GetLevel(levelValue)
}
func (this *memberService) GetNextLevel(levelValue int) *member.MemberLevel {
	return this._memberRep.GetNextLevel(levelValue)
}

func (this *memberService) GetRelation(memberId int) member.MemberRelation {
	return *this._memberRep.GetRelation(memberId)
}

func (this *memberService) Login(usr, pwd string) (bool, *member.ValueMember, error) {
	val := this._memberRep.GetMemberValueByUsr(usr)
	if val == nil {
		return false, nil, errors.New("会员不存在")
	}

	if val.Pwd != domain.EncodeMemberPwd(usr, pwd) {
		return false, nil, errors.New("会员用户或密码不正确")
	}

	if val.State == 0 {
		return false, nil, errors.New("会员已停用")
	}

	val.LastLoginTime = time.Now().Unix()
	this._memberRep.SaveMember(val)

	return true, val, nil
}

func (this *memberService) CheckUsrExist(usr string) bool {
	return this._memberRep.CheckUsrExist(usr)
}

func (this *memberService) GetAccount(memberId int) *member.Account {
	m := this._memberRep.CreateMember(&member.ValueMember{Id: memberId})
	return m.GetAccount()
}

func (this *memberService) GetBank(memberId int) *member.BankInfo {
	m := this._memberRep.CreateMember(&member.ValueMember{Id: memberId})
	b := m.GetBank()
	return &b
}

func (this *memberService) SaveBankInfo(v *member.BankInfo) error {
	m := this._memberRep.CreateMember(&member.ValueMember{Id: v.MemberId})
	return m.SaveBank(v)
}

// 获取返现记录
func (this *memberService) QueryIncomeLog(memberId, page, size int,
	where, orderby string) (num int, rows []map[string]interface{}) {
	return this._query.QueryIncomeLog(memberId, page, size, where, orderby)
}

// 查询分页订单
func (this *memberService) QueryPagerOrder(memberId, page, size int,
	where, orderby string) (num int, rows []map[string]interface{}) {
	return this._query.QueryPagerOrder(memberId, page, size, where, orderby)
}

/*********** 收货地址 ***********/
func (this *memberService) GetDeliverAddrs(memberId int) []member.DeliverAddress {
	return this._memberRep.GetDeliverAddrs(memberId)
}

//获取配送地址
func (this *memberService) GetDeliverAddrById(memberId,
	deliverId int) *member.DeliverAddress {
	m := this._memberRep.CreateMember(&member.ValueMember{Id: memberId})
	v := m.GetDeliver(deliverId).GetValue()
	return &v
}

//保存配送地址
func (this *memberService) SaveDeliverAddr(memberId int, e *member.DeliverAddress) (int, error) {
	m := this._memberRep.CreateMember(&member.ValueMember{Id: memberId})
	var v member.IDeliver
	if e.Id > 0 {
		v = m.GetDeliver(e.Id)
		v.SetValue(e)
	} else {
		v = m.CreateDeliver(e)
	}
	return v.Save()
}

//删除配送地址
func (this *memberService) DeleteDeliverAddr(memberId int, deliverId int) error {
	m := this._memberRep.CreateMember(&member.ValueMember{Id: memberId})
	return m.DeleteDeliver(deliverId)
}

func (this *memberService) ModifyPassword(memberId int, oldPwd, newPwd string) error {
	m, _ := this._memberRep.GetMember(memberId)
	return m.ModifyPassword(newPwd, oldPwd)
}
