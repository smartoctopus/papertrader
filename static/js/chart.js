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

  // Attribute functions
  get instrument() {
    return this.getAttribute("instrument");
  }

  static get observedAttributes() {
    return ["instrument"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name === "instrument") {
      this.loadData();
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
    });

    // Add a simple line series with example data
    this.series = this.chart.addSeries(CandlestickSeries, {
      upColor: "#26a69a",
      downColor: "#ef5350",
      borderVisible: false,
      wickUpColor: "#26a69a",
      wickDownColor: "#ef5350",
    });

    // Handle resizing
    window.addEventListener("resize", () => {
      this.chart.applyOptions({
        width: container.clientWidth,
        height: container.clientHeight,
      });
    });
  }

  loadData() {
    fetch(`/data/ohlcv?instrument=${this.instrument}`)
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok " + response.statusText);
        }
        return response.json();
      })
      .then((data) => {
        console.log(data);

        this.series.setData(data);
        this.chart.timeScale().fitContent();
      })
      .catch((error) => {
        console.error("There was a problem with the fetch operation:", error);
      });
  }
}

// Define the custom element
window.customElements.define("lightweight-chart", LightweightChartWC);
