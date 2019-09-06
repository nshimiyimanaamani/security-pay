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
              <option :value="null">select sector</option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="cell" v-show="select.sector">
            <b-form-select v-model="select.cell" :options="select.cellOptions">
              <option :value="null">select cell</option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="village" v-show="select.sector && select.cell">
            <b-form-select v-model="select.village" :options="select.villageOptions">
              <option :value="null">select village</option>
            </b-form-select>
          </b-form-group>
          <b-dropdown-divider></b-dropdown-divider>
          <b-button variant="primary" size="sm" @click.prevent="filter">Ok</b-button>
          <b-button variant="danger" size="sm" @click.prevent="clearFilter">Clear</b-button>
        </b-dropdown-form>
      </b-dropdown>
      <b-button @click.prevent="download" class="download btn-info">Download</b-button>
    </div>
    <b-table id="data-table" bordered striped hover small :items="tableItems" :fields="fields">
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
      select: {
        sector: null,
        cell: null,
        village: null,
        sectorOptions: [],
        cellOptions: [],
        villageOptions: []
      },
      size: "5px",
      fields: {
        owner: {
          label: "Full name",
          sortable: true
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
        this.tableItems = this.items;
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
          console.log(res.data);
          this.loading.request = false;
          this.items = res.data.properties
          // for (const key in res.data.properties) {
          //   if (res.data.properties.hasOwnProperty(key)) {
          //     const element = res.data.properties[key];
          //     this.axios
          //       .get(`${this.endpoint}/properties/owners/${element.owner}`)
          //       .then(result => {
          //         element.owner = result.data;
          //         this.items = [...this.items, element];
          //         const width = (
          //           (this.items.length * 100) /
          //           res.data.properties.length
          //         ).toFixed();
          //         for (let i = 0; i <= width; i++) {
          //           if (i > this.width) {
          //             this.width = i;
          //           }
          //         }
          //         if (res.data.properties.length == this.items.length) {
          //           this.loading.request = false;
          //           setTimeout(() => {
          //             this.loading.progress = false;
          //           }, 1000);
          //         }
          //       })
          //       .catch(err => {
          //         console.log(err);
          //         this.loading.progress = false;
          //       });
          //   }
          // }
        })
        .catch(err => {
          console.log(err);
          this.loading.request = false;
        });
    },
    filter() {
      if (this.select.sector) {
        if (this.select.cell) {
          if (this.select.village) {
            const filtered = this.items.filter(
              item =>
                item.sector == this.select.sector &&
                item.cell == this.select.cell &&
                item.village == this.select.village
            );
            this.tableItems = filtered;
            this.selected = this.select.village;
          } else if (!this.select.village) {
            const filtered = this.items.filter(
              item =>
                item.sector == this.select.sector &&
                item.cell == this.select.cell
            );
            this.tableItems = filtered;
            this.selected = this.select.cell;
          }
        } else if (!this.select.cell) {
          const filtered = this.items.filter(
            item => item.sector == this.select.sector
          );
          this.tableItems = filtered;
          this.selected = this.select.sector;
        }
        this.$refs.dropdown.hide(true);
      } else if (!this.select.sector) {
        this.tableItems = this.items;
        this.$refs.dropdown.hide(true);
      }
    },
    clearFilter() {
      this.select.sector = null;
      this.select.cell = null;
      this.select.village = null;
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
  margin-left: auto;
  outline: none;
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
  margin-bottom: 5px;
  font-weight: bold;
}
.filter-dropdown select {
  width: 100%;
  margin: 0;
  height: 30px;
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
}
.filter-dropdown .dropdown-menu {
  min-width: 200px;
  margin: 0 2px 0;
}
.dropdown-menu form button {
  font-size: 12px !important;
  padding: 4px 10px !important;
  margin-left: 6px;
}
.dropdown-menu form {
  outline: none !important;
}
</style>
