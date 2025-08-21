import { LineSeries } from "./lightweight-charts.js";

export class SMA {
  constructor(length) {
    this.length = length;
  }

  setup(chart, candleData) {
    const data = [];

    for (let i = 0; i < candleData.length; i++) {
      if (i < this.length) {
        // Provide whitespace data points until the MA can be calculated
        data.push({ time: candleData[i].time });
      } else {
        let sum = 0;
        for (let j = 0; j < this.length; j++) {
          sum += parseFloat(candleData[i - j].close);
        }
        const maValue = sum / this.length;
        data.push({ time: candleData[i].time, value: maValue });
      }
    }

    if (this.series) {
      chart.removeSeries(this.series);
    }

    this.series = chart.addSeries(LineSeries, {
      color: "#2962FF",
      lineWidth: 1,
    });

    this.series.setData(data);
  }

  update(candleData) {
    const last = candleData.length - 1;

    let sum = 0;
    for (let i = 0; i < this.length; i++) {
      sum += parseFloat(candleData[last - i].close);
    }
    const maValue = sum / this.length;

    this.series.update({ time: candleData[last].time, value: maValue });
  }
}
