<template>
  <div id="chart2container">
    <div class="chart2Title">
      <i class="fa fa-th-large"></i>
      <h1>Remera total collected</h1>
      <span class="fa fa-cog"></span>
    </div>
    <div class="chart2">
      <canvas id="Chart-2"></canvas>
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
      let Chart1Container = document.getElementById("Chart-2").getContext("2d");
      Chart.defaults.global.defaultFontSize = 20;
      let chart1 = new Chart(Chart1Container, {
        type: "pie",
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
                  display: false
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
@import url("../assets/css/chart2.css");
</style>

