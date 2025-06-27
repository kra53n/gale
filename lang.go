package main

import (
	_ "path/filepath"
	"strconv"
	"strings"
)

// #cgo CFLAGS: -Ivendor
// #cgo LDFLAGS: -lm
// #include "gale.h"
// import "C"

var program string = `crop 10% 20%
> imgs/картинк*
> scan/imgs/here/*

rotate left
rotate right`

const (
	ImgFormatNotSupported = iota
	ImgFormatJpg
	ImgFormatPng
	ImgFormatBmp
	ImgFormatTga
	ImgFormatPsd
	ImgFormatHdr
	ImgFormatGif
	ImgFormatPnm
)

type LexState int

const (
	LexStateNil LexState = iota
	LexStateScanFilepath
	LexStateScanWord
)

type LexWordKind int

const (
	LexBaseWord LexWordKind = iota
	LexFilepathWord
	LexNewLWord
)

type LexWord struct {
	beg, end int
	kind     LexWordKind
}

func (w LexWord) String() string {
	return strings.TrimSpace(program[w.beg:w.end])
}

type Words []LexWord

var words Words

func (w Words) iterator() func() (LexWord, bool) {
	var i int
	return func() (LexWord, bool) {
		if i < len(w) {
			i++
			return w[i-1], true
		}
		return LexWord{}, false
	}
}

type RuleKind int

const (
	NilRule RuleKind = iota
	Crop
	RotateLeft
	RotateRight
	FlipVertically
	FlipHorizontally
)

type MeasurementUnit int

const (
	Pixel MeasurementUnit = iota
	Percent
)

type CropParam struct {
	value int
	unit  MeasurementUnit
}

type CropParams struct {
	top, right, bottom, left CropParam
}

type Rule struct {
	kind   RuleKind
	params interface{}
}

var generalRules []Rule

type Img string

type Block struct {
	imgs  []Img
	rules []Rule
}

var blocks []Block

func main() {
	lex()
	parse()
	exec()
}

func exec() {
	for _, block := range blocks {
		block.applyRules()
	}
}

func (b *Block) applyRules() {
	for _, img := range b.imgs {
		img.open()
		img.applyRules(generalRules)
		img.applyRules(b.rules)
		img.close()
	}
}

func (i *Img) open() {

}

func (i *Img) close() {

}

func (i *Img) applyRules(rules []Rule) {
	for _, r := range rules {
		switch r.kind {
		case Crop:
			if params, ok := r.params.(CropParams); ok {
				i.crop(params)
			}
		case RotateLeft:
			i.rotateLeft()
		case RotateRight:
			i.rotateRight()
		case FlipVertically:
			i.flipVertically()
		case FlipHorizontally:
			i.flipHorizontally()
		}
	}
}

func (i *Img) crop(params CropParams) {

}

func (i *Img) rotateLeft() {

}

func (i *Img) rotateRight() {

}

func (i *Img) flipVertically() {

}

func (i *Img) flipHorizontally() {

}

func lex() {
	var s LexState
	var w LexWord
	for i, r := range program {
		switch r {
		case '\r', '\n':
			if s == LexStateScanWord || s == LexStateScanFilepath {
				w.end += 1
				words = append(words, w)
			}
			s = LexStateNil
			if len(words) > 0 && words[len(words)-1].kind != LexNewLWord {
				words = append(words, LexWord{kind: LexNewLWord})
			}
		case ' ', '\t':
			if s == LexStateScanWord {
				w.end += 1
				words = append(words, w)
			}
			if s == LexStateScanFilepath {
				w.end += 1
				continue
			}
			s = LexStateNil
		case '>':
			if s != LexStateNil {
				panic("error")
			}
			s = LexStateScanFilepath
			w = LexWord{beg: i + 1, end: i, kind: LexFilepathWord}
		default:
			if s == LexStateScanWord || s == LexStateScanFilepath {
				w.end += 1
			} else {
				w = LexWord{beg: i, end: i}
				s = LexStateScanWord
			}
		}
	}
	if s == LexStateScanFilepath || s == LexStateScanWord {
		w.end += 1
		words = append(words, w)
	}
	// for _, word := range words {
	// 	fmt.Println(program[word.beg:word.end])
	// }
}

func parse() {
	var generalRulesWasScanned bool
	var rule Rule
	var block Block
	iter := words.iterator()
	for w, ok := iter(); ok; w, ok = iter() {
		s := w.String()
		switch w.kind {
		case LexBaseWord:
			switch s {
			case "crop", "c":
				rule = parseCrop(iter)
			case "rotate":
				if w, ok = iter(); !ok || w.kind != LexBaseWord {
					panic("bad")
				}
				switch w.String() {
				case "right":
					rule = Rule{
						kind: RotateRight,
					}
				case "left":
					rule = Rule{
						kind: RotateLeft,
					}
				default:
					panic("bad")
				}
			case "flip":
				if w, ok = iter(); !ok || w.kind != LexBaseWord {
					panic("bad")
				}
				switch w.String() {
				case "vertically":
					rule = Rule{
						kind: FlipVertically,
					}
				case "horizontally":
					rule = Rule{
						kind: FlipHorizontally,
					}
				default:
					panic("bad")
				}
			case "l":
				rule = Rule{kind: RotateLeft}
			case "r":
				rule = Rule{kind: RotateRight}
			case "v":
				rule = Rule{kind: FlipVertically}
			case "h":
				rule = Rule{kind: FlipHorizontally}
			}
		case LexFilepathWord:
			generalRulesWasScanned = true
			if len(block.rules) == 0 {
				block.imgs = append(block.imgs, Img(s))
			} else {
				blocks = append(blocks, block)
			}
			continue
		case LexNewLWord:
			continue
		}
		if generalRulesWasScanned {
			block.rules = append(block.rules, rule)
		} else {
			generalRules = append(generalRules, rule)
		}
	}
	blocks = append(blocks, block)
}

func parseCrop(iter func() (LexWord, bool)) Rule {
	var w LexWord
	var ok bool
	ws := make([]LexWord, 0, 5)
	for w, ok = iter(); ok && w.kind == LexBaseWord; w, ok = iter() {
		ws = append(ws, w)
	}
	if len(ws) == 0 || len(ws) == 3 || len(ws) >= 5 {
	}
	switch len(ws) {
	case 1:
		return Rule{
			kind: Crop,
			params: CropParams{
				top:    parseCropParam(ws[0]),
				right:  parseCropParam(ws[0]),
				bottom: parseCropParam(ws[0]),
				left:   parseCropParam(ws[0]),
			},
		}
	case 2:
		return Rule{
			kind: Crop,
			params: CropParams{
				top:    parseCropParam(ws[0]),
				right:  parseCropParam(ws[1]),
				bottom: parseCropParam(ws[0]),
				left:   parseCropParam(ws[1]),
			},
		}
	case 4:
		return Rule{
			kind: Crop,
			params: CropParams{
				top:    parseCropParam(ws[0]),
				right:  parseCropParam(ws[1]),
				bottom: parseCropParam(ws[2]),
				left:   parseCropParam(ws[3]),
			},
		}
	case 0:
		fallthrough
	case 3:
		fallthrough
	default:
		panic("wrong numbers for crop action, possible only 1, 2 or 4")
	}
	return Rule{}
}

func parseCropParam(w LexWord) CropParam {
	var unit MeasurementUnit
	var value int
	var digitEndScanning int
	s := w.String()
	for digitEndScanning < len(s) {
		c := s[digitEndScanning]
		if !('0' <= c && c <= '9') {
			break
		}
		digitEndScanning++
	}
	if digitEndScanning == 0 {
		panic("bad")
	}
	value, err := strconv.Atoi(s[:digitEndScanning])
	if err != nil {
		panic("bad")
	}
	if digitEndScanning < len(s) {
		switch s[digitEndScanning:] {
		case "px":
			unit = Pixel
		case "%":
			unit = Percent
		default:
			panic("bad")
		}
	}
	return CropParam{
		value: value,
		unit:  unit,
	}
}

// func read_matched_imgs_rotate_left_and_save() {
// 	matches, err := filepath.Glob("imgs/*")
// 	if err != nil {
// 		fmt.Println("ERROR  occured when trying to match files")
// 	}
// 	fmt.Println(matches)
// 	for _, filepath := range matches {
// 		cFilepath := C.CString(filepath)
// 		fmt.Println("Filepath:", cFilepath)
// 		res := C.gale_load_img(cFilepath)
// 		C.gale_rotate_img_left(&res)
// 		C.gale_save_img_as(&res, cFilepath, ImgFormatJpg)
// 	}
// 	fmt.Println(ImgFormatPng)
// }
