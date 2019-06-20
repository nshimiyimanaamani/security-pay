<template>
  <div id="chart1container">
    <div class="chart1Title">
      <i class="fa fa-th-large"></i>
      <h1>Remera collecting acct.</h1>
      <span class="fa fa-cog"></span>
    </div>
    <div class="chart1">
      <canvas id="Chart-1"></canvas>
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
      let Chart1Container = document.getElementById("Chart-1").getContext("2d");
      Chart.defaults.global.defaultFontSize = 20;
      let chart1 = new Chart(Chart1Container, {
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
              left: 50,
              right: 15,
              top: 55,
              bottom: 15
            }
          }
        }
      });
    }
  }
};
</script>


<style>
@import url("../assets/css/chart1.css");
</style>
