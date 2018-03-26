package inflect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// baseAcronyms comes from https://en.wikipedia.org/wiki/List_of_information_technology_acronymss
const baseAcronyms = `ACK,ACL,ADSL,AES,ANSI,API,ARP,ATM,BGP,BSS,CAT,CCITT,CHAP,CIDR,CIR,CLI,CPE,CPU,CRC,CRT,CSMA,CMOS,DCE,DEC,DES,DHCP,DNS,DRAM,DSL,DSLAM,DTE,DMI,EHA,EIA,EIGRP,EOF,ESS,FCC,FCS,FDDI,FTP,GBIC,gbps,GEPOF,HDLC,HTTP,HTTPS,IANA,ICMP,IDF,IDS,IEEE,IETF,IMAP,IP,IPS,ISDN,ISP,kbps,LACP,LAN,LAPB,LAPF,LLC,MAC,MAN,Mbps,MC,MDF,MIB,MoCA,MPLS,MTU,NAC,NAT,NBMA,NIC,NRZ,NRZI,NVRAM,OSI,OSPF,OUI,PAP,PAT,PC,PIM,PIM,PCM,PDU,POP3,POP,POST,POTS,PPP,PPTP,PTT,PVST,RADIUS,RAM,RARP,RFC,RIP,RLL,ROM,RSTP,RTP,RCP,SDLC,SFD,SFP,SLARP,SLIP,SMTP,SNA,SNAP,SNMP,SOF,SRAM,SSH,SSID,STP,SYN,TDM,TFTP,TIA,TOFU,UDP,URL,URI,USB,UTP,VC,VLAN,VLSM,VPN,W3C,WAN,WEP,WiFi,WPA,WWW`

// Rule used by rulesets
type Rule struct {
	suffix      string
	replacement string
	exact       bool
}

// Ruleset a Ruleset is the config of pluralization rules
// you can extend the rules with the Add* methods
type Ruleset struct {
	uncountables map[string]bool
	plurals      []*Rule
	singulars    []*Rule
	humans       []*Rule
	acronyms     []*Rule
}

// create a blank ruleset. Unless you are going to
// build your own rules from scratch you probably
// won't need this and can just use the defaultRuleset
// via the global inflect.* methods
func NewRuleset() *Ruleset {
	rs := new(Ruleset)
	rs.uncountables = make(map[string]bool)
	rs.plurals = make([]*Rule, 0)
	rs.singulars = make([]*Rule, 0)
	rs.humans = make([]*Rule, 0)
	rs.acronyms = make([]*Rule, 0)
	return rs
}

// create a new ruleset and load it with the default
// set of common English pluralization rules
func NewDefaultRuleset() *Ruleset {
	rs := NewRuleset()
	rs.AddPlural("s", "s")
	rs.AddPlural("testis", "testes")
	rs.AddPlural("axis", "axes")
	rs.AddPlural("octopus", "octopi")
	rs.AddPlural("virus", "viri")
	rs.AddPlural("octopi", "octopi")
	rs.AddPlural("viri", "viri")
	rs.AddPlural("alias", "aliases")
	rs.AddPlural("status", "statuses")
	rs.AddPlural("Status", "Statuses")
	rs.AddPlural("bus", "buses")
	rs.AddPlural("buffalo", "buffaloes")
	rs.AddPlural("tomato", "tomatoes")
	rs.AddPlural("tum", "ta")
	rs.AddPlural("ium", "ia")
	rs.AddPlural("ta", "ta")
	rs.AddPlural("ia", "ia")
	rs.AddPlural("sis", "ses")
	rs.AddPlural("lf", "lves")
	rs.AddPlural("rf", "rves")
	rs.AddPlural("afe", "aves")
	rs.AddPlural("bfe", "bves")
	rs.AddPlural("cfe", "cves")
	rs.AddPlural("dfe", "dves")
	rs.AddPlural("efe", "eves")
	rs.AddPlural("gfe", "gves")
	rs.AddPlural("hfe", "hves")
	rs.AddPlural("ife", "ives")
	rs.AddPlural("jfe", "jves")
	rs.AddPlural("kfe", "kves")
	rs.AddPlural("lfe", "lves")
	rs.AddPlural("mfe", "mves")
	rs.AddPlural("nfe", "nves")
	rs.AddPlural("ofe", "oves")
	rs.AddPlural("pfe", "pves")
	rs.AddPlural("qfe", "qves")
	rs.AddPlural("rfe", "rves")
	rs.AddPlural("sfe", "sves")
	rs.AddPlural("tfe", "tves")
	rs.AddPlural("ufe", "uves")
	rs.AddPlural("vfe", "vves")
	rs.AddPlural("wfe", "wves")
	rs.AddPlural("xfe", "xves")
	rs.AddPlural("yfe", "yves")
	rs.AddPlural("zfe", "zves")
	rs.AddPlural("hive", "hives")
	rs.AddPlural("quy", "quies")
	rs.AddPlural("by", "bies")
	rs.AddPlural("cy", "cies")
	rs.AddPlural("dy", "dies")
	rs.AddPlural("fy", "fies")
	rs.AddPlural("gy", "gies")
	rs.AddPlural("hy", "hies")
	rs.AddPlural("jy", "jies")
	rs.AddPlural("ky", "kies")
	rs.AddPlural("ly", "lies")
	rs.AddPlural("my", "mies")
	rs.AddPlural("ny", "nies")
	rs.AddPlural("py", "pies")
	rs.AddPlural("qy", "qies")
	rs.AddPlural("ry", "ries")
	rs.AddPlural("sy", "sies")
	rs.AddPlural("ty", "ties")
	rs.AddPlural("vy", "vies")
	rs.AddPlural("wy", "wies")
	rs.AddPlural("xy", "xies")
	rs.AddPlural("zy", "zies")
	rs.AddPlural("x", "xes")
	rs.AddPlural("ch", "ches")
	rs.AddPlural("ss", "sses")
	rs.AddPlural("sh", "shes")
	rs.AddPlural("matrix", "matrices")
	rs.AddPlural("vertix", "vertices")
	rs.AddPlural("indix", "indices")
	rs.AddPlural("matrex", "matrices")
	rs.AddPlural("vertex", "vertices")
	rs.AddPlural("index", "indices")
	rs.AddPlural("mouse", "mice")
	rs.AddPlural("louse", "lice")
	rs.AddPlural("mice", "mice")
	rs.AddPlural("lice", "lice")
	rs.AddPlural("ress", "resses")
	rs.AddPluralExact("ox", "oxen", true)
	rs.AddPluralExact("oxen", "oxen", true)
	rs.AddPluralExact("quiz", "quizzes", true)
	rs.AddSingular("s", "")
	rs.AddSingular("ss", "ss")
	rs.AddSingular("news", "news")
	rs.AddSingular("ta", "tum")
	rs.AddSingular("ia", "ium")
	rs.AddSingular("analyses", "analysis")
	rs.AddSingular("bases", "basis")
	rs.AddSingularExact("basis", "basis", true)
	rs.AddSingular("diagnoses", "diagnosis")
	rs.AddSingularExact("diagnosis", "diagnosis", true)
	rs.AddSingular("parentheses", "parenthesis")
	rs.AddSingular("prognoses", "prognosis")
	rs.AddSingular("synopses", "synopsis")
	rs.AddSingular("theses", "thesis")
	rs.AddSingular("analyses", "analysis")
	rs.AddSingularExact("analysis", "analysis", true)
	rs.AddSingular("aves", "afe")
	rs.AddSingular("bves", "bfe")
	rs.AddSingular("cves", "cfe")
	rs.AddSingular("dves", "dfe")
	rs.AddSingular("eves", "efe")
	rs.AddSingular("gves", "gfe")
	rs.AddSingular("hves", "hfe")
	rs.AddSingular("ives", "ife")
	rs.AddSingular("jves", "jfe")
	rs.AddSingular("kves", "kfe")
	rs.AddSingular("lves", "lfe")
	rs.AddSingular("mves", "mfe")
	rs.AddSingular("nves", "nfe")
	rs.AddSingular("oves", "ofe")
	rs.AddSingular("pves", "pfe")
	rs.AddSingular("qves", "qfe")
	rs.AddSingular("rves", "rfe")
	rs.AddSingular("sves", "sfe")
	rs.AddSingular("tves", "tfe")
	rs.AddSingular("uves", "ufe")
	rs.AddSingular("vves", "vfe")
	rs.AddSingular("wves", "wfe")
	rs.AddSingular("xves", "xfe")
	rs.AddSingular("yves", "yfe")
	rs.AddSingular("zves", "zfe")
	rs.AddSingular("hives", "hive")
	rs.AddSingular("tives", "tive")
	rs.AddSingular("lves", "lf")
	rs.AddSingular("rves", "rf")
	rs.AddSingular("quies", "quy")
	rs.AddSingular("bies", "by")
	rs.AddSingular("cies", "cy")
	rs.AddSingular("dies", "dy")
	rs.AddSingular("fies", "fy")
	rs.AddSingular("gies", "gy")
	rs.AddSingular("hies", "hy")
	rs.AddSingular("jies", "jy")
	rs.AddSingular("kies", "ky")
	rs.AddSingular("lies", "ly")
	rs.AddSingular("mies", "my")
	rs.AddSingular("nies", "ny")
	rs.AddSingular("pies", "py")
	rs.AddSingular("qies", "qy")
	rs.AddSingular("ries", "ry")
	rs.AddSingular("sies", "sy")
	rs.AddSingular("ties", "ty")
	rs.AddSingular("vies", "vy")
	rs.AddSingular("wies", "wy")
	rs.AddSingular("xies", "xy")
	rs.AddSingular("zies", "zy")
	rs.AddSingular("series", "series")
	rs.AddSingular("movies", "movie")
	rs.AddSingular("xes", "x")
	rs.AddSingular("ches", "ch")
	rs.AddSingular("sses", "ss")
	rs.AddSingular("shes", "sh")
	rs.AddSingular("mice", "mouse")
	rs.AddSingular("lice", "louse")
	rs.AddSingular("buses", "bus")
	rs.AddSingularExact("bus", "bus", true)
	rs.AddSingular("oes", "o")
	rs.AddSingular("shoes", "shoe")
	rs.AddSingular("crises", "crisis")
	rs.AddSingularExact("crisis", "crisis", true)
	rs.AddSingular("axes", "axis")
	rs.AddSingularExact("axis", "axis", true)
	rs.AddSingular("testes", "testis")
	rs.AddSingularExact("testis", "testis", true)
	rs.AddSingular("octopi", "octopus")
	rs.AddSingularExact("octopus", "octopus", true)
	rs.AddSingular("viri", "virus")
	rs.AddSingularExact("virus", "virus", true)
	rs.AddSingular("statuses", "status")
	rs.AddSingular("Statuses", "Status")
	rs.AddSingularExact("status", "status", true)
	rs.AddSingularExact("Status", "Status", true)
	rs.AddSingular("aliases", "alias")
	rs.AddSingularExact("alias", "alias", true)
	rs.AddSingularExact("oxen", "ox", true)
	rs.AddSingular("vertices", "vertex")
	rs.AddSingular("indices", "index")
	rs.AddSingular("matrices", "matrix")
	rs.AddSingularExact("quizzes", "quiz", true)
	rs.AddSingular("databases", "database")
	rs.AddSingular("resses", "ress")
	rs.AddSingular("ress", "ress")
	rs.AddIrregular("person", "people")
	rs.AddIrregular("man", "men")
	rs.AddIrregular("child", "children")
	rs.AddIrregular("sex", "sexes")
	rs.AddIrregular("move", "moves")
	rs.AddIrregular("zombie", "zombies")
	rs.AddIrregular("Status", "Statuses")
	rs.AddIrregular("status", "statuses")
	rs.AddUncountable("equipment")
	rs.AddUncountable("information")
	rs.AddUncountable("rice")
	rs.AddUncountable("money")
	rs.AddUncountable("species")
	rs.AddUncountable("series")
	rs.AddUncountable("fish")
	rs.AddUncountable("sheep")
	rs.AddUncountable("jeans")
	rs.AddUncountable("police")

	acronyms := strings.Split(baseAcronyms, ",")
	for _, acr := range acronyms {
		rs.AddAcronym(acr)
	}

	return rs
}

func (rs *Ruleset) Uncountables() map[string]bool {
	return rs.uncountables
}

// add a pluralization rule
func (rs *Ruleset) AddPlural(suffix, replacement string) {
	rs.AddPluralExact(suffix, replacement, false)
}

// add a pluralization rule with full string match
func (rs *Ruleset) AddPluralExact(suffix, replacement string, exact bool) {
	// remove uncountable
	delete(rs.uncountables, suffix)
	// create rule
	r := new(Rule)
	r.suffix = suffix
	r.replacement = replacement
	r.exact = exact
	// prepend
	rs.plurals = append([]*Rule{r}, rs.plurals...)
}

// add a singular rule
func (rs *Ruleset) AddSingular(suffix, replacement string) {
	rs.AddSingularExact(suffix, replacement, false)
}

// same as AddSingular but you can set `exact` to force
// a full string match
func (rs *Ruleset) AddSingularExact(suffix, replacement string, exact bool) {
	// remove from uncountable
	delete(rs.uncountables, suffix)
	// create rule
	r := new(Rule)
	r.suffix = suffix
	r.replacement = replacement
	r.exact = exact
	rs.singulars = append([]*Rule{r}, rs.singulars...)
}

// Human rules are applied by humanize to show more friendly
// versions of words
func (rs *Ruleset) AddHuman(suffix, replacement string) {
	r := new(Rule)
	r.suffix = suffix
	r.replacement = replacement
	rs.humans = append([]*Rule{r}, rs.humans...)
}

// Add any inconsistent pluralizing/singularizing rules
// to the set here.
func (rs *Ruleset) AddIrregular(singular, plural string) {
	delete(rs.uncountables, singular)
	delete(rs.uncountables, plural)
	rs.AddPlural(singular, plural)
	rs.AddPlural(plural, plural)
	rs.AddSingular(plural, singular)
}

// if you use acronym you may need to add them to the ruleset
// to prevent Underscored words of things like "HTML" coming out
// as "h_t_m_l"
func (rs *Ruleset) AddAcronym(word string) {
	r := new(Rule)
	r.suffix = word
	r.replacement = rs.Titleize(strings.ToLower(word))
	rs.acronyms = append(rs.acronyms, r)
}

// add a word to this ruleset that has the same singular and plural form
// for example: "rice"
func (rs *Ruleset) AddUncountable(word string) {
	rs.uncountables[strings.ToLower(word)] = true
}

func (rs *Ruleset) isUncountable(word string) bool {
	// handle multiple words by using the last one
	words := strings.Split(word, " ")
	if _, exists := rs.uncountables[strings.ToLower(words[len(words)-1])]; exists {
		return true
	}
	return false
}

//isAcronym returns if a word is acronym or not.
func (rs *Ruleset) isAcronym(word string) bool {
	for _, rule := range rs.acronyms {
		if rule.suffix == word {
			return true
		}
	}

	return false
}

func (rs *Ruleset) PluralizeWithSize(word string, size int) string {
	if size == 1 {
		return rs.Singularize(word)
	}
	return rs.Pluralize(word)
}

// returns the plural form of a singular word
func (rs *Ruleset) Pluralize(word string) string {
	if len(word) == 0 {
		return word
	}
	lWord := strings.ToLower(word)
	if rs.isUncountable(lWord) {
		return word
	}

	var candidate string
	for _, rule := range rs.plurals {
		if rule.exact {
			if lWord == rule.suffix {
				// Capitalized word
				if lWord[0] != word[0] && lWord[1:] == word[1:] {
					return Capitalize(rule.replacement)
				}
				return rule.replacement
			}
			continue
		}

		if strings.EqualFold(word, rule.suffix) {
			candidate = rule.replacement
		}

		if strings.HasSuffix(word, rule.suffix) {
			return replaceLast(word, rule.suffix, rule.replacement)
		}
	}

	if candidate != "" {
		return candidate
	}
	return word + "s"
}

// returns the singular form of a plural word
func (rs *Ruleset) Singularize(word string) string {
	if len(word) <= 1 {
		return word
	}
	lWord := strings.ToLower(word)
	if rs.isUncountable(lWord) {
		return word
	}

	var candidate string

	for _, rule := range rs.singulars {
		if rule.exact {
			if lWord == rule.suffix {
				// Capitalized word
				if lWord[0] != word[0] && lWord[1:] == word[1:] {
					return Capitalize(rule.replacement)
				}
				return rule.replacement
			}
			continue
		}

		if strings.EqualFold(word, rule.suffix) {
			candidate = rule.replacement
		}

		if strings.HasSuffix(word, rule.suffix) {
			return replaceLast(word, rule.suffix, rule.replacement)
		}
	}

	if candidate != "" {
		return candidate
	}

	return word
}

// uppercase first character
func (rs *Ruleset) Capitalize(word string) string {
	if strings.ToLower(word) == "id" {
		return "ID"
	}
	return strings.ToUpper(word[:1]) + word[1:]
}

// "dino_party" -> "DinoParty"
func (rs *Ruleset) Camelize(word string) string {
	if strings.ToLower(word) == "id" {
		return "ID"
	}
	words := splitAtCaseChangeWithTitlecase(word)
	return strings.Join(words, "")
}

// same as Camelcase but with first letter downcased
func (rs *Ruleset) CamelizeDownFirst(word string) string {
	word = Camelize(word)
	return strings.ToLower(word[:1]) + word[1:]
}

// Capitalize every word in sentence "hello there" -> "Hello There"
func (rs *Ruleset) Titleize(word string) string {
	words := splitAtCaseChangeWithTitlecase(word)
	result := strings.Join(words, " ")

	var acronymWords []string
	for index, word := range words {
		if len(word) == 1 {
			acronymWords = append(acronymWords, word)
		}

		if len(word) > 1 || index == len(words)-1 || len(acronymWords) > 1 {
			acronym := strings.Join(acronymWords, "")
			if !rs.isAcronym(acronym) {
				acronymWords = acronymWords[:len(acronymWords)]
				continue
			}

			result = strings.Replace(result, strings.Join(acronymWords, " "), acronym, 1)
			acronymWords = []string{}
		}
	}

	return result
}

func (rs *Ruleset) safeCaseAcronyms(word string) string {
	// convert an acronym like HTML into Html
	for _, rule := range rs.acronyms {
		word = strings.Replace(word, rule.suffix, rule.replacement, -1)
	}
	return word
}

func (rs *Ruleset) separatedWords(word, sep string) string {
	word = rs.safeCaseAcronyms(word)
	words := splitAtCaseChange(word)
	return strings.Join(words, sep)
}

// lowercase underscore version "BigBen" -> "big_ben"
func (rs *Ruleset) Underscore(word string) string {
	return rs.separatedWords(word, "_")
}

// First letter of sentence capitalized
// Uses custom friendly replacements via AddHuman()
func (rs *Ruleset) Humanize(word string) string {
	word = replaceLast(word, "_id", "") // strip foreign key kinds
	// replace and strings in humans list
	for _, rule := range rs.humans {
		word = strings.Replace(word, rule.suffix, rule.replacement, -1)
	}
	sentence := rs.separatedWords(word, " ")

	r, n := utf8.DecodeRuneInString(sentence)
	return string(unicode.ToUpper(r)) + sentence[n:]
}

// an underscored foreign key name "Person" -> "person_id"
func (rs *Ruleset) ForeignKey(word string) string {
	return rs.Underscore(rs.Singularize(word)) + "_id"
}

// a foreign key (with an underscore) "Person" -> "personid"
func (rs *Ruleset) ForeignKeyCondensed(word string) string {
	return rs.Underscore(word) + "id"
}

// Rails style pluralized table names: "SuperPerson" -> "super_people"
func (rs *Ruleset) Tableize(word string) string {
	return rs.Pluralize(rs.Underscore(rs.Typeify(word)))
}

var notUrlSafe *regexp.Regexp = regexp.MustCompile(`[^\w\d\-_ ]`)

// param safe dasherized names like "my-param"
func (rs *Ruleset) Parameterize(word string) string {
	return ParameterizeJoin(word, "-")
}

// param safe dasherized names with custom separator
func (rs *Ruleset) ParameterizeJoin(word, sep string) string {
	word = strings.ToLower(word)
	word = rs.Asciify(word)
	word = notUrlSafe.ReplaceAllString(word, "")
	word = strings.Replace(word, " ", sep, -1)
	if len(sep) > 0 {
		squash, err := regexp.Compile(sep + "+")
		if err == nil {
			word = squash.ReplaceAllString(word, sep)
		}
	}
	word = strings.Trim(word, sep+" ")
	return word
}

var lookalikes map[string]*regexp.Regexp = map[string]*regexp.Regexp{
	"A":  regexp.MustCompile(`À|Á|Â|Ã|Ä|Å`),
	"AE": regexp.MustCompile(`Æ`),
	"C":  regexp.MustCompile(`Ç`),
	"E":  regexp.MustCompile(`È|É|Ê|Ë`),
	"G":  regexp.MustCompile(`Ğ`),
	"I":  regexp.MustCompile(`Ì|Í|Î|Ï|İ`),
	"N":  regexp.MustCompile(`Ñ`),
	"O":  regexp.MustCompile(`Ò|Ó|Ô|Õ|Ö|Ø`),
	"S":  regexp.MustCompile(`Ş`),
	"U":  regexp.MustCompile(`Ù|Ú|Û|Ü`),
	"Y":  regexp.MustCompile(`Ý`),
	"ss": regexp.MustCompile(`ß`),
	"a":  regexp.MustCompile(`à|á|â|ã|ä|å`),
	"ae": regexp.MustCompile(`æ`),
	"c":  regexp.MustCompile(`ç`),
	"e":  regexp.MustCompile(`è|é|ê|ë`),
	"g":  regexp.MustCompile(`ğ`),
	"i":  regexp.MustCompile(`ì|í|î|ï|ı`),
	"n":  regexp.MustCompile(`ñ`),
	"o":  regexp.MustCompile(`ò|ó|ô|õ|ö|ø`),
	"s":  regexp.MustCompile(`ş`),
	"u":  regexp.MustCompile(`ù|ú|û|ü|ũ|ū|ŭ|ů|ű|ų`),
	"y":  regexp.MustCompile(`ý|ÿ`),
}

// transforms Latin characters like é -> e
func (rs *Ruleset) Asciify(word string) string {
	for repl, regex := range lookalikes {
		word = regex.ReplaceAllString(word, repl)
	}
	return word
}

var tablePrefix *regexp.Regexp = regexp.MustCompile(`^[^.]*\.`)

// "something_like_this" -> "SomethingLikeThis"
func (rs *Ruleset) Typeify(word string) string {
	word = tablePrefix.ReplaceAllString(word, "")
	return rs.Camelize(rs.Singularize(word))
}

// "SomeText" -> "some-text"
func (rs *Ruleset) Dasherize(word string) string {
	return rs.separatedWords(word, "-")
}

// "1031" -> "1031st"
func (rs *Ruleset) Ordinalize(str string) string {
	number, err := strconv.Atoi(str)
	if err != nil {
		return str
	}
	switch abs(number) % 100 {
	case 11, 12, 13:
		return fmt.Sprintf("%dth", number)
	default:
		switch abs(number) % 10 {
		case 1:
			return fmt.Sprintf("%dst", number)
		case 2:
			return fmt.Sprintf("%dnd", number)
		case 3:
			return fmt.Sprintf("%drd", number)
		}
	}
	return fmt.Sprintf("%dth", number)
}

func (rs *Ruleset) ForeignKeyToAttribute(str string) string {
	w := rs.Camelize(str)
	if strings.HasSuffix(w, "Id") {
		return strings.TrimSuffix(w, "Id") + "ID"
	}
	return w
}

func (rs *Ruleset) LoadReader(r io.Reader) error {
	m := map[string]string{}
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		return fmt.Errorf("could not decode inflection JSON from reader: %s", err)
	}
	for s, p := range m {
		defaultRuleset.AddIrregular(s, p)
	}
	return nil
}

/////////////////////////////////////////
// the default global ruleset
//////////////////////////////////////////

var defaultRuleset *Ruleset

func LoadReader(r io.Reader) error {
	return defaultRuleset.LoadReader(r)
}

func init() {
	defaultRuleset = NewDefaultRuleset()

	pwd, _ := os.Getwd()
	cfg := filepath.Join(pwd, "inflections.json")
	if p := os.Getenv("INFLECT_PATH"); p != "" {
		cfg = p
	}
	if _, err := os.Stat(cfg); err == nil {
		b, err := ioutil.ReadFile(cfg)
		if err != nil {
			fmt.Printf("could not read inflection file %s (%s)\n", cfg, err)
			return
		}
		if err = defaultRuleset.LoadReader(bytes.NewReader(b)); err != nil {
			fmt.Println(err)
		}
	}
}

func Uncountables() map[string]bool {
	return defaultRuleset.Uncountables()
}

func AddPlural(suffix, replacement string) {
	defaultRuleset.AddPlural(suffix, replacement)
}

func AddSingular(suffix, replacement string) {
	defaultRuleset.AddSingular(suffix, replacement)
}

func AddHuman(suffix, replacement string) {
	defaultRuleset.AddHuman(suffix, replacement)
}

func AddIrregular(singular, plural string) {
	defaultRuleset.AddIrregular(singular, plural)
}

func AddAcronym(word string) {
	defaultRuleset.AddAcronym(word)
}

func AddUncountable(word string) {
	defaultRuleset.AddUncountable(word)
}

func Pluralize(word string) string {
	return defaultRuleset.Pluralize(word)
}

func PluralizeWithSize(word string, size int) string {
	return defaultRuleset.PluralizeWithSize(word, size)
}

func Singularize(word string) string {
	return defaultRuleset.Singularize(word)
}

func Capitalize(word string) string {
	return defaultRuleset.Capitalize(word)
}

func Camelize(word string) string {
	return defaultRuleset.Camelize(word)
}

func CamelizeDownFirst(word string) string {
	return defaultRuleset.CamelizeDownFirst(word)
}

func Titleize(word string) string {
	return defaultRuleset.Titleize(word)
}

func Underscore(word string) string {
	return defaultRuleset.Underscore(word)
}

func Humanize(word string) string {
	return defaultRuleset.Humanize(word)
}

func ForeignKey(word string) string {
	return defaultRuleset.ForeignKey(word)
}

func ForeignKeyCondensed(word string) string {
	return defaultRuleset.ForeignKeyCondensed(word)
}

func Tableize(word string) string {
	return defaultRuleset.Tableize(word)
}

func Parameterize(word string) string {
	return defaultRuleset.Parameterize(word)
}

func ParameterizeJoin(word, sep string) string {
	return defaultRuleset.ParameterizeJoin(word, sep)
}

func Typeify(word string) string {
	return defaultRuleset.Typeify(word)
}

func Dasherize(word string) string {
	return defaultRuleset.Dasherize(word)
}

func Ordinalize(word string) string {
	return defaultRuleset.Ordinalize(word)
}

func Asciify(word string) string {
	return defaultRuleset.Asciify(word)
}

func ForeignKeyToAttribute(word string) string {
	return defaultRuleset.ForeignKeyToAttribute(word)
}

// helper funcs

func reverse(s string) string {
	o := make([]rune, utf8.RuneCountInString(s))
	i := len(o)
	for _, c := range s {
		i--
		o[i] = c
	}
	return string(o)
}

func isSpacerChar(c rune) bool {
	switch {
	case c == rune("_"[0]):
		return true
	case c == rune(" "[0]):
		return true
	case c == rune(":"[0]):
		return true
	case c == rune("-"[0]):
		return true
	}
	return false
}

func splitAtCaseChange(s string) []string {
	words := make([]string, 0)
	word := make([]rune, 0)
	for _, c := range s {
		spacer := isSpacerChar(c)
		if len(word) > 0 {
			if unicode.IsUpper(c) || spacer {
				words = append(words, string(word))
				word = make([]rune, 0)
			}
		}
		if !spacer {
			word = append(word, unicode.ToLower(c))
		}
	}
	words = append(words, string(word))
	return words
}

func splitAtCaseChangeWithTitlecase(s string) []string {
	words := make([]string, 0)
	word := make([]rune, 0)

	for _, c := range s {
		spacer := isSpacerChar(c)
		if len(word) > 0 {
			if unicode.IsUpper(c) || spacer {
				words = append(words, string(word))
				word = make([]rune, 0)
			}
		}
		if !spacer {
			if len(word) > 0 {
				word = append(word, unicode.ToLower(c))
			} else {
				word = append(word, unicode.ToUpper(c))
			}
		}
	}

	words = append(words, string(word))
	return words
}

func replaceLast(s, match, repl string) string {
	// reverse strings
	srev := reverse(s)
	mrev := reverse(match)
	rrev := reverse(repl)
	// match first and reverse back
	return reverse(strings.Replace(srev, mrev, rrev, 1))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
