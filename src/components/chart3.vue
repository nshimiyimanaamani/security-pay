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
      let chart3 = new Chart(Chart3Container, {
        type: "bar",
        data: {
          labels: ["BK Acc", "MTN", "AIRTEL"],
          datasets: [
            {
              label: "RWF",
              data: [1, 3, 4],
              backgroundColor: ["grey", "yellow", "red"]
            }
          ]
        },
        options: {
          scales: {
            yAxes: [
              {
                ticks: {
                  max: 4,
                  min: 0,
                  stepSize: 1,
                  callback: function(label, index, labels) {
                    switch (label) {
                      case 0:
                        return "0";
                      case 1:
                        return "50M";
                      case 2:
                        return "100M";
                      case 3:
                        return "150M";
                      case 4:
                        return "200M";
                    }
                  }
                }
              }
            ]
          },
          maintainAspectRatio: false,
          legend: {
            display: false,
            label: {
              fontsize: 23
            },
            tooltip: {
              display: false
            }
          },
          layout: {
            padding: {
              left: 30,
              right: 40,
              top: 55,
              bottom: 20
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
