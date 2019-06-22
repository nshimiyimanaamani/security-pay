<template>
  <div id="chart1container">
    <div class="chart1Title">
      <i class="fa fa-th-large"></i>
      <h1 class>Remera collecting acct.</h1>
      <span class="fa fa-cog"></span>
    </div>
    <div class="chart1">
      <canvas id="Chart-1"></canvas>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      datas: [2, 4, 3],
      labels: ["BK Acc", "MTN", "AIRTEL"]
    };
  },
  mounted() {
    this.drawChart();
  },
  methods: {
    drawChart(datas) {
      let Chart1Container = document.getElementById("Chart-1").getContext("2d");
      Chart.defaults.global.defaultFontSize = 18;
      Chart.defaults.global.defaultFontFamily = "Helvetica Neue, sans-serif";
      let chart1 = new Chart(Chart1Container, {
        type: "bar",
        data: {
          labels: this.labels,
          datasets: [
            {
              data: this.datas,
              backgroundColor: "#58c5ad"
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
              fontsize: 23
            }
          },
          layout: {
            padding: {
              left: 15,
              right: 15,
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
@import url("../assets/css/chart1.css");
</style>
