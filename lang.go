package main

import (
	"fmt"
	_ "path/filepath"
)

// #cgo CFLAGS: -Ivendor
// #cgo LDFLAGS: -lm
// #include "gale.h"
// import "C"

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
	kind LexWordKind
}

type Lexer struct {
	s LexState
	words []LexWord
}

var l Lexer

const program string = `crop 10px 20px 30px 40px
> imgs/картинк*
> scan/imgs/here/*

rotate left
rotate right`


func main() {
	lex()
}

func lex() {
	var w LexWord
	for i, r := range program {
		switch r {
		case '\r', '\n':
			if l.s == LexStateScanWord || l.s == LexStateScanFilepath {
				w.end += 1
				l.words = append(l.words, w)
			}
			l.s = LexStateNil
			if len(l.words) > 0 && l.words[len(l.words)-1].kind != LexNewLWord {
				l.words = append(l.words, LexWord{kind: LexNewLWord})
			}
		case ' ', '\t':
			if l.s == LexStateScanWord {
				w.end += 1
				l.words = append(l.words, w)
			}
			if l.s == LexStateScanFilepath {
				w.end += 1
				continue
			}
			l.s = LexStateNil
		case '>':
			if l.s != LexStateNil {
				panic("error")
			}
			l.s = LexStateScanFilepath
			w = LexWord {beg: i+1, end: i, kind: LexFilepathWord}
		default:
			if l.s == LexStateScanWord || l.s == LexStateScanFilepath {
				w.end += 1
			} else {
				w = LexWord {beg: i, end: i}
				l.s = LexStateScanWord
			}
		}
	}
	if l.s == LexStateScanFilepath || l.s == LexStateScanWord {
		w.end += 1
		l.words = append(l.words, w)
	}
	for _, word := range l.words {
		fmt.Println(program[word.beg:word.end])
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
