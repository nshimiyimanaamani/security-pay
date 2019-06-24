<template>
  <div id="chart3container">
    <div class="chart3Title">
      <i class="fa fa-th-large"></i>
      <h1>Remera sector</h1>
      <span class="fa fa-cog"></span>
    </div>
    <div class="chart3">
      <canvas id="Chart-3"></canvas>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      pay: {
        CELLA: [3, 2, 1, 20],
        CELLB: [3, 2, 1, 40],
        CELLC: [3, 2, 1, 30],
        CELLD: [3, 2, 1, 50],
        CELLE: [3, 2, 1, 20],
        CELLF: [3, 2, 1, 10],
        CELLG: [3, 2, 1, 30]
      }
    };
  },
  mounted() {
    this.drawChart(this.pay);
  },
  methods: {
    drawChart(payData) {
      let Chart3Container = document.getElementById("Chart-3").getContext("2d");
      Chart.defaults.global.defaultFontSize = 20;
      Chart.defaults.scale.ticks.beginAtZero = true;

      var customTooltips = function(tooltip) {
        // Tooltip Element
        var tooltipEl = document.getElementById("chartjs-tooltip");
        if (!tooltipEl) {
          tooltipEl = document.createElement("div");
          tooltipEl.id = "chartjs-tooltip";
          tooltipEl.innerHTML = "<table></table>";
          this._chart.canvas.parentNode.appendChild(tooltipEl);
        }
        // Hide if no tooltip
        if (tooltip.opacity === 0) {
          tooltipEl.style.opacity = 0;
          return;
        }
        // Set caret Position
        tooltipEl.classList.remove("no-transform");
        if (tooltip.yAlign) {
          tooltipEl.classList.add(tooltip.yAlign);
        } else {
          tooltipEl.classList.add("no-transform");
        }
        function getBody(bodyItem) {
          return bodyItem.lines;
        }
        if (getBody) {
          var innerHtml = "<thead>";
          var title = tooltip.title.toString() || "";
          innerHtml += "<tr><th>" + title + "</th></tr>";
          innerHtml += "</thead><tbody>";
          var index = tooltip.dataPoints[0].index;
          var body = Object.values(payData);
          var key = [
            "Fully payed:",
            "Half payed:",
            "Not Payed:",
            "Total percentage:"
          ];
          var span = '<span class="chartjs-tooltip-key"></span>';
          for (var i = 0; i < key.length; i++) {
            innerHtml += "<tr><td>" + key[i] + body[index][i] + "</td></tr>";
          }
          innerHtml += "</tbody>";
          var tableRoot = tooltipEl.querySelector("table");
          tableRoot.innerHTML = innerHtml;
        }
        var positionY = this._chart.canvas.offsetBottom;
        var positionX = this._chart.canvas.offsetLeft;
        // Display, position, and set styles for font
        tooltipEl.style.opacity = 1;
        tooltipEl.style.left = positionX + tooltip.caretX + "px";
        tooltipEl.style.top =
          tooltip.dataPoints[0].y - tooltip.height * 2 + "px";
        tooltipEl.style.fontFamily = tooltip._bodyFontFamily;
        tooltipEl.style.fontSize = tooltip.bodyFontSize + "px";
        tooltipEl.style.fontStyle = tooltip._bodyFontStyle;
        tooltipEl.style.padding =
          tooltip.yPadding + "px " + tooltip.xPadding + "px";
      };
      let chart3 = new Chart(Chart3Container, {
        type: "bar",
        data: {
          labels: [
            "CELL A",
            "CELL B",
            "CELL C",
            "CELL D",
            "CELL E",
            "CELL F",
            "CELL G"
          ],
          datasets: [
            {
              type: "line",
              lineTension: 0,
              label: "Rwf ",
              data: [
                payData.CELLA[3],
                payData.CELLB[3],
                payData.CELLC[3],
                payData.CELLD[3],
                payData.CELLE[3],
                payData.CELLF[3],
                payData.CELLG[3]
              ],
              backgroundColor: "transparent",
              borderColor: "#46A9D4",
              borderWidth: 5,
              borderDash: [5, 5],
              pointBackgroundColor: "#fff",
              pointBorderWidth: 1,
              pointRadius: 8,
              pointHoverRadius: 8
            },
            {
              label: "RWF",
              data: [
                payData.CELLA[3],
                payData.CELLB[3],
                payData.CELLC[3],
                payData.CELLD[3],
                payData.CELLE[3],
                payData.CELLF[3],
                payData.CELLG[3]
              ],
              backgroundColor: "#58c5ad",
              borderColor: "#58c5ad",
              borderWidth: 1
            }
          ]
        },
        options: {
          tooltips: {
            // filter: function(tooltipItem) {
            //   return tooltipItem.datasetIndex === 0;
            // },
            enabled: false,
            mode: "index",
            position: "nearest",
            custom: customTooltips,
            displayColors: false
          },

          scales: {
            yAxes: [
              {
                ticks: {
                  callback: function(value, index, values) {
                    return value + " ";
                  },
                  max: 100,
                  min: 0,
                  stepSize: 20
                }
              }
            ],
            xAxes: [
              {
                barPercentage: 0.99,
                categoryPercentage: 1,
                gridLines: { display: false }
              }
            ]
          },
          maintainAspectRatio: false,
          legend: {
            display: false,
            label: {
              fontsize: 20
            },
            tooltip: {
              display: false
            }
          },
          layout: {
            padding: {
              left: 20,
              right: 50,
              top: 25,
              bottom: 0
            }
          }
        }
      });
    }
  }
};
</script>


<style>
@import url("../assets/css/chart3.css");
</style>
