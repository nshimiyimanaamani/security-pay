<script>
import { Doughnut, mixins } from "vue-chartjs";
const { reactiveProp } = mixins;
export default {
  extends: Doughnut,
  mixins: [reactiveProp],
  props: ["options"],
  mounted() {
    Chart.defaults.global.defaultFontFamily = '"Montserrat", sans-serif';
    Chart.defaults.global.defaultFontColor = "#2d2d2d";
    Chart.defaults.global.defaultFontSize = 12.3;
    this.renderChart(this.chartData, this.options);
  }
};
</script>