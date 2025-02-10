package main

import (
	"gonum.org/v1/gonum/interp"
	"gonum.org/v1/gonum/stat"
	pl "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"sort"
)

type plotData struct {
	xs []float64
	ys []float64
}

func newPlotData(xs, ys []float64) plotData {
	return plotData{
		xs: xs,
		ys: ys,
	}
}

func (p plotData) Len() int {
	return len(p.xs)
}

func (p plotData) Less(i, j int) bool {
	return p.xs[i] < p.xs[j]
}

func (p plotData) Swap(i, j int) {
	p.xs[i], p.xs[j] = p.xs[j], p.xs[i]
	p.ys[i], p.ys[j] = p.ys[j], p.ys[i]
}

func (p plotData) XY(i int) (x, y float64) {
	return p.xs[i], p.ys[i]
}

func plot() {
	// xs1 := []float64{
	// 	1.0425,
	// 	1.0424,
	// 	1.0424,
	// 	1.0425,
	// 	1.0427,
	// 	1.0432,
	// 	1.0439,
	// 	1.0452,
	// 	1.0473,
	// 	1.0528,
	// 	1.1336,
	// }
	//
	// ys1 := []float64{
	// 	0.7674,
	// 	0.7759,
	// 	0.7779,
	// 	0.7809,
	// 	0.7857,
	// 	0.7877,
	// 	0.7917,
	// 	0.7962,
	// 	0.8000,
	// 	0.8049,
	// 	0.8176,
	// }

	xs2 := []float64{
		1.0710,
		1.0708,
		1.0711,
		1.0715,
		1.0721,
		1.0727,
		1.0738,
		1.0754,
		1.0790,
		1.0875,
		1.1685,
	}

	ys2 := []float64{
		0.7844,
		0.7932,
		0.7946,
		0.7976,
		0.8021,
		0.8047,
		0.8092,
		0.8128,
		0.8147,
		0.8166,
		0.8221,
	}

	xs3 := []float64{
		1.1214,
		1.1211,
		1.1215,
		1.1218,
		1.1225,
		1.1241,
		1.1267,
		1.1301,
		1.1366,
		1.1492,
		1.2159,
	}

	ys3 := []float64{
		0.8041,
		0.8120,
		0.8135,
		0.8150,
		0.8189,
		0.8211,
		0.8245,
		0.8277,
		0.8288,
		0.8295,
		0.8307,
	}

	xs4 := []float64{
		1.2021,
		1.2018,
		1.2023,
		1.2035,
		1.2047,
		1.2074,
		1.2109,
		1.2156,
		1.2244,
		1.2390,
		1.2844,
	}

	ys4 := []float64{
		0.8176,
		0.8257,
		0.8265,
		0.8284,
		0.8324,
		0.8335,
		0.8380,
		0.8406,
		0.8427,
		0.8433,
		0.8436,
	}

	// xs5 := []float64{
	// 	1.12141326,
	// 	1.12123169,
	// 	1.12165313,
	// 	1.12170917,
	// 	1.12252964,
	// 	1.12348236,
	// 	1.12576891,
	// 	1.12912026,
	// 	1.13459900,
	// 	1.14703825,
	// 	1.21456753,
	// }
	//
	// ys5 := []float64{
	// 	0.80406044,
	// 	0.80465550,
	// 	0.79859057,
	// 	0.79781999,
	// 	0.79641271,
	// 	0.79576397,
	// 	0.79527510,
	// 	0.79513581,
	// 	0.79367285,
	// 	0.79275649,
	// 	0.78602057,
	// }

	// xs5 := []float64{
	// 	1.3256,
	// 	1.3246,
	// 	1.3251,
	// 	1.3265,
	// 	1.3281,
	// 	1.3314,
	// 	1.3360,
	// 	1.3413,
	// 	1.3509,
	// 	1.3650,
	// 	1.3867,
	// }
	//
	// ys5 := []float64{
	// 	0.8348,
	// 	0.8442,
	// 	0.8448,
	// 	0.8463,
	// 	0.8500,
	// 	0.8517,
	// 	0.8564,
	// 	0.8580,
	// 	0.8594,
	// 	0.8589,
	// 	0.8599,
	// }

	xs6 := []float64{
		1.04250731,
		1.07098599,
		1.12141326,
		1.20205654,
		1.32556581,
	}

	ys7 := []float64{
		0.76738930,
		0.78441537,
		0.80406044,
		0.81763303,
		0.83476191,
	}

	// runner.RunAll(50256, 1<<15, 1<<14, 1<<13)

	p := pl.New()

	// p.Title.Text = ""

	p.X.Label.Text = "Fertility"
	p.Y.Label.Text = "Merge Layer"

	p.X.Min = 1.05
	p.X.Max = 1.32
	p.Y.Min = 0.76
	p.Y.Max = 0.86

	// p.Y.Scale = pl.LogScale{}
	// p.Y.Tick.Marker = pl.LogTicks{}

	drawSpline := func(xs, ys []float64, clr color.RGBA, label string) {
		// m := (ys[1] - ys[0]) / (xs[1] - xs[0]) // slope of first two points
		// x, y := xs[0]-(ys[0]/m), 0.0           // phantom origin
		// xs = append([]float64{x}, xs...)
		// ys = append([]float64{y}, ys...)

		// xs = append([]float64{0}, xs...)
		// ys = append([]float64{0}, ys...)

		pts := newPlotData(xs, ys)

		sort.Sort(pts)

		predictor := interp.FittablePredictor(&interp.FritschButland{})

		if err := predictor.Fit(xs, ys); err != nil {
			log.Fatal(err)
		}

		line := plotter.NewFunction(func(x float64) float64 {
			return predictor.Predict(x)
		})

		line.Color = clr
		line.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}

		p.Add(line)
	}

	draw := func(xs, ys []float64, color color.RGBA, label string) {
		intercept, slope := stat.LinearRegression(xs, ys, nil, false)

		scatterData := make(plotter.XYs, len(xs))

		for i := range xs {
			scatterData[i].X = xs[i]
			scatterData[i].Y = ys[i]
		}

		scatter, err := plotter.NewScatter(scatterData)

		scatter.Color = color

		if err != nil {
			log.Fatal(err)
		}

		scatter.Radius = 2

		lineFunc := plotter.NewFunction(func(x float64) float64 {
			return intercept + slope*x
		})

		lineFunc.Color = color

		pts := newPlotData(xs, ys)

		line, err := plotter.NewLine(pts)

		if err != nil {
			log.Fatal(err)
		}

		line.Color = color

		p.Add(scatter)
		// p.Add(scatter, lineFunc)

		if label != "baseline" {
			p.Add(line)
		}

		p.Legend.Add(label, scatter)
	}

	// draw(xs1, ys1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, "2^17")
	draw(xs2, ys2, color.RGBA{R: 0, G: 0, B: 255, A: 255}, "2^16")
	draw(xs3, ys3, color.RGBA{R: 0, G: 255, B: 0, A: 255}, "2^15")
	draw(xs4, ys4, color.RGBA{R: 255, G: 0, B: 255, A: 255}, "2^14")
	// draw(xs5, ys5, color.RGBA{R: 0, G: 255, B: 255, A: 255}, "2^13")
	draw(xs6, ys7, color.RGBA{R: 0, G: 0, B: 0, A: 255}, "baseline")

	// TODO use dots circles ans squares
	// TODO use linear step for baseline but include starting points

	// drawSpline(xs1, ys1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, "2^17")
	// drawSpline(xs2, ys2, color.RGBA{R: 0, G: 0, B: 255, A: 255}, "2^16")
	// drawSpline(xs3, ys3, color.RGBA{R: 0, G: 255, B: 0, A: 255}, "2^15")
	// drawSpline(xs4, ys4, color.RGBA{R: 255, G: 0, B: 255, A: 255}, "2^14")
	// drawSpline(xs5, ys5, color.RGBA{R: 0, G: 255, B: 255, A: 255}, "2^15 dummy")
	drawSpline(xs6, ys7, color.RGBA{R: 0, G: 0, B: 0, A: 255}, "baseline")

	if err := p.Save(8*vg.Inch, 5*vg.Inch, "assets/plot.svg"); err != nil {
		log.Fatal(err)
	}
}
