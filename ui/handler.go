package ui

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func (s *server) handleGetIndex() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				Theme: types.ThemeInfographic,
				Width: "80%",
			}),
			charts.WithTitleOpts(opts.Title{
				Title:    "Total Bytes Received/Send",
				Subtitle: "Random data generated for testing purpose only!",
			}),
			charts.WithTooltipOpts(opts.Tooltip{
				Show: true,
			}),
			charts.WithLegendOpts(opts.Legend{
				Show: true,
			}),
		)

		line.SetXAxis([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24"}).
			AddSeries("Bytes Received", generateLineItems()).
			AddSeries("Bytes Send", generateLineItems()).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{
					Smooth: false,
				}),
				charts.WithAreaStyleOpts(opts.AreaStyle{
					Opacity: 0.2,
				}),
			)
		if err := line.Render(writer); err != nil {
			log.Println(err.Error())
		}
	}
}

func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i <= 24; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(30000)}) //nolint:gosec
	}
	return items
}
