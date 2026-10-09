// Package tools: fungsi utilitas murni (nggak tau soal WhatsApp), gampang dites.
package tools

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"wabot/config"
)

var start = time.Now()

func Uptime() string {
	d := time.Since(start).Round(time.Second)
	return fmt.Sprintf("%dh %dm %ds", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
}

// Clock = jam lokal sesuai TZ_OFFSET_HOURS.
func Clock() string {
	off := config.Get().TZOffsetHours
	t := time.Now().UTC().Add(time.Duration(off) * time.Hour)
	return fmt.Sprintf("%s (UTC%+d)", t.Format("15:04:05"), off)
}

// Calc mengevaluasi ekspresi aritmetika: + - * / % ^ dan kurung. Tanpa dependency.
func Calc(expr string) (string, error) {
	if strings.TrimSpace(expr) == "" {
		return "", errors.New("Contoh: calc 2*(3+4)^2")
	}
	p := &parser{s: expr}
	v, err := p.sum()
	if err != nil {
		return "", err
	}
	if p.skip(); p.i < len(p.s) {
		return "", fmt.Errorf("Ekspresi salah di %q", p.s[p.i:])
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "", errors.New("Hasil tidak valid")
	}
	return strconv.FormatFloat(v, 'g', 12, 64), nil
}

type parser struct {
	s string
	i int
}

func (p *parser) skip() {
	for p.i < len(p.s) && p.s[p.i] == ' ' {
		p.i++
	}
}

func (p *parser) eat(c byte) bool {
	p.skip()
	if p.i < len(p.s) && p.s[p.i] == c {
		p.i++
		return true
	}
	return false
}

// sum := prod (('+'|'-') prod)*
func (p *parser) sum() (float64, error) {
	v, err := p.prod()
	for err == nil {
		switch {
		case p.eat('+'):
			var r float64
			if r, err = p.prod(); err == nil {
				v += r
			}
		case p.eat('-'):
			var r float64
			if r, err = p.prod(); err == nil {
				v -= r
			}
		default:
			return v, nil
		}
	}
	return 0, err
}

// prod := unary (('*'|'/'|'%') unary)*
func (p *parser) prod() (float64, error) {
	v, err := p.unary()
	for err == nil {
		var op byte
		switch {
		case p.eat('*'):
			op = '*'
		case p.eat('/'):
			op = '/'
		case p.eat('%'):
			op = '%'
		default:
			return v, nil
		}
		var r float64
		if r, err = p.unary(); err != nil {
			break
		}
		if r == 0 && op != '*' {
			return 0, errors.New("Pembagian dengan nol")
		}
		switch op {
		case '*':
			v *= r
		case '/':
			v /= r
		default:
			v = math.Mod(v, r)
		}
	}
	return 0, err
}

// unary := '-' unary | '+' unary | pow       (jadi -2^2 = -(2^2))
func (p *parser) unary() (float64, error) {
	if p.eat('-') {
		v, err := p.unary()
		return -v, err
	}
	if p.eat('+') {
		return p.unary()
	}
	return p.pow()
}

// pow := atom ('^' unary)?                   (right-associative)
func (p *parser) pow() (float64, error) {
	b, err := p.atom()
	if err != nil || !p.eat('^') {
		return b, err
	}
	e, err := p.unary()
	return math.Pow(b, e), err
}

func (p *parser) atom() (float64, error) {
	if p.eat('(') {
		v, err := p.sum()
		if err != nil {
			return 0, err
		}
		if !p.eat(')') {
			return 0, errors.New("Kurung tutup kurang")
		}
		return v, nil
	}
	p.skip()
	j := p.i
	for j < len(p.s) && (p.s[j] >= '0' && p.s[j] <= '9' || p.s[j] == '.') {
		j++
	}
	if j == p.i {
		return 0, errors.New("Ekspresi salah, contoh: 2*(3+4)^2")
	}
	v, err := strconv.ParseFloat(p.s[p.i:j], 64)
	p.i = j
	if err != nil {
		return 0, errors.New("Angka tidak valid")
	}
	return v, nil
}
