import { createChart, CandlestickSeries } from "./lightweight-charts.js";
import { EMA } from "./indicators.js";

// Define styles for the web component
const elementStyles = `
  :host {
    display: block;
    width: 100%;
  }
  :host([hidden]) {
    display: none;
  }
  .chart-container {
    height: 100%;
    width: 100%;
  }
`;

class LightweightChartWC extends HTMLElement {
  constructor() {
    super();
    // Attach shadow DOM
    this.attachShadow({ mode: "open" });
  }

  // Attribute functions
  get instrument() {
    return this.getAttribute("instrument");
  }

  get timeframe() {
    return this.getAttribute("timeframe");
  }

  static get observedAttributes() {
    return ["instrument", "timeframe"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name === "instrument" || name === "timeframe") {
      this.loadData();
      this.setupSocket();
    }
  }

  connectedCallback() {
    // Create container div
    const container = document.createElement("div");
    container.className = "chart-container";

    // Attach styles
    const style = document.createElement("style");
    style.textContent = elementStyles;

    // Append everything to the shadowRoot
    this.shadowRoot.append(style, container);

    // Create the chart
    this.chart = createChart(container, {
      layout: {
        textColor: "black",
        background: { type: "solid", color: "white" },
      },
      timeScale: {
        timeVisible: true,
      },
    });

    // Add candlestick series
    this.series = this.chart.addSeries(CandlestickSeries, {
      upColor: "#26a69a",
      downColor: "#ef5350",
      borderVisible: false,
      wickUpColor: "#26a69a",
      wickDownColor: "#ef5350",
    });

    // Load indicators
    this.indicators = [new EMA(8, "#FF0000"), new EMA(21, "#2962FF")];

    // Handle resizing
    window.addEventListener("resize", () => {
      this.chart.applyOptions({
        width: container.clientWidth,
        height: container.clientHeight,
      });
    });
  }

  addTick(tick) {
    const createCandle = (time, price) => ({
      time,
      open: price,
      high: price,
      low: price,
      close: price,
    });

    const updateCandle = (candle, val) => ({
      time: candle.time,
      close: val,
      open: candle.open,
      low: Math.min(candle.low, val),
      high: Math.max(candle.high, val),
    });

    let lastCandle = this.series.data().at(-1);

    if (tick.time != lastCandle.time) {
      this.series.update(createCandle(tick.time, tick.price));
    } else {
      this.series.update(updateCandle(lastCandle, tick.price));
    }

    for (let indicator of this.indicators) {
      indicator.update(this.series.data());
    }
  }

  setupSocket() {
    if (this.socket) {
      this.socket.close();
    }

    this.socket = new WebSocket(
      `/data/tick?instrument=${this.instrument}&timeframe=${this.timeframe}`,
    );

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.error) {
          console.error("WebSocket error:", data.error);
          console.error("Details:", data.details);
        } else {
          this.addTick(data);
          this.dispatchEvent(
            new CustomEvent("newtick", {
              detail: {
                value: data.price,
              },
            }),
          );
        }
      } catch (e) {
        console.error("Failed to parse WebSocket message:", e);
      }
    };

    this.socket.onerror = (event) => {
      console.error("WebSocket error:", {
        url: this.socket.url,
        readyState: this.socket.readyState,
        time: new Date().toISOString(),
      });
    };

    this.socket.onclose = (event) => {
      if (event.code === 1006) {
        console.warn("WebSocket closed abnormally (code 1006)");
      }
    };
  }

  loadData() {
    fetch(
      `/data/ohlcv?instrument=${this.instrument}&timeframe=${this.timeframe}`,
    )
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok " + response.statusText);
        }
        return response.json();
      })
      .then((data) => {
        this.series.setData(data);
        this.chart.timeScale().fitContent();

        for (let indicator of this.indicators) {
          indicator.setup(this.chart, this.series.data());
        }
      })
      .catch((error) => {
        console.error("There was a problem with the fetch operation:", error);
      });
  }
}

// Define the custom element
window.customElements.define("lightweight-chart", LightweightChartWC);
