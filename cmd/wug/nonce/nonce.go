package nonce

import (
	"sort"
)

var Able = [50]string{
	"actignable",
	"anilicable",
	"anvastable",
	"chalinable",
	"comfolvable",
	"compechable",
	"condumable",
	"contaitable",
	"corgervable",
	"covornable",
	"cresucable",
	"enocutable",
	"expeaceable",
	"expelocable",
	"expernable",
	"fispoceable",
	"fupeactable",
	"fusuperable",
	"imalatable",
	"impalvable",
	"inbeadable",
	"inedifiable",
	"infoustable",
	"intoundable",
	"intountable",
	"inveicable",
	"irediocable",
	"mecoushable",
	"parendable",
	"peplaicable",
	"praleckable",
	"preneckable",
	"prequakable",
	"previnable",
	"previtable",
	"puneadable",
	"pustameable",
	"redeptable",
	"rempadable",
	"retaleable",
	"sempoivable",
	"swimitable",
	"tegornable",
	"unaclerable",
	"unalintable",
	"undeperable",
	"unutintable",
	"unvatrable",
	"unvediable",
	"utililable",
}

var Ish = [50]string{
	"badyish",
	"beavish",
	"breyish",
	"carmish",
	"clangish",
	"clurlish",
	"cunkish",
	"devevish",
	"direish",
	"doutish",
	"dwaplish",
	"fadyish",
	"fawkish",
	"fevetish",
	"fevewish",
	"fevilish",
	"frietish",
	"friquish",
	"ghumpish",
	"gireish",
	"goguish",
	"higetish",
	"knarish",
	"laretish",
	"lureish",
	"lurmish",
	"moguish",
	"peftish",
	"preanish",
	"prienish",
	"purerish",
	"radyish",
	"reckish",
	"redyish",
	"rourfish",
	"shigeish",
	"skierish",
	"slarish",
	"slownish",
	"slundish",
	"slungish",
	"snoulish",
	"sonkish",
	"tivilish",
	"turgeish",
	"wabyish",
	"waguish",
	"wainish",
	"wawkish",
	"woungish",
}

var Ive = [50]string{
	"atecusive",
	"cogective",
	"conovative",
	"cormasive",
	"cuminitive",
	"decertive",
	"deflosive",
	"defrertive",
	"dejovative",
	"depulsive",
	"dermasive",
	"dignitive",
	"dimusitive",
	"exhauctive",
	"expecative",
	"extuctive",
	"gederative",
	"imimative",
	"impuctive",
	"indetative",
	"nogensive",
	"nombasive",
	"nonvuptive",
	"nutensive",
	"obsensive",
	"pedititive",
	"pedulsive",
	"pepulative",
	"pransitive",
	"prediasive",
	"prititive",
	"protrative",
	"pumbative",
	"recentive",
	"recumotive",
	"rejeptive",
	"ruchontive",
	"seceptive",
	"sejensive",
	"serposive",
	"submiative",
	"submictive",
	"submistive",
	"sumpertive",
	"sumurative",
	"suprective",
	"tecensive",
	"tendusive",
	"tredictive",
	"vederative",
}

var Ous = [50]string{
	"adodagious",
	"adupendous",
	"anoninous",
	"aurtiguous",
	"cazardous",
	"coivonous",
	"creninous",
	"dardulous",
	"dexarious",
	"erenymous",
	"eretulous",
	"euphitious",
	"eutrigeous",
	"faluminous",
	"fapturous",
	"glamalous",
	"glumonous",
	"gluninous",
	"gropenious",
	"hibeguous",
	"honoderous",
	"indaminous",
	"iniragious",
	"insicious",
	"lasavenous",
	"leamogous",
	"ligegious",
	"liratonous",
	"luticorous",
	"malicinous",
	"meglarious",
	"momogorous",
	"mystuorous",
	"nomeneous",
	"oblicious",
	"pecacious",
	"plalorous",
	"poncorous",
	"prolacious",
	"ralygerous",
	"ravarious",
	"reamorous",
	"rebelorous",
	"slaicitous",
	"suspibious",
	"tefigious",
	"trospurous",
	"undicitous",
	"vexuteous",
	"vombageous",
}

var All []string

func init() {
	All = make([]string, 0, 200)

	All = append(All, Able[:]...)
	All = append(All, Ish[:]...)
	All = append(All, Ive[:]...)
	All = append(All, Ous[:]...)

	sort.Strings(All)
}
