package myast

import (
	"fmt"
)

type RegexpNode interface {
	Log()
	get() string
}

type Repeat struct {
	node RegexpNode
}

func (n Repeat) get() string {
	return "Repeat(" + n.node.get() + ")"
}

func (n Repeat) Log() {
	fmt.Println(n.get())
}

type Optional struct {
	node RegexpNode
}

func (o Optional) get() string {
	return "Optional(" + o.node.get() + ")"
}

func (o Optional) Log() {
	fmt.Println(o.get())
}

type Concat struct {
	leftNode  RegexpNode
	rightNode RegexpNode
}

func (c Concat) get() string {
	return "Concat(" + c.leftNode.get() + "," + c.rightNode.get() + ")"
}

func (c Concat) Log() {
	fmt.Println(c.get())
}

type Alternate struct {
	leftNode  RegexpNode
	rightNode RegexpNode
}

func (a Alternate) get() string {
	return "Alternate(" + a.leftNode.get() + "," + a.rightNode.get() + ")"
}

func (a Alternate) Log() {
	fmt.Println(a.get())
}

// func (a Alternate) getRegExp() string {
// 	return "(" + a.leftNode.get() + "|" + a.rightNode.get() + ")"
// }

type NegativeCharacterGroup struct {
	node RegexpNode
}

func (n NegativeCharacterGroup) get() string {
	return "NegativeGroup(^" + n.node.get() + ")"
}

func (n NegativeCharacterGroup) Log() {
	fmt.Println(n.get())
}

type PositiveCharacterGroup struct {
	node RegexpNode
}

func (p PositiveCharacterGroup) get() string {
	return "PositiveGroup(" + p.node.get() + ")"
}

func (p PositiveCharacterGroup) Log() {
	fmt.Println(p.get())
}

type Literal struct {
	value rune
}

func (l Literal) get() string {
	// return string("Literal(" + l.value + ")")
	return string(l.value)
}

func (l Literal) Log() {
	fmt.Println(l.get())
}

type Digit struct {
}

func (d Digit) get() string {
	return "AnyDigit()"
}

func (d Digit) Log() {
	fmt.Println(d.get())
}

type AlphaNum struct {
}

func (a AlphaNum) get() string {
	return "AnyAlphaNum()"
}

func (a AlphaNum) Log() {
	fmt.Println(a.get())
}

type AnchorStart struct {
}

func (a AnchorStart) get() string {
	return "AnchorStart()"
}

func (a AnchorStart) Log() {
	fmt.Println(a.get())
}

type AnchorEnd struct {
}

func (a AnchorEnd) get() string {
	return "AnchorEnd()"
}

func (a AnchorEnd) Log() {
	fmt.Println(a.get())
}

type Undefined struct {
}

func (u Undefined) get() string {
	return "Undefined()"
}

func (u Undefined) Log() {
	fmt.Println(u.get())
}
