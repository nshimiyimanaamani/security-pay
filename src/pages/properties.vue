<template>
  <b-container class="table-container px-5 max-width">
    <vue-title title="Paypack | Properties" />
    <h4 class="title text-center">
      {{title}}
      <b-button class="add-property mb-1 font-15" variant="info" @click="addProperty.show = true">
        <i class="fas fa-plus-circle"></i> Property
      </b-button>
    </h4>
    <hr />
    <b-row class="my-1 align-items-end px-3">
      <b-dropdown
        id="dropdown-dropright"
        dropright
        variant="info"
        ref="dropdown"
        class="filter-dropdown mr-auto"
      >
        <template slot="button-content">Filter By</template>
        <b-dropdown-form>
          <b-card-body class="p-2">
            <b-form-group label="cell">
              <b-form-select v-model="select.cell" :options="cellOptions">
                <template v-slot:first>
                  <option :value="null" disabled>Select cell</option>
                </template>
              </b-form-select>
            </b-form-group>
            <b-form-group label="village" v-if="select.cell">
              <b-form-select v-model="select.village" :options="villageOptions">
                <template v-slot:first>
                  <option :value="null" disabled>Select village</option>
                </template>
              </b-form-select>
            </b-form-group>
          </b-card-body>
          <b-card-body class="p-2">
            <b-form-group>
              <template v-slot:label>
                <b>Choose columns to display:</b>
                <br />
                <b-form-checkbox
                  v-model="select.selectAll"
                  aria-describedby="columns"
                  aria-controls="columns"
                  @change="allSelected"
                >{{ (select.selectAll) ? 'Un-select All' : 'Select All' }}</b-form-checkbox>
              </template>
              <b-form-checkbox-group
                id="columns"
                v-model="select.shownColumn"
                :options="columns"
                size="sm"
                name="columns"
                stacked
              ></b-form-checkbox-group>
            </b-form-group>
          </b-card-body>
        </b-dropdown-form>
        <b-button variant="primary" size="sm" @click.prevent="tableItems = filter()">Ok</b-button>
        <b-button variant="danger" size="sm" @click.prevent="clearFilter">Clear</b-button>
      </b-dropdown>

      <div>
        <b-form-input
          placeholder="search user..."
          class="rounded-2 font-15"
          type="search"
          size="sm"
          v-model="search.name"
          list="search-user-id"
        ></b-form-input>
        <datalist id="search-user-id">
          <option v-for="name in search.datalist" :key="name">{{ name }}</option>
        </datalist>
      </div>

      <b-button @click="loadData" variant="info" size="sm" class="ml-1 font-15">Refresh</b-button>
    </b-row>

    <b-table
      id="data-table"
      bordered
      striped
      hover
      small
      responsive
      :items="tableItems"
      :fields="fields"
      :busy="loading.request"
      :sort-by.sync="sortBy"
      :show-empty="!loading.request"
      :current-page="pagination.currentPage"
      :per-page="pagination.perPage"
      @row-contextmenu="editHouse"
    >
      <template v-slot:cell(due)="data">{{Number(data.item.due).toLocaleString()}} Rwf</template>
      <template v-slot:cell(owner)="data">{{data.item.owner.fname +" "+ data.item.owner.lname}}</template>
      <template v-slot:cell(occupied)="data">
        <b-card-text
          class="text-right m-0"
        >{{data.item.occupied ? data.item.occupied?'Yes':'No' : "No"}}</b-card-text>
      </template>
      <template v-slot:cell(index)="data">
        <article class="text-center">{{data.index + 1}}</article>
      </template>
      <template v-slot:table-busy>
        <div class="text-center my-2">
          <b-spinner small class="align-middle" />&nbsp;
          <strong>Loading...</strong>
        </div>
      </template>
      <template v-slot:empty>
        <h5
          class="text-center font-15 my-4"
        >{{search.name ? search.name+' "is not in the list"':'No Property Found!'}}</h5>
      </template>
      <template v-slot:custom-foot v-if="!loading.request">
        <b-tr v-if="select.shownColumn.includes('Amount')">
          <b-td v-for="index in select.shownColumn" :key="index" variant="primary">
            <div
              v-if="index == select.shownColumn[select.shownColumn.indexOf('Amount')-1]"
              class="text-danger"
            >
              <strong>TOTAL:</strong>
            </div>
            <div v-if="index == 'Amount'">
              <strong>{{totals(tableItems)}} Rwf</strong>
            </div>
          </b-td>
        </b-tr>
      </template>
    </b-table>
    <b-pagination
      size="sm"
      align="center"
      v-model="pagination.currentPage"
      :total-rows="pagination.totalRows"
      :per-page="pagination.perPage"
      class="my-0"
      pills
      v-if="!loading.request && pagination.totalRows/pagination.perPage > 1"
    ></b-pagination>
    <add-property
      :show="addProperty.show"
      v-on:closeModal="addProperty.show =false"
      v-on:refresh="loadData()"
    />
    <b-modal id="updateModal" v-model="updateModal.show" hide-footer>
      <template v-slot:modal-title>Modify House</template>
      <update-house
        v-if="updateModal.show"
        :item="updateModal.item"
        :option="updateModal.option"
        v-on:closeModal="closeUpdateModal"
      />
    </b-modal>
    <vue-simple-context-menu
      :elementId="'rightmenu'"
      :options="rightMenu.options"
      :ref="'rightMenu'"
      @option-clicked="showUpdateModal"
    ></vue-simple-context-menu>
  </b-container>
</template>
<script>
import updateHouse from "../components/updateHouse.vue";
import addPropertyModal from "../components/modals/addPropertyModal.vue";
import jsPDF from "jspdf";
const { Village } = require("rwanda");
const { isPhoneNumber } = require("rwa-validator");
import "jspdf-autotable";
export default {
  name: "reports",
  components: {
    "update-house": updateHouse,
    "add-property": addPropertyModal
  },
  data() {
    return {
      addProperty: {
        show: false
      },
      selected: null,
      title: null,
      width: 0,
      options: [],
      color: "#333333bd",
      loading: {
        progress: false,
        request: false
      },
      rightMenu: {
        options: [{ name: "Edit House", slug: "edit" }]
      },
      updateModal: {
        show: false,
        item: [],
        option: []
      },
      modal: {
        show: false,
        switch: false,
        loading: false,
        title: "Search House Owner",
        btnContent: "Search",
        form: {
          fname: null,
          lname: null,
          phone: null,
          id: null,
          due: "500"
        },
        select: {
          cell: null,
          village: null
        }
      },
      search: {
        name: "",
        datalist: []
      },
      select: {
        sector: "Remera",
        cell: null,
        village: null,
        sectorOptions: [],
        cellOptions: [],
        villageOptions: [],
        shownColumn: [],
        selectAll: true
      },
      size: "5px",
      sortBy: "owner",
      sortDesc: false,
      fields: [
        { key: "index", label: "NO" },
        { key: "owner", label: "Full name", sortable: true },
        { key: "id", label: "House Code" },
        { key: "owner.phone", label: "Phone Number" },
        { key: "address.sector", label: "sector", sortable: true },
        { key: "address.cell", label: "Cell", sortable: true },
        { key: "address.village", label: "Village", sortable: true },
        { key: "occupied", label: "Rented", sortable: true },
        { key: "due", label: "Amount" }
      ],
      items: [],
      tableItems: [],
      pagination: {
        perPage: 15,
        currentPage: 1,
        totalRows: 1,
        show: false
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    sectorOptions() {
      return [this.activeSector];
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      const cell = this.select.cell;
      if (cell) {
        return Village("Kigali", "Gasabo", "Remera", cell);
      } else {
        return [];
      }
    },
    activeSector() {
      return this.capitalize(this.$store.getters.getActiveSector);
    },
    activeCell() {
      return this.capitalize(this.$store.getters.getActiveCell);
    },
    columns() {
      let array = [];
      this.fields.forEach(i => array.push(i.label));
      return array;
    },
    checkNumber() {
      const n = this.modal.form.phone;
      return n ? isPhoneNumber(n) : null;
    },
    user() {
      return this.$store.getters.userDetails;
    }
  },
  watch: {
    items() {
      handler: {
        this.tableItems = this.items;
      }
    },
    "select.shownColumn"() {
      handler: {
        this.select.selectAll =
          this.columns.length == this.select.shownColumn.length ? true : false;
      }
    },
    selected() {
      if (this.selected) {
        this.title = `List of users in ${this.selected}`;
      } else {
        this.title = `List of users in ${this.activeSector}`;
      }
    },
    "search.name"() {
      handler: {
        this.pagination.currentPage = 1;
        const searchedName = this.search.name;
        this.tableItems = this.filter().filter(obj => {
          const name = this.lc(obj.owner.fname + " " + obj.owner.lname);
          this.search.datalist = [...new Set([...this.search.datalist, name])];
          if (name.search(new RegExp(searchedName, "i")) != -1) {
            return (
              obj.owner.fname.includes(obj.owner.fname) ||
              obj.owner.lname.includes(obj.owner.lname)
            );
          }
        });
        while (this.search.datalist.length > 7) {
          this.search.datalist.pop();
        }
      }
    }
  },
  mounted() {
    this.loadData();
    this.select.shownColumn = this.columns;
    this.title = `List of users in ${this.activeSector}`;
  },
  methods: {
    loadData() {
      this.loading.request = true;
      const axios = this.axios;
      var promise;
      if (this.user.role.toLowerCase() == "basic") {
        promise =
          this.endpoint + `/properties?cell=${this.activeCell}&offset=0&limit=`;
      } else {
        promise =
          this.endpoint +
          `/properties?sector=${this.activeSector}&offset=0&limit=`;
      }
      axios
        .get(promise + `0`)
        .then(result => {
          axios
            .get(promise + `${result.data.Total}`)
            .then(res => {
              this.items = [...res.data.Properties];
              this.pagination.totalRows = res.data.Total;
            })
            .catch(err => {
              if (navigator.onLine) {
                const error = err.response
                  ? err.response.data.error || err.response.data
                  : "an error occured";
                this.$snotify.error(error);
              } else {
                this.$snotify.error("Please connect to the internet");
              }
            })
            .finally(() => {
              this.loading.request = false;
            });
        })
        .catch(err => {
          this.loading.request = false;
          if (navigator.onLine) {
            const error = err.response
              ? err.response.data.error || err.response.data
              : "an error occured";
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
        });
    },
    editHouse(house, index, evt) {
      evt.preventDefault();
      this.$refs.rightMenu.showMenu(evt, house);
    },
    showUpdateModal(data) {
      if (data) {
        this.updateModal.item = data.item ? data.item : {};
        this.updateModal.option = data.option ? data.option : {};
        this.updateModal.show = data.item ? true : false;
      }
    },
    closeUpdateModal() {
      this.loadData();
      this.updateModal.show = false;
    },
    totals(data) {
      if (data) {
        let total = 0;
        data.forEach(element => {
          total += Number(element.due);
        });
        return total.toLocaleString();
      }
    },
    filter() {
      this.$refs.dropdown.hide(true);
      //disabling some of the columns
      this.fields.forEach(value => {
        if (!this.select.shownColumn.includes(value.label)) {
          value.tdClass = "d-none";
          value.thClass = "d-none";
        } else {
          delete value.tdClass;
          delete value.thClass;
        }
      });
      const sector = this.activeSector;
      const cell = this.select.cell ? this.select.cell.toLowerCase() : "";
      const village = this.select.village
        ? this.select.village.toLowerCase()
        : "";
      this.tableItems = this.items;
      return this.tableItems.filter(item => {
        return (
          item.address.sector.search(new RegExp(sector, "i")) != -1 &&
          item.address.cell.search(new RegExp(cell, "i")) != -1 &&
          item.address.village.search(new RegExp(village, "i")) != -1
        );
      });
    },
    allSelected(checked) {
      this.select.shownColumn = checked ? this.columns.slice() : [];
    },
    clearFilter() {
      this.select.cell = null;
      this.select.village = null;
      this.tableItems = this.items;
      this.selected = null;
      this.$refs.dropdown.hide(true);
    },
    capitalize(string) {
      string.toLowerCase();
      return string.charAt(0).toUpperCase() + string.slice(1);
    },
    lc(a) {
      return a.toLowerCase();
    },
    confirm(message) {
      return this.$bvModal.msgBoxConfirm(message, {
        title: "Please Confirm",
        buttonSize: "sm",
        okVariant: "danger",
        okTitle: "YES",
        cancelTitle: "NO",
        footerClass: "p-3",
        hideHeaderClose: false,
        centered: true
      });
    }
  }
};
</script>
<style>
.table-container {
  position: relative;
  min-height: 100%;
}

hr {
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}

.table-container > h4.title {
  font-size: 20px;
  text-transform: capitalize;
}

.title .add-property {
  float: right;
  padding: 5px 10px;
  height: fit-content;
  border: none;
  margin-top: -5px;
}

.modal-body form button {
  float: right;
  margin-left: 10px;
  padding: 3px 15px;
}

.controllers {
  margin: 10px 0;
  height: 30px;
  display: flex;
}

.controllers .download {
  margin-left: 10px;
  outline: none;
  padding: 0.25rem 10px;
  height: fit-content;
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
  padding-bottom: 5px;
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
  margin-bottom: 5px;
}

.filter-dropdown hr {
  margin-top: 10px;
  margin-bottom: 10px;
}

.filter-dropdown button {
  float: right;
  height: fit-content;
  padding: 0.25rem 10px;
  font-size: 15px;
}

.filter-dropdown .dropdown-menu {
  min-width: 200px;
  margin: 0 2px 0;
}

.filter-dropdown .dropdown-menu > button {
  font-size: 13px !important;
  padding: 5px 20px !important;
  margin: 0 10px 0 0;
  width: fit-content;
}

.filter-dropdown .dropdown-menu form {
  outline: none !important;
  display: flex !important;
  width: 400px !important;
  padding: 5px !important;
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
  font-size: 14px;
}
</style>
