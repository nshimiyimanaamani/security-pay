<script>
import { Bar, mixins } from "vue-chartjs";
const { reactiveProp } = mixins;
export default {
  extends: Bar,
  mixins: [reactiveProp],
  data() {
    return {
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          yAxes: [
            {
              ticks: {
                min: 0,
                stepSize: 10,
                callback: (label, index, labels) => {
                  switch (label) {
                    case label:
                      return label + "M";
                  }
                }
              }
            }
          ],
          xAxes: [
            {
              barPercentage: 0.95,
              categoryPercentage: 1,
              gridLines: { display: false }
            }
          ]
        },
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
        },
        tooltips: {
          displayColors: false,
          callbacks: {
            label: function(tooltipItem, data) {
              var label = [data.datasets[tooltipItem.datasetIndex].label] || "";

              if (label) {
                label = new Array();
                label.push(`Abishyuye: 30`);
                label.push("Abatarishyura: 40");
              }
              return label;
            }
          }
        }
      }
    };
  },
  mounted() {
    this.renderChart(this.chartData, this.options);
  }
};
</script>
