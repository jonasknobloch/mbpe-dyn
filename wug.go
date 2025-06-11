package mbpe

type AdjectiveSuffixEvaluator struct {
	able          [50]string
	ish           [50]string
	ive           [50]string
	ous           [50]string
	reportAverage bool
}

func NewAdjectiveSuffixEvaluator() *AdjectiveSuffixEvaluator {
	able := [50]string{
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

	ish := [50]string{
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

	ive := [50]string{
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

	ous := [50]string{
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

	return &AdjectiveSuffixEvaluator{
		able:          able,
		ish:           ish,
		ive:           ive,
		ous:           ous,
		reportAverage: false,
	}
}

func (a *AdjectiveSuffixEvaluator) SetReportAverage(reportAverage bool) {
	a.reportAverage = reportAverage
}

func (a *AdjectiveSuffixEvaluator) Eval(tokenizer Tokenizer, maxRank int) ([]float64, error) {
	m, ok := tokenizer.model.(*MBPE)

	if !ok {
		panic("unexpected model type")
	}

	eval := func(dict [50]string, suffix string) float64 {
		n := 0

		for _, v := range dict {
			t := m.tokenize("Ġ"+v, nil, maxRank)

			r := m.ToString(t)

			if r[len(r)-1] == suffix {
				n++
			}
		}

		return float64(n) / float64(len(dict))
	}

	able := eval(a.able, "able")
	ive := eval(a.ive, "ive")
	ish := eval(a.ish, "ish")
	ous := eval(a.ous, "ous")

	if a.reportAverage {
		return []float64{(able + ive + ish + ous) / 4.0}, nil
	}

	return []float64{able, ive, ish, ous}, nil
}
