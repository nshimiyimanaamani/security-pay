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
  mounted() {
    this.drawChart();
  },
  methods: {
    drawChart() {
      let Chart3Container = document.getElementById("Chart-3").getContext("2d");
      Chart.defaults.global.defaultFontSize = 20;
      Chart.defaults.scale.ticks.beginAtZero = true;
      console.log(Chart);
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
              data: [3, 4, 2, 1, 3, 4, 5],
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
              data: [3, 4, 2, 1, 3, 4, 5],
              backgroundColor: "#88D8C6",
              borderColor: "#7fc7b7",
              borderWidth: 1
            }
          ]
        },
        options: {
          tooltips: {
            filter: function(tooltipItem) {
              return tooltipItem.datasetIndex === 0;
            }
          },
          scales: {
            yAxes: [
              {
                ticks: {
                  callback: function(value, index, values) {
                    console.log(value, index, values);
                    return value + "%";
                  },
                  max: 5,
                  min: 0,
                  stepSize: 1
                }
              }
            ],
            xAxes: [{ barPercentage: 1, categoryPercentage: 1 }]
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
