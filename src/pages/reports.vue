<template>
  <div class="table-container">
    <h4 class="title text-center" v-show="selected">{{title}}</h4>
    <div class="controllers">
      <b-dropdown
        id="dropdown-dropright"
        dropright
        variant="info"
        ref="dropdown"
        class="filter-dropdown"
      >
        <template slot="button-content">Filter By</template>
        <b-dropdown-form>
          <b-form-group label="sector">
            <b-form-select v-model="select.sector" :options="select.sectorOptions">
              <option value>select sector</option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="cell" v-show="select.sector">
            <b-form-select v-model="select.cell" :options="select.cellOptions">
              <option value>select cell</option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="village" v-show="select.sector && select.cell">
            <b-form-select v-model="select.village" :options="select.villageOptions">
              <option value>select village</option>
            </b-form-select>
          </b-form-group>
          <b-dropdown-divider></b-dropdown-divider>
          <b-button variant="primary" size="sm" @click.prevent="tableItems = filter()">Ok</b-button>
          <b-button variant="danger" size="sm" @click.prevent="clearFilter">Clear</b-button>
        </b-dropdown-form>
      </b-dropdown>
      <div class="search">
        <b-form-input
          placeholder="search user..."
          size="sm"
          v-model="search.name"
          list="search-datalist-id"
        ></b-form-input>
        <b-button variant="info" @click="search.name = ''">
          <i class="fa fa-times"></i>
        </b-button>
        <datalist id="search-datalist-id">
          <option v-for="name in search.datalist" :key="name">{{ name }}</option>
        </datalist>
      </div>
      <b-button @click.prevent="download" class="download btn-info">Download</b-button>
    </div>
    <b-table id="data-table" bordered striped hover small :items="tableItems" :fields="fields">
      <template v-slot:cell(due)="data">{{data.item.due}} Rwf</template>
      <template v-slot:cell(index)="data">
        <article style="text-align:center">{{data.index + 1}}</article>
      </template>
    </b-table>
    <pulse-loader class="reports-loader" :loading="loading.request" :color="color" :size="size"></pulse-loader>
  </div>
</template>

<script>
import jsPDF from "jspdf";
import "jspdf-autotable";
export default {
  name: "reports",
  data() {
    return {
      selected: "",
      width: 0,
      options: [],
      color: "#333333bd",
      loading: {
        progress: false,
        request: false
      },
      search: {
        name: "",
        datalist: []
      },
      select: {
        sector: "",
        cell: "",
        village: "",
        sectorOptions: [],
        cellOptions: [],
        villageOptions: []
      },
      size: "5px",
      fields: [
        {
          key: "index",
          label: "NO"
        },
        {
          key: "owner",
          label: "Full name",
          sortable: true
        },
        {
          key: "sector",
          label: "sector",
          sortable: true
        },
        {
          key: "cell",
          label: "cell",
          sortable: true
        },
        {
          key: "village",
          label: "village",
          sortable: true
        },
        {
          key: "due",
          label: "amount",
          sortable: false
        }
      ],
      items: [],
      tableItems: []
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    title() {
      return `List of users in ${this.selected}`;
    },
    activeSector() {
      const Fsector =
        this.$store.getters.getActiveSector.charAt(0).toUpperCase() +
        this.$store.getters.getActiveSector.slice(1);
      return Fsector;
    }
  },
  watch: {
    items() {
      handler: {
        this.tableItems = this.filter();
        this.select.sectorOptions = new Array();
        if (this.items.length > 0) {
          this.selected = this.items[0].sector;
        }
        this.items.forEach(element => {
          if (this.select.sectorOptions.indexOf(element.sector) == -1) {
            this.select.sectorOptions.push(element.sector);
          }
        });
      }
    },
    "select.sector"() {
      handler: {
        if (this.select.sector) {
          this.select.cellOptions = [];
          const cellOptions = this.items.filter(
            sec => sec.sector == this.select.sector
          );
          cellOptions.forEach(element => {
            if (this.select.cellOptions.indexOf(element.cell) == -1) {
              this.select.cellOptions = [
                ...this.select.cellOptions,
                element.cell
              ];
            }
          });
        } else {
          this.select.cellOptions = [];
        }
      }
    },
    "select.cell"() {
      handler: {
        if (this.select.cell) {
          this.select.villageOptions = [];
          const villageOptions = this.items.filter(
            res => res.cell == this.select.cell
          );
          villageOptions.forEach(element => {
            if (this.select.villageOptions.indexOf(element.village) == -1) {
              this.select.villageOptions = [
                ...this.select.villageOptions,
                element.village
              ];
            }
          });
        } else {
          this.select.villageOptions = [];
        }
      }
    },
    "search.name"() {
      handler: {
        this.search.datalist = new Array();
        this.tableItems = this.filter().filter(obj => {
          this.search.datalist.push(obj.owner);
          return obj.owner
            .toLowerCase()
            .includes(this.search.name.toLowerCase());
        });
        while (this.search.datalist.length > 5) {
          this.search.datalist.pop();
        }
      }
    }
  },
  mounted() {
    this.loadData();
  },
  methods: {
    loadData() {
      this.loading.request = true;
      this.axios
        .get(
          `${
            this.endpoint
          }/properties/sectors/${this.activeSector.toLowerCase()}?offset=1&limit=100`
        )
        .then(res => {
          this.items = new Array();
          this.loading.request = false;
          this.items = res.data.properties;
          console.log(this.items);
        })
        .catch(err => {
          console.log(err);
          this.loading.request = false;
        });
    },
    filter() {
      this.$refs.dropdown.hide(true);
      return this.items.filter(item => {
        return (
          item.sector
            .toLowerCase()
            .includes(this.select.sector.toLowerCase()) &&
          item.cell.toLowerCase().includes(this.select.cell.toLowerCase()) &&
          item.village.toLowerCase().includes(this.select.village.toLowerCase())
        );
      });
    },
    clearFilter() {
      this.select.sector = "";
      this.select.cell = "";
      this.select.village = "";
      this.tableItems = this.items;
      this.selected = this.activeSector;
      this.$refs.dropdown.hide(true);
    },
    download() {
      if (this.loading.progress) {
        this.$snotify.error("Please wait while the list is being completed");
      } else if (!this.loading.progress && !this.loading.request) {
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
  display: flex;
}
.controllers .download {
  margin-left: 10px;
  outline: none;
  padding: 2px 10px;
  height: fit-content;
}
.reports-loader::before {
  display: inline;
  font-weight: bold;
  content: "loading";
  user-select: none;
}
.controllers button {
  outline: none !important;
}
.progress {
  border-radius: 10px;
  height: 10px;
  width: 50%;
  margin: auto;
}
.progress-bar {
  font-size: 10px;
  line-height: 11px;
  background: #4394da;
}
.subtitle {
  font-weight: bold;
  width: fit-content;
  margin-left: auto;
  font-size: 13px;
  margin-bottom: auto;
}
.filter-dropdown legend {
  font-size: 14px;
  font-weight: bold;
}
.filter-dropdown select {
  width: 100%;
  font-size: 14px;
  margin: 0;
  height: auto;
  border-radius: 3px;
  border-color: #cacaca;
}
.filter-dropdown .form-group {
  margin-bottom: 10px;
}
.filter-dropdown hr {
  margin-top: 10px;
  margin-bottom: 10px;
}
.filter-dropdown button {
  float: right;
  height: fit-content;
  padding: 2px 10px;
  font-size: 16px;
}
.filter-dropdown .dropdown-menu {
  min-width: 200px;
  margin: 0 2px 0;
}
.dropdown-menu form button {
  font-size: 12px !important;
  padding: 4px 10px !important;
  margin-left: 6px;
  width: 45%;
}
.dropdown-menu form {
  outline: none !important;
}
.table-container .controllers .search {
  display: flex;
  margin-left: auto;
}
.table-container .controllers .search input {
  height: 100%;
  border-radius: 5px 0 0 5px;
}
.table-container .controllers .search > button {
  border-radius: 0 4px 4px 0;
  color: white;
  height: fit-content;
  padding: 2px 10px;
}
table thead th {
  font-size: 14px;
}
table td {
  text-transform: capitalize;
  font-size: 14px;
}
</style>
