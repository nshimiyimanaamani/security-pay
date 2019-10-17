<script>
import { Doughnut, mixins } from "vue-chartjs";
const { reactiveProp } = mixins;
export default {
  extends: Doughnut,
  mixins: [reactiveProp],
  data() {
    return {
      options: {
        responsive: true,
        maintainAspectRatio: false,
        hover: { mode: null },
        tooltips: { enabled: false },
        legend: { display: false },
        layout: {
          padding: {
            left: 10,
            right: 10,
            top: 10,
            bottom: 10
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