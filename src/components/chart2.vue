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
      Chart.defaults.global.legend.position = "right";
      let chart1 = new Chart(Chart1Container, {
        type: "doughnut",
        data: {
          labels: ["Umutekano", ""],
          datasets: [
            {
              label: "RWF",
              data: [5, 3],
              backgroundColor: ["#58C5AD", "#F9F9F9"]
            }
          ]
        },
        options: {
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
            label: {
              fontsize: 23
            },
            tooltip: {
              display: false
            }
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
      });
    }
  }
};
</script>


<style>
@import url("../assets/css/chart2.css");
</style>

