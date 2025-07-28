import { createChart, CandlestickSeries } from "./lightweight-charts.js";

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
    });

    // Add a simple line series with example data
    this.series = this.chart.addSeries(CandlestickSeries, {
      upColor: "#26a69a",
      downColor: "#ef5350",
      borderVisible: false,
      wickUpColor: "#26a69a",
      wickDownColor: "#ef5350",
    });
    const data = [
      { open: 10, high: 10.63, low: 9.49, close: 9.55, time: 1642427876 },
      { open: 9.55, high: 10.3, low: 9.42, close: 9.94, time: 1642514276 },
      { open: 9.94, high: 10.17, low: 9.92, close: 9.78, time: 1642600676 },
      { open: 9.78, high: 10.59, low: 9.18, close: 9.51, time: 1642687076 },
      { open: 9.51, high: 10.46, low: 9.1, close: 10.17, time: 1642773476 },
      { open: 10.17, high: 10.96, low: 10.16, close: 10.47, time: 1642859876 },
      { open: 10.47, high: 11.39, low: 10.4, close: 10.81, time: 1642946276 },
      { open: 10.81, high: 11.6, low: 10.3, close: 10.75, time: 1643032676 },
      { open: 10.75, high: 11.6, low: 10.49, close: 10.93, time: 1643119076 },
      { open: 10.93, high: 11.53, low: 10.76, close: 10.96, time: 1643205476 },
    ];
    this.series.setData(data);

    this.chart.timeScale().fitContent();

    // Handle resizing
    window.addEventListener("resize", () => {
      this.chart.applyOptions({
        width: container.clientWidth,
        height: container.clientHeight,
      });
    });
  }

  disconnectedCallback() {
    // Cleanup listeners or destroy chart if needed
  }
}

// Define the custom element
window.customElements.define("lightweight-chart", LightweightChartWC);
