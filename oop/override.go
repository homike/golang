package oop

import "log"

type People struct {
}

func (p *People) PrePing() {
	log.Println("pre ping")
}

func (p *People) Ping() {
	log.Println("ping")
}

func (p *People) Start() {
	p.PrePing()

	p.Ping()
}

type Man struct {
	People
}

func (g *Man) Ping() {
	log.Println("pong")
}

func RunOOP() {
	m := &Man{}
	//m.PrePing()
	//m.Ping()
	m.Start()
}
