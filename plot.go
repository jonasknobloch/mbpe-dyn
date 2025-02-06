package main

import (
	"gonum.org/v1/gonum/stat"
	pl "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
)

func plot() {
	xs1 := []float64{
		1.1214,
		1.1223,
		1.1718,
		1.1241,
		1.2159,
		1.1639,
		1.2036,
	}

	ys1 := []float64{
		0.4877,
		0.4992,
		0.5199,
		0.5115,
		0.5311,
		0.4936,
		0.5147,
	}

	xs2 := []float64{
		1.1235,
		1.1847,
		1.1972,
		1.4130,
	}

	ys2 := []float64{
		0.4850,
		0.5267,
		0.5601,
		0.5632,
	}

	p := pl.New()

	// p.Title.Text = ""

	p.X.Label.Text = "fertility"
	p.Y.Label.Text = "F1"

	p.X.Min = 1.0
	p.X.Max = 2.0
	p.Y.Min = 0.0
	p.Y.Max = 1.0

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

		p.Add(scatter, lineFunc)

		p.Legend.Add(label, scatter)
	}

	draw(xs1, ys1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, "out")
	draw(xs2, ys2, color.RGBA{R: 0, G: 0, B: 255, A: 255}, "out-mbpe")

	if err := p.Save(8*vg.Inch, 5*vg.Inch, "assets/plot.svg"); err != nil {
		log.Fatal(err)
	}
}
