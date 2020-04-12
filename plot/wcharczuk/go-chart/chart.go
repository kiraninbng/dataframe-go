// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package chart

import (
	"context"
	"time"

	"github.com/wcharczuk/go-chart"

	dataframe "github.com/rocketlaunchr/dataframe-go"
)

// S converts a SeriesFloat64 to a chart.Series for usage with the "github.com/wcharczuk/go-chart" package.
// Currently x can be nil, a SeriesFloat64 or a SeriesTime. nil values in the x and y Series are ignored.
func S(ctx context.Context, y *dataframe.SeriesFloat64, x interface{}, r ...dataframe.Range) (chart.Series, error) {

	var out chart.Series

	if len(r) == 0 {
		r = append(r, dataframe.Range{})
	}

	yNRows := y.NRows(dataframe.DontLock)

	start, end, err := r[0].Limits(yNRows)
	if err != nil {
		return nil, err
	}

	switch xx := x.(type) {
	case nil:
		cs := chart.ContinuousSeries{Name: y.Name(dataframe.DontLock)}

		xVals := []float64{}
		yVals := []float64{}

		// Remove nil values
		for i, j := 0, start; j < end+1; i, j = i+1, j+1 {

			if err := ctx.Err(); err != nil {
				return nil, err
			}

			yval := y.Values[j]

			if dataframe.IsValidFloat64(yval) {
				yVals = append(yVals, yval)
				xVals = append(xVals, float64(i))
			}
		}

		cs.XValues = xVals
		cs.YValues = yVals

		out = cs
	case *dataframe.SeriesFloat64:

		cs := chart.ContinuousSeries{Name: y.Name(dataframe.DontLock)}

		xVals := []float64{}
		yVals := []float64{}

		// Remove nil values
		for i, j := 0, start; j < end+1; i, j = i+1, j+1 {

			if err := ctx.Err(); err != nil {
				return nil, err
			}

			yval := y.Values[j]
			xval := xx.Values[j]

			if dataframe.IsValidFloat64(yval) {
				// Check x val is valid
				if dataframe.IsValidFloat64(xval) {
					yVals = append(yVals, yval)
					xVals = append(xVals, xval)
				}
			}
		}

		cs.XValues = xVals
		cs.YValues = yVals

		out = cs
	case *dataframe.SeriesTime:

		cs := chart.TimeSeries{Name: y.Name(dataframe.DontLock)}

		xVals := []time.Time{}
		yVals := []float64{}

		// Remove nil values
		for i, j := 0, start; j < end+1; i, j = i+1, j+1 {

			if err := ctx.Err(); err != nil {
				return nil, err
			}

			yval := y.Values[j]
			xval := xx.Values[j]

			if dataframe.IsValidFloat64(yval) {
				// Check x val is valid
				if xval != nil {
					yVals = append(yVals, yval)
					xVals = append(xVals, *xval)
				}
			}
		}

		cs.XValues = xVals
		cs.YValues = yVals

		out = cs
	default:
		panic("unrecognized x")
	}

	return out, nil
}
