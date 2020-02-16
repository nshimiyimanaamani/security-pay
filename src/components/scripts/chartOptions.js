const options = {
  chart1: {
    tooltips: { enabled: false },
    hover: { mode: null },
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
    }
  },
  chart2: {
    responsive: true,
    maintainAspectRatio: false,
    cutoutPercentage: 80,
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
  },
  chart3: {
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
                  return label + "%";
              }
            }
          }
        }
      ],
      xAxes: [
        {
          gridLines: { display: false },
          ticks: { fontSize: 12.5 }
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
            label.push(`Percentage: ${tooltipItem.value}%`);
          }
          return label;
        }
      }
    }
  }
};
export default options;
