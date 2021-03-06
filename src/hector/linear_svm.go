package hector

import (
	"math"
	"strconv"
	"math/rand"
	"fmt"
)

//This algorithm implement L1 Linear SVM described in "A Dual Coordinate Descent Method for Large-scale Linear SVM"
//You can download the paper from http://ntu.csie.org/~cjlin/papers/cddual.pdf
type LinearSVM struct {
	sv []*Vector
	y []float64
	a []float64
	b float64
	C float64
	e float64
	w *Vector

	xx []float64
}

func (c *LinearSVM) Init(params map[string]string){
	c.C,_ = strconv.ParseFloat(params["c"], 64)
	c.e,_ = strconv.ParseFloat(params["e"], 64)

	c.w = NewVector()
}

func (c *LinearSVM) Predict(sample *Sample) float64 {
	x := sample.GetFeatureVector()
	return c.PredictVector(x)
}

func (c *LinearSVM) PredictVector(x *Vector) float64 {
	ret := c.w.Dot(x)
	return ret
}

func (c *LinearSVM) Train(dataset *DataSet) {
	c.sv = []*Vector{}
	c.y = []float64{}
	c.a = []float64{}
	for k, sample := range dataset.Samples {
		x := sample.GetFeatureVector()
		c.sv = append(c.sv, x)
		c.xx = append(c.xx, x.Dot(x))
		if sample.Label > 0.0 {
			c.y = append(c.y, 1.0)
		} else {
			c.y = append(c.y, -1.0)
		}
		c.a = append(c.a, c.C * rand.Float64() * 0.0)
		c.w.AddVector(x, c.y[k] * c.a[k])
	}

	da0 := 0.0
	for {
		da := 0.0
		for i, ai := range c.a {
			g := c.y[i] * c.w.Dot(c.sv[i]) - 1.0
			pg := g
			if ai < 1e-9 {
				pg = math.Min(0.0, g)
			} else if ai > c.C - 1e-9 {
				pg = math.Max(0.0, g)
			}

			if math.Abs(pg) > 1e-9 {
				ai0 := ai
				ai = math.Min(math.Max(0, ai - g / c.xx[i]), c.C)
				c.w.AddVector(c.sv[i], (ai - ai0) * c.y[i])
				da += math.Abs(ai - ai0)
			}
		}
		da /= float64(len(c.a))
		fmt.Println(da)
		if da < c.e || math.Abs(da - da0) < 1e-3 {
			break
		}
		da0 = da
	}
}