<template>
  <b-card id="chart1container">
    <b-card-header class="chart1Title">
      <i class="fa fa-th-large"></i>
      <h1 class>Remera collecting acct.</h1>
      <span class="fa fa-cog"></span>
    </b-card-header>
    <div class="chart1">
      <canvas id="Chart-1"></canvas>
    </div>
  </b-card>
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
    if(this.$route.name == "dashboard"){
    this.drawChart();
    }
  },
  methods: {
    drawChart(datas) {
      let Chart1Container = document.getElementById("Chart-1").getContext("2d");
      Chart.defaults.global.defaultFontSize = 15;
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
          tooltips: { enabled: false },
          hover: { mode: null },
          legend: false,
          plugins: {
            datalabels: {
              align: "end",
              anchor: "end",
              color: "white",
              font: {
                size: 13,
                weight: 600
              },
              offset: 4,
              padding: 5,
              formatter: function(value) {
                return Math.round(value * 10) / 10;
              }
            }
          },
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
            display: false
          },
          layout: {
            padding: {
              left: 10,
              right: 20,
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


<style scoped>
@import url("../assets/css/chart1.css");
</style>
