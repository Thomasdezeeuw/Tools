package main

import "regexp"

var (
	singleLineEnd     = ".*\n?"
	multiLineContents = "(.|[\r\n])*?"

	singleLineShell = regexp.MustCompile("#" + singleLineEnd)
	singleLineC     = regexp.MustCompile("//" + singleLineEnd)
	multiLineC      = regexp.MustCompile(`/\*` + multiLineContents + `\*/`)
)

// TODO(Thomas): better notation of the languages
var languages = map[string]*language{
	"actionscript": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"as"},
	},
	"asp": {
		[]*regexp.Regexp{regexp.MustCompile("'" + singleLineEnd)},
		nil,
		[]string{"asa", "asp"},
	},
	"c": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"c", "h"},
	},
	"c#": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"cs"},
	},
	"c++": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"c++", "cpp", "cp", "cc", "hh"},
	},
	"clojure": {
		[]*regexp.Regexp{regexp.MustCompile(";" + singleLineEnd)},
		nil,
		[]string{"clj"},
	},
	"css": {
		nil,
		[]*regexp.Regexp{multiLineC},
		[]string{"css"},
	},
	"d": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC, regexp.MustCompile(`/\+` + multiLineContents + `\+/`)},
		[]string{"d", "di"},
	},
	"erlang": {
		[]*regexp.Regexp{regexp.MustCompile("%" + singleLineEnd)},
		nil,
		[]string{"erl", "hrl"},
	},
	"go": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"go"},
	},
	"dot": {
		[]*regexp.Regexp{singleLineC, singleLineShell},
		[]*regexp.Regexp{multiLineC},
		[]string{"dot", "DOT"},
	},
	"groovy": {
		[]*regexp.Regexp{singleLineC, singleLineShell},
		[]*regexp.Regexp{multiLineC},
		[]string{"groovy", "gvy"},
	},
	"haskell": {
		[]*regexp.Regexp{regexp.MustCompile("{-" + multiLineContents + "-}")},
		[]*regexp.Regexp{regexp.MustCompile("--" + singleLineEnd)},
		[]string{"hs"},
	},
	"html": {
		nil,
		[]*regexp.Regexp{regexp.MustCompile("<!--" + multiLineContents + "-->")},
		[]string{"html", "htm", "shtml", "xhtml", "phtml", "tmpl", "tpl"},
	},
	"java": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"java"},
	},
	"javascript": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"js", "jsx"},
	},
	"lisp": {
		[]*regexp.Regexp{regexp.MustCompile("#|" + multiLineContents + "|#")},
		[]*regexp.Regexp{regexp.MustCompile(";" + singleLineEnd)},
		[]string{"lisp", "cl", "l"},
	},
	"lua": {
		nil,
		[]*regexp.Regexp{regexp.MustCompile("--" + singleLineEnd)},
		[]string{"lua"},
	},
	"objective-c": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"m", "mm", "M"},
	},
	"ocaml": {
		[]*regexp.Regexp{regexp.MustCompile(`(\*` + multiLineContents + `\*)`)},
		nil,
		[]string{"ml", "mli", "mll"},
	},
	"pascal": {
		[]*regexp.Regexp{regexp.MustCompile(`(\*` + multiLineContents + `\*)`), regexp.MustCompile("{" + multiLineContents + "}")},
		[]*regexp.Regexp{regexp.MustCompile("--" + singleLineEnd)},
		[]string{"pas", "p"},
	},
	"perl": {
		[]*regexp.Regexp{singleLineShell},
		[]*regexp.Regexp{regexp.MustCompile("^=" + multiLineContents + "^=cut")},
		[]string{"pl", "pm"},
	},
	"php": {
		[]*regexp.Regexp{singleLineC, singleLineShell},
		[]*regexp.Regexp{multiLineC},
		[]string{"php"},
	},
	"python": {
		[]*regexp.Regexp{singleLineShell},
		[]*regexp.Regexp{regexp.MustCompile("\"\"\"" + multiLineContents + "\"\"\""), regexp.MustCompile("'''" + multiLineContents + "'''")},
		[]string{"py", "rpy", "cpy", "pyw"},
	},
	"r": {
		[]*regexp.Regexp{singleLineShell},
		nil,
		[]string{"R", "r", "s", "S"},
	},
	"ruby": {
		[]*regexp.Regexp{singleLineShell},
		[]*regexp.Regexp{regexp.MustCompile("=begin" + multiLineContents + "=end")},
		[]string{"rb", "rbx", "rjs"},
	},
	"rust": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"rs"},
	},
	"scala": {
		[]*regexp.Regexp{singleLineC},
		[]*regexp.Regexp{multiLineC},
		[]string{"scala"},
	},
	"shell": {
		[]*regexp.Regexp{singleLineShell},
		nil,
		[]string{"sh", "bash", "zsh"},
	},
	"sql": {
		[]*regexp.Regexp{regexp.MustCompile("---" + singleLineEnd)},
		[]*regexp.Regexp{multiLineC},
		[]string{"sql"},
	},
	"unkown": {},
}
