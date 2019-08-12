<template>
  <div class="table-container">
    <h4 class="title text-center">{{title}}</h4>
    <div class="controllers">
      <b-dropdown id="dropdown-dropright" dropright variant="info">
        <template slot="button-content">Filter by: {{select}}</template>
        <b-dropdown-item @click="setSelected('sector')">sector</b-dropdown-item>
        <b-dropdown-item @click="setSelected('cell')">cell</b-dropdown-item>
        <b-dropdown-item @click="setSelected('village')">village</b-dropdown-item>
      </b-dropdown>
      <b-form-select :options="options" v-model="filter" v-show="select"></b-form-select>
      <b-button @click.prevent="download" class="download btn-info">Download</b-button>
    </div>
    <b-table
      id="data-table"
      bordered
      hover
      small
      :items="items"
      :fields="fields"
      :filter="filter"
      :filterIncludedFields="filterOn"
    >
      <template slot="owner" slot-scope="data">{{ data.value.fname }} {{ data.value.lname }}</template>
    </b-table>
    <pulse-loader class="reports-loader" :loading="loading" :color="color" :size="size"></pulse-loader>
  </div>
</template>

<script>
import jsPDF from "jspdf";
import "jspdf-autotable";
export default {
  name: "reports",
  data() {
    return {
      text: `List of users in`,
      selected: "",
      select: "",
      options: [],
      filterOn: [],
      filter: null,
      loading: false,
      color: "#3db3fa",
      size: "12px",
      fields: {
        owner: {
          label: "Full name",
          sortable: true
        },
        "owner.phone": {
          label: "Phone number",
          sortable: false
        },
        sector: {
          label: "sector",
          sortable: true
        },
        cell: {
          label: "cell",
          sortable: true
        },
        village: {
          label: "village",
          sortable: true
        },
        due: {
          label: "amount",
          sortable: false
        }
      },
      items: []
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    title() {
      return `${this.text} ${this.selected}`;
    },
    activeSector() {
      const Fsector =
        this.$store.getters.getActiveSector.charAt(0).toUpperCase() +
        this.$store.getters.getActiveSector.slice(1);
      return Fsector;
    }
  },
  watch: {
    filter() {
      handler: {
        this.selected = this.filter;
      }
      immediate: true;
    }
  },
  mounted() {
    if (!this.selected) {
      this.selected = this.activeSector;
    }
    this.loadData();
  },
  methods: {
    loadData() {
      this.loading = true;
      this.axios
        .get(
          `${
            this.endpoint
          }/properties/sectors/${this.activeSector.toLowerCase()}?offset=1&limit=10`
        )
        .then(res => {
          this.items = new Array();
          for (const key in res.data.properties) {
            if (res.data.properties.hasOwnProperty(key)) {
              const element = res.data.properties[key];
              this.axios
                .get(`${this.endpoint}/properties/owners/${element.owner}`)
                .then(res => {
                  element.owner = res.data;
                  this.items = [...this.items, element];
                })
                .catch(err => {
                  console.log(err);
                });
            }
          }
          this.loading = false;
        })
        .catch(err => {
          console.log(err);
          this.loading = false;
        });
    },
    download() {
      if (this.items.length <= 0) {
        this.$snotify.error("No data available to download");
      } else if (this.items.length > 0) {
        const doc = new jsPDF();
        const date = new Date();
        const months = [
          "January",
          "February",
          "March",
          "April",
          "May",
          "June",
          "July",
          "August",
          "September",
          "October",
          "November",
          "December"
        ];
        const month = months[date.getMonth()];
        const year = date.getFullYear();
        var pageWidth = doc.internal.pageSize.width;
        doc.setFontSize(12);
        doc.text(`${this.activeSector} sector`, 14, 20);
        doc.text(`on ${date.getDate()} ${month}, ${year}`, pageWidth - 50, 20);
        doc.text(this.title, pageWidth / 3 + 10, 30);
        doc.autoTable({
          html: "#data-table",
          startY: 40,
          showHead: "firstPage",
          bodyStyles: {
            fillColor: [255, 255, 255],
            textColor: 10
          },
          headStyles: {
            fillColor: [255, 255, 255],
            textColor: 10
          },
          styles: {
            lineColor: [0, 0, 0],
            lineWidth: 0.2
          },
          theme: "plain"
        });
        doc.save(`${this.title} of ${month}, ${year}.pdf`);
      }
    },
    setSelected(that) {
      this.select = that;
      this.filterOn = new Array();
      this.filterOn.push(that);
      var options = new Array();
      options.push({ value: null, text: `select ${that}` });
      this.items.forEach(element => {
        if (options.indexOf(element[that]) == -1) {
          options.push(element[that]);
        }
      });
      this.options = options;
    }
  }
};
</script>

<style>
.table-container {
  padding: 40px;
}
.controllers {
  margin: 10px 0;
  height: 30px;
}
.controllers .download {
  float: right;
  outline: none;
}
.reports-loader::after {
  display: block;
  font-weight: bold;
  content: "Please wait...";
  user-select: none;
}
.controllers select {
  padding: 6px;
  margin: 0 10px;
  border-radius: 5px;
  width: 150px;
  outline: none;
}
.controllers button {
  outline: none !important;
}
</style>
