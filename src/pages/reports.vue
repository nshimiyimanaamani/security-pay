<template>
  <div class="table-container">
    <vue-title title="Paypack | Properties" />
    <h4 class="title text-center">
      {{title}}
      <b-button class="add-property mb-1" variant="info" @click="modal.show = ! modal.show">
        <i class="fas fa-plus-circle"></i> Property
      </b-button>
    </h4>
    <hr />
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
          <b-card-body class="p-2">
            <b-form-group label="sector">
              <b-form-select v-model="select.sector" :options="sectorOptions">
                <template v-slot:first>
                  <option :value="null" disabled>Select sector</option>
                </template>
              </b-form-select>
            </b-form-group>
            <b-form-group label="cell" v-show="select.sector">
              <b-form-select v-model="select.cell" :options="cellOptions">
                <template v-slot:first>
                  <option :value="null" disabled>Select cell</option>
                </template>
              </b-form-select>
            </b-form-group>
            <b-form-group label="village" v-show="select.sector && select.cell">
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
      <div class="search">
        <b-form-input
          placeholder="search user..."
          size="sm"
          v-model="search.name"
          list="search-datalist-id"
        ></b-form-input>
        <b-button variant="info" style="height: 100%" @click="search.name = ''">
          <i class="fa fa-times"></i>
        </b-button>
        <datalist id="search-datalist-id">
          <option v-for="name in search.datalist" :key="name">{{ name }}</option>
        </datalist>
      </div>
      <b-button @click.prevent="download" class="download btn-info">Download</b-button>
    </div>
    <b-table
      id="data-table"
      bordered
      striped
      hover
      small
      :items="tableItems"
      :fields="fields"
      :busy="loading.request"
      :sort-by.sync="sortBy"
      :sort-desc.sync="sortDesc"
      show-empty
      :current-page="pagination.currentPage"
      :per-page="pagination.perPage"
    >
      <template v-slot:cell(due)="data">{{Number(data.item.due).toLocaleString()}} Rwf</template>
      <template v-slot:cell(owner)="data">{{data.item.owner.fname +" "+ data.item.owner.lname}}</template>
      <template v-slot:cell(index)="data">
        <article class="text-center">{{data.index + 1}}</article>
      </template>
      <template v-slot:table-busy>
        <div class="text-center my-2">
          <b-spinner class="align-middle"></b-spinner>
          <strong>Loading...</strong>
        </div>
      </template>
      <template v-slot:empty="scope">
        <h5
          class="text-center my-4"
        >{{search.name ? search.name+' "is not availble in the list"':'No user Found!'}}</h5>
      </template>
      <template v-slot:custom-foot="items" v-if="!loading.request">
        <b-tr v-if="select.shownColumn.includes('Amount')">
          <b-td v-for="index in select.shownColumn" :key="index" variant="primary">
            <div
              v-if="index == select.shownColumn[select.shownColumn.indexOf('Amount')-1]"
              class="text-danger"
            >
              <strong>TOTAL:</strong>
            </div>
            <div v-if="index == 'Amount'">
              <strong>{{totals(items.items)}} Rwf</strong>
            </div>
          </b-td>
        </b-tr>
      </template>
    </b-table>
    <b-pagination
      size="sm"
      v-model="pagination.currentPage"
      :total-rows="pagination.totalRows"
      :per-page="pagination.perPage"
      align="fill"
      class="my-0"
      v-if="!loading.request"
    />
    <div class="add-property-modal" v-show="modal.show">
      <!-- Modal content -->
      <b-card class="mb-2 modal-body">
        <h5 class="text-center mb-1">{{modal.title}}</h5>
        <hr />
        <b-form @submit.prevent="search_user" @reset="resetModal">
          <b-form-group
            id="input-group-1"
            class="mb-2"
            label="First Name:"
            label-for="input-1"
            description="Amazina ya nyiri inzu (*Ntabwo ari ukodesheje)"
          >
            <b-form-input
              id="input-1"
              v-model="modal.form.fname"
              required
              placeholder="First name"
              :disabled="modal.switch"
            ></b-form-input>
          </b-form-group>
          <b-form-group id="input-group-2" class="mb-2" label="Last Name:" label-for="input-2">
            <b-form-input
              id="input-2"
              v-model="modal.form.lname"
              :disabled="modal.switch"
              required
              placeholder="Last name"
            ></b-form-input>
          </b-form-group>
          <b-form-group id="input-group-3" class="mb-2" label="Phone Number:" label-for="input-3">
            <b-form-input
              id="input-3"
              v-model="modal.form.phone"
              :state="checkNumber"
              :disabled="modal.switch"
              required
              type="number"
              placeholder="Phone number"
            ></b-form-input>
            <b-form-invalid-feedback
              :state="checkNumber"
            >Phone number must be greater than or equal 10.</b-form-invalid-feedback>
          </b-form-group>
          <b-form-group
            id="input-group-4"
            :label="'Due: '+ modal.form.due +' Rwf'"
            label-for="range-1"
            v-show="modal.switch"
            class="mb-2"
          >
            <b-form-input
              id="range-1"
              v-model="modal.form.due"
              type="range"
              min="0"
              max="10000"
              step="500"
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="input-group-5"
            class="mb-2"
            label="Cell:"
            label-for="input-4"
            v-show="modal.switch"
          >
            <b-form-select v-model="select.cell" :options="cellOptions" class="mb-0">
              <template v-slot:first>
                <option :value="null" disabled>select a cell</option>
              </template>
            </b-form-select>
          </b-form-group>
          <b-form-group
            id="input-group-6"
            label="Village:"
            label-for="input-5"
            v-show="modal.switch"
            class="mb-3"
          >
            <b-form-select v-model="select.village" :options="villageOptions" class="mb-0">
              <template v-slot:first>
                <option :value="null" disabled>select a village</option>
              </template>
            </b-form-select>
          </b-form-group>
          <b-button type="submit" variant="primary">
            {{modal.loading ? modal.btnContent+'ing' : modal.btnContent}}
            <b-spinner v-show="modal.loading" small type="grow"></b-spinner>
          </b-button>
          <b-button type="reset" variant="danger">cancel</b-button>
        </b-form>
      </b-card>
    </div>
  </div>
</template>
<script>
import jsPDF from "jspdf";
const { District, Sector, Cell, Village } = require("rwanda");
import "jspdf-autotable";
export default {
  name: "reports",
  data() {
    return {
      selected: null,
      title: null,
      width: 0,
      options: [],
      color: "#333333bd",
      loading: {
        progress: false,
        request: false
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
        sector: null,
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
        { key: "id", label: "House Code", sortable: false },
        { key: "owner.phone", label: "Phone Number", sortable: false },
        { key: "address.sector", label: "sector", sortable: true },
        { key: "address.cell", label: "Cell", sortable: true },
        { key: "address.village", label: "Village", sortable: true },
        { key: "due", label: "Amount", sortable: false }
      ],
      items: [],
      tableItems: [],
      pagination: {
        perPage: 15,
        currentPage: 1,
        totalRows: 1
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
      const sector = this.select.sector;
      if (sector) {
        return Cell("Kigali", "Gasabo", sector);
      } else {
        return Cell("Kigali", "Gasabo", "Remera");
      }
    },
    villageOptions() {
      const sector = this.select.sector;
      const cell = this.select.cell;
      if (sector && cell) {
        return Village("Kigali", "Gasabo", sector, cell);
      } else {
        if (cell) {
          return Village("Kigali", "Gasabo", "Remera", cell);
        } else {
          return [];
        }
      }
    },
    activeSector() {
      return this.capitalize(this.$store.getters.getActiveSector);
    },
    columns() {
      let array = [];
      this.fields.forEach(i => array.push(i.label));
      return array;
    },
    checkNumber() {
      if (this.modal.form.phone) {
        if (this.modal.form.phone.length >= 1) {
          return this.modal.form.phone.length >= 10;
        } else {
          return null;
        }
      }
    }
  },
  watch: {
    items() {
      this.tableItems = this.items;
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
    },
    tableItems() {
      handler: {
        this.pagination.totalRows = this.tableItems.length;
      }
    }
  },
  mounted() {
    this.loadData();
    this.select.shownColumn = this.columns;
    this.title = `List of users in ${this.activeSector}`;
    console.log(sessionStorage.email.toString());
  },
  methods: {
    loadData() {
      this.loading.request = true;
      this.axios
        .get(this.endpoint + "/properties?sector=remera&offset=0&limit=10")
        .then(res => {
          this.items = new Array();
          this.items = res.data.Properties;
          this.loading.request = false;
        })
        .catch(err => {
          console.log(err);
        });
    },
    search_user() {
      if (!this.modal.switch) {
        const fname = this.capitalize(this.modal.form.fname.trim());
        const lname = this.capitalize(this.modal.form.lname.trim());
        const phone = this.modal.form.phone.trim();
        this.modal.loading = true;
        this.axios
          .get(
            this.endpoint +
              "/owners/search?fname=" +
              fname +
              "&lname=" +
              lname +
              "&phone=" +
              phone
          )
          .then(res => {
            this.modal.loading = false;
            this.modal.switch = true;
            this.modal.title = "Register Property";
            this.modal.btnContent = "Register";
            this.modal.form.id = res.data.id;
          })
          .catch(err => {
            this.modal.loading = false;
            const message =
              fname +
              " " +
              lname +
              " is not registered! Do you want to register this user?";
            this.confirm(message).then(state => {
              if (state === true) {
                this.modal.loading = true;
                this.axios
                  .post(`${this.endpoint}/owners`, {
                    fname: fname,
                    lname: lname,
                    phone: phone
                  })
                  .then(res => {
                    this.modal.loading = false;
                    this.modal.switch = true;
                    this.modal.title = "Register Property";
                    this.modal.btnContent = "Register";
                    this.modal.form.id = res.data.id;
                    this.$snotify.info(
                      `User created. proceeding to registration...`
                    );
                  });
              }
            });
          });
      } else if (this.modal.switch) {
        this.modal.loading = true;
        console.log(this.modal, this.select);
        this.axios
          .post(this.endpoint + "/properties", {
            owner: {
              id: this.modal.form.id
            },
            address: {
              cell: this.select.cell,
              village: this.select.village,
              sector: "remera"
            },
            due: this.modal.form.due.toString(),
            recorded_by: "6de63fec-5f4a-4867-ae4c-3f3af70c9166"
          })
          .then(res => {
            this.resetModal();
            this.loadData();
            this.$snotify.info(`Property Registered successfully!`);
          })
          .catch(err => {
            this.modal.loading = false;
            this.$snotify.error(`Property Registration Failed!`);
          });
      }
    },
    resetModal() {
      this.modal.show = false;
      this.modal.switch = false;
      this.modal.loading = false;
      this.modal.title = "Search House Owner";
      this.modal.btnContent = "Search";
      this.modal.form.fname = null;
      this.modal.form.lname = null;
      this.modal.form.phone = null;
      this.modal.form.id = null;
      this.modal.form.due = null;
      this.modal.select.cell = null;
      this.modal.select.village = null;
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
      return this.items.filter(item => {
        return (
          item.address.sector
            .toLowerCase()
            .includes(this.select.sector.toLowerCase()) &&
          item.address.cell
            .toLowerCase()
            .includes(this.select.cell.toLowerCase()) &&
          item.address.village
            .toLowerCase()
            .includes(this.select.village.toLowerCase())
        );
      });
    },
    allSelected(checked) {
      this.select.shownColumn = checked ? this.columns.slice() : [];
    },
    clearFilter() {
      this.select.sector = null;
      this.select.cell = null;
      this.select.village = null;
      this.tableItems = this.items;
      this.selected = null;
      this.$refs.dropdown.hide(true);
    },
    download() {
      if (this.loading.request && !this.tableItems.length) {
        this.$snotify.error(
          "No Data available to download! refresh page to retry"
        );
      } else if (!this.loading.request && this.tableItems.length) {
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
    capitalize(string) {
      string.toLowerCase();
      return string.charAt(0).toUpperCase() + string.slice(1);
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
  padding: 15px 40px 5px;
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

.add-property-modal {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 100%;
  margin: auto;
  background: #000000cc;
  z-index: 100;
}

.add-property-modal .modal-body {
  position: sticky;
  -ms-flex: 1 1 auto;
  -webkit-box-flex: 1;
  flex: 1 1 auto;
  padding: 0;
  width: 40%;
  top: 5rem;
  margin: auto;
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
  padding: 2px 10px;
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
  padding: 2px 10px;
  font-size: 16px;
}

.filter-dropdown .dropdown-menu {
  min-width: 200px;
  margin: 0 2px 0;
}

.dropdown-menu > button {
  font-size: 13px !important;
  padding: 5px 20px !important;
  margin: 0 10px 0 0;
  width: fit-content;
}

.dropdown-menu form {
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
  text-transform: capitalize;
  font-size: 14px;
}
</style>
