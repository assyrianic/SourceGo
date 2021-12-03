package main

import "sourcemod"

const SLOTCOUNT = 3

type Player struct {
	userID, pauseTick, buttons int
	origin, velocity Vec3
	nextAttackPrimary, nextAttackSecondary [SLOTCOUNT]float
	stasisTick float
	isInStasis bool
}

func (p *Player) get() int {
	return GetClientOfUserId(p.userID)
}

func (p *Player) set(client Entity) {
	if isValidClient(client) {
		p.userID = GetClientUserid(client)
	} else {
		p.userID = 0
	}
}

func (p *Player) toggleStasis() {
	if !p.isInStasis {
		p.enableStasis()
	} else {
		p.disableStasis()
	}
}

func (p *Player) enableStasis() {
	p.isInStasis = true
	p.freeze()
	p.pauseAttack()
	p.displayLasers()
}

func (p *Player) disableStasis() {
	p.isInStasis = false
	p.unfreeze()
	p.resumeAttack()
}

func (p *Player) freeze() {
	client := p.get()
	if !client {
		return
	}

	freezeProjectiles(client)

	GetClientAbsOrigin(client, p.origin)

	getClientAbsVelocity(client, p.velocity)

	SetEntityMoveType(client, MOVETYPE_NONE)

	p.stasisTick = GetGameTime()
	p.buttons = GetClientButtons(client)
	p.pauseTick = GetGameTickCount()
}

func (p *Player) unfreeze() {
	client := p.get()
	if !client {
		return
	}

	unfreezeProjectiles(client)
	SetEntityMoveType(client, MOVETYPE_WALK)
	TeleportEntity(client, p.origin, NULL_VECTOR, p.velocity)
}

func (p *Player) pauseAttack() {
	client := p.get()
	if !client {
		return
	}

	for slot := 0; slot < SLOTCOUNT; slot++ {
		var weapon int = GetPlayerWeaponSlot(client, slot)
		if !IsValidEntity(weapon) {
			continue
		}

		if HasEntProp(weapon, Prop_Send, "m_flNextPrimaryAttack") {
			p.nextAttackPrimary[slot] = GetEntPropFloat(weapon, Prop_Send, "m_flNextPrimaryAttack")
			SetEntPropFloat(weapon, Prop_Send, "m_flNextPrimaryAttack", 999999999.0)
		}

		if HasEntProp(weapon, Prop_Send, "m_flNextSecondaryAttack") {
			p.nextAttackSecondary[slot] = GetEntPropFloat(weapon, Prop_Send, "m_flNextSecondaryAttack")
			SetEntPropFloat(weapon, Prop_Send, "m_flNextSecondaryAttack", 999999999.0)
		}
	}
}

func (p *Player) resumeAttack() {
	client := p.get()
	if !client {
		return
	}

	for slot := 0; slot < SLOTCOUNT; slot++ {
		var weapon int = GetPlayerWeaponSlot(client, slot)
		if !IsValidEntity(weapon) {
			continue
		}
		
		if HasEntProp(weapon, Prop_Send, "m_flNextPrimaryAttack") {
			SetEntPropFloat( weapon, Prop_Send, "m_flNextPrimaryAttack", p.nextAttackPrimary[slot]+(GetGameTime()-p.stasisTick))
		}
		if HasEntProp(weapon, Prop_Send, "m_flNextSecondaryAttack") {
			SetEntPropFloat(weapon, Prop_Send, "m_flNextSecondaryAttack", p.nextAttackSecondary[slot]+(GetGameTime()-p.stasisTick))
		}
	}
}

func (p *Player) displayLasers() {
	client := p.get()
	if !client {
		return
	}

	var angles, eyePos, temp, fwrd, end Vec3
	GetClientEyeAngles(client, angles)
	GetClientEyePosition(client, eyePos)
	GetVectorAngles(p.velocity, &temp)
	
	GetAngleVectors(temp, &fwrd, &NULL_VECTOR, &NULL_VECTOR)
	ScaleVector(&fwrd, GetVectorLength(p.velocity)*0.2)
	AddVectors(p.origin, fwrd, &temp)
	doLaserBeam(client, p.origin, temp)

	GetAngleVectors(angles, &fwrd, &NULL_VECTOR, &NULL_VECTOR)
	ScaleVector(&fwrd, 80.0)

	AddVectors(eyePos, fwrd, &end)
	doLaserBeam(client, eyePos, end, 255, 20, 20)

	temp = eyePos
	temp[2] -= 30.0
	doLaserBeam(client, temp, end, 255, 20, 20)
}
