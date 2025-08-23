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

export class EMA {
  constructor(length, color) {
    this.length = length;
    this.color = color;

    this.multiplier = 2 / (length + 1);
    this.emaPrev = null;
    this.emaCurrent = null;
    this.lastProcessedIndex = -1;
  }

  setup(chart, candleData) {
    const data = [];

    for (let i = 0; i < candleData.length; i++) {
      const close = parseFloat(candleData[i].close);

      if (i < this.length - 1) {
        // Provide whitespace data points until the MA can be calculated
        data.push({ time: candleData[i].time });
      } else if (i === this.length - 1) {
        let sum = 0;
        for (let j = 0; j < this.length; j++) {
          sum += parseFloat(candleData[i - j].close);
        }
        this.emaPrev = sum / this.length;
        data.push({ time: candleData[i].time, value: this.emaPrev });
      } else {
        // EMA = (Close - EMA_prev) * multiplier + EMA_prev
        this.emaCurrent =
          (close - this.emaPrev) * this.multiplier + this.emaPrev;

        if (i < candleData.length - 1) {
          this.emaPrev = this.emaCurrent;
        }

        data.push({ time: candleData[i].time, value: this.emaCurrent });
      }

      this.lastProcessedIndex = i;
    }

    if (this.series) {
      chart.removeSeries(this.series);
    }

    this.series = chart.addSeries(LineSeries, {
      color: this.color,
      lineWidth: 1,
    });

    this.series.setData(data);
  }

  update(candleData) {
    const lastIndex = candleData.length - 1;
    const close = parseFloat(candleData[lastIndex].close);

    if (lastIndex > this.lastProcessedIndex) {
      // New candle formed, update emaPrev with the previous completed candle's value
      this.emaPrev = this.emaCurrent;
      this.lastProcessedIndex = lastIndex;
    }

    this.emaCurrent = (close - this.emaPrev) * this.multiplier + this.emaPrev;
    this.series.update({
      time: candleData[lastIndex].time,
      value: this.emaCurrent,
    });
  }
}
