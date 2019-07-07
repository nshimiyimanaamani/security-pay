<template>
  <div id="chart2container">
    <div class="chart2Title">
      <i class="fa fa-th-large"></i>
      <h1>Remera total collected</h1>
      <span class="fa fa-cog"></span>
    </div>
    <div class="chart2">
      <canvas id="Chart-2"></canvas>
      <div id="legend"></div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      Percentage: 60
    };
  },
  mounted() {
    this.drawChart();
  },
  methods: {
    drawChart() {
      let chartcanvas = document.getElementById("Chart-2").getContext("2d");
      Chart.defaults.global.defaultFontSize = 15;
      Chart.defaults.global.legend.position = "right";
      Chart.defaults.global.legend.labels.boxWidth = 0;
      Chart.defaults.global.tooltips.enabled = false;
      var value = this.Percentage;
      var chartData = {
        type: "doughnut",
        data: {
          labels: [`UMUTEKANO: ${value}% `],
          datasets: [
            {
              data: [value, 100 - value],
              backgroundColor: ["#58C5AD", "#f9f9f9"],
              hoverBackgroundColor: ["#58C5AD", "#f9f9f9"],
              hoverBorderColor: ["#58C5AD", "#ffffff"]
            }
          ]
        },
        options: {
          elements: {
            center: {
              text: `${value}%`,
              color: "#58C5AD", // Default is #000000
              fontStyle: "Arial", // Default is Arial
              sidePadding: 20 // Defualt is 20 (as a percentage)
            }
          },
          cutoutPercentage: 80,
          scales: {
            ticks: {
              display: false,
              gridLines: {
                display: false
              }
            }
          },
          maintainAspectRatio: false,
          legend: {
            display: false
          },
          legendCallback: function(chart) {
            // Return the HTML string here.
            console.log(chart.data.datasets[0]);
            var text = [];
            text.push('<ul class="' + chart.id + '-legend">');
            for (var i = 0; i < chart.data.datasets[0].data.length; i++) {
              text.push(
                '<li><span id="legend-' +
                  i +
                  '-item" style="background-color:' +
                  chart.data.datasets[0].backgroundColor[i] +
                  '"   onclick="updateDataset(event, ' +
                  "'" +
                  i +
                  "'" +
                  ')">'
              );
              if (chart.data.labels[i]) {
                text.push(chart.data.labels[i]);
              }
              text.push("</span></li>");
            }
            text.push("</ul>");
            return text.join("");
          },
          layout: {
            padding: {
              left: 10,
              right: 15,
              top: 15,
              bottom: 15
            }
          }
        }
      };
      window.chart = new Chart(chartcanvas, chartData);
      chart.generateLegend();
      Chart.pluginService.register({
        beforeDraw: function(chart) {
          if (chart.config.options.elements.center) {
            var ctx = chart.chart.ctx;
            var centerConfig = chart.config.options.elements.center;
            var fontStyle = centerConfig.fontStyle || "arial";
            var txt = centerConfig.text;
            var color = centerConfig.color || "#000";
            var sidePadding = centerConfig.sidePadding || 20;
            var sidePaddingCalculated =
              (sidePadding / 100) * (chart.innerRadius * 2);
            ctx.font = "45px " + fontStyle;
            var stringWidth = ctx.measureText(txt).width;
            var elementWidth = chart.innerRadius * 2 - sidePaddingCalculated;
            var widthRatio = elementWidth / stringWidth;
            var newFontSize = Math.floor(30 * widthRatio);
            var elementHeight = chart.innerRadius * 2;
            var fontSizeToUse = Math.min(newFontSize, elementHeight);
            ctx.textAlign = "center";
            ctx.textBaseline = "middle";
            var centerX = (chart.chartArea.left + chart.chartArea.right) / 2;
            var centerY = (chart.chartArea.top + chart.chartArea.bottom) / 2;
            ctx.font = fontSizeToUse + "px " + fontStyle;
            ctx.fillStyle = color;
            ctx.fillText(txt, centerX, centerY);
          }
        }
      });
    }
  }
};
</script>


<style>
@import url("../assets/css/chart2.css");
</style>

