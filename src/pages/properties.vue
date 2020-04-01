<template>
  <b-container class="table-container px-5 py-3 mw-100">
    <vue-title title="Paypack | Properties" />
    <h4 class="title d-flex justify-content-between flex-nowrap">
      List of properties in {{selected}}
      <b-button class="add-property mb-1 font-14" variant="info" @click="addProperty.show = true">
        <i class="fa fa-plus-circle"></i> Property
      </b-button>
    </h4>
    <hr />

    <!-- filters -->
    <b-row class="my-1 align-items-end px-3 flex-nowrap">
      <b-dropdown
        id="dropdown-dropright"
        dropright
        variant="info"
        ref="dropdown"
        class="filter-dropdown mr-auto"
      >
        <template slot="button-content">
          <p class="font-14 d-inline">Filter By</p>
        </template>
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
              </template>
              <b-form-checkbox-group
                id="columns"
                v-model="select.postColumns"
                :options="columns"
                size="sm"
                name="columns"
                stacked
              ></b-form-checkbox-group>
            </b-form-group>
          </b-card-body>
        </b-dropdown-form>
        <b-button variant="primary" size="sm" @click.prevent="disableColumns">Ok</b-button>
        <b-button variant="danger" size="sm" @click.prevent="clearFilter">Clear</b-button>
      </b-dropdown>

      <div>
        <b-form-input
          placeholder="search user..."
          class="rounded-2 font-14"
          type="search"
          size="sm"
          v-model="search.name"
          list="search-user-id"
        ></b-form-input>
        <datalist id="search-user-id">
          <option v-for="name in search.datalist" :key="name">{{ name }}</option>
        </datalist>
      </div>
      <b-button size="sm" variant="info" class="ml-1" @click="loadData">
        <i class="fa fa-sync-alt" :class="{'fa-spin':loading.request}"></i>
      </b-button>
      <b-button
        @click="downloadList"
        :disabled="filteredItems.length ? false:true"
        variant="info"
        size="sm"
        class="ml-1 font-14"
      >Download</b-button>
    </b-row>
    <!-- end of filters -->

    <b-table
      id="data-table"
      bordered
      striped
      hover
      small
      responsive
      :key="key"
      :items="shownItems"
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
          class="text-right m-0 text-center"
        >{{data.item.occupied ? data.item.occupied?'Yes':'No' : "No"}}</b-card-text>
      </template>
      <template v-slot:cell(index)="data">
        <article class="text-center">{{data.index + 1}}</article>
      </template>
      <template v-slot:table-busy>
        <div class="text-center my-2 p-3">
          <b-spinner small class="align-middle" />&nbsp;
          <strong>Loading...</strong>
        </div>
      </template>
      <template v-slot:empty>
        <h5
          class="text-center font-14 my-4"
        >{{search.name ? search.name+' "is not in the list"':'No Property Found!'}}</h5>
      </template>
      <template v-slot:custom-foot v-if="!loading.request">
        <b-tr v-if="select.shownColumn.includes('Amount')">
          <b-td v-for="(index,i) in select.shownColumn" :key="i" variant="primary">
            <div
              class="text-danger text-right"
              v-if="select.shownColumn.length > 1 && i===select.shownColumn.length-2 "
            >
              <strong>Total:</strong>
            </div>
            <div v-if="i===select.shownColumn.length-1">
              <strong>{{totals(filteredItems)}} Rwf</strong>
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
    <vue-menu
      :elementId="'rightmenu'"
      :options="rightMenu.options"
      :ref="'rightMenu'"
      @option-clicked="showContextMenu"
    ></vue-menu>

    <message
      :phones="message.sendTo"
      v-if="message.show"
      v-on:modal-closed="message.show= false"
      v-on:sent="messageSent"
    />
  </b-container>
</template>
<script>
const download = import(
  /* webpackChunkName: "downloadScript" */ "../components/download scripts/downloadProperties"
);
import { Village } from "rwanda";
import { isPhoneNumber } from "rwa-validator";
export default {
  name: "reports",
  components: {
    "update-house": () =>
      import(
        /* webpackChunkName: "updateHouse" */ "../components/updateHouse.vue"
      ),
    "add-property": () =>
      import(
        /* webpackChunkName: "addProperty" */ "../components/modals/addPropertyModal.vue"
      ),
    message: () =>
      import(/* webpackChunkName: "message" */ "../components/modals/message")
  },
  data() {
    return {
      addProperty: {
        show: false
      },
      selected: null,
      width: 0,
      options: [],
      color: "#333333bd",
      loading: {
        progress: false,
        request: false
      },
      rightMenu: {
        options: [
          { name: "Edit House", slug: "edit" },
          { name: "Delete House", slug: "delete" },
          { name: "Send Message", slug: "send" }
        ]
      },
      updateModal: {
        show: false,
        item: [],
        option: []
      },
      message: {
        sendTo: [],
        show: false
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
        postColumns: []
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
      key: 1,
      pagination: {
        perPage: 15,
        currentPage: 1,
        totalRows: 1,
        show: false
      }
    };
  },
  computed: {
    sectorOptions() {
      return [this.activeSector];
    },
    filteredItems() {
      const sector = this.activeSector.trim().toLowerCase();
      const cell = this.select.cell || "";
      const village = this.select.village || "";
      this.selected = village || cell || sector;
      return this.items.filter(item => {
        if (sector && cell && village) {
          return (
            item.address.sector.toLowerCase().trim() == sector &&
            item.address.cell.toLowerCase().trim() == cell.toLowerCase() &&
            item.address.village.toLowerCase().trim() == village.toLowerCase()
          );
        }
        if (sector && cell && !village) {
          return (
            item.address.sector.toLowerCase().trim() == sector &&
            item.address.cell.toLowerCase().trim() == cell.toLowerCase()
          );
        }
        if (sector && !cell && !village) {
          return item.address.sector.toLowerCase().trim() == sector;
        }
      });
    },
    shownItems() {
      if (this.items) {
        this.pagination.currentPage = 1;
        this.pagination.totalRows = this.filteredItems.length;
        while (this.search.datalist.length > 7) {
          this.search.datalist.pop();
        }
        return this.filteredItems.filter(item => {
          return (item.owner.fname + " " + item.owner.lname)
            .toLowerCase()
            .includes(this.search.name.toLowerCase());
        });
      } else [];
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
      return this.fields.map(i => i.label);
    },
    checkNumber() {
      return this.modal.form.phone
        ? isPhoneNumber(this.modal.form.phone)
        : null;
    },
    user() {
      return this.$store.getters.userDetails;
    }
  },
  mounted() {
    this.loadData();
    this.select.postColumns = this.columns;
    if (this.user.role.toLowerCase() == "basic") {
      this.selected = this.activeCell;
    } else {
      this.selected = this.activeSector;
    }
  },
  methods: {
    async loadData() {
      this.loading.request = true;
      var promise = await this.getUrl();
      var total = await this.getTotal();
      // total = (total / 10).toFixed();
      this.axios
        .get(promise + `${total}`)
        .then(res => {
          this.items = [...res.data.Properties];
          this.pagination.totalRows = total;
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.loading.request = false;
        });
    },
    async getTotal() {
      var promise = await this.getUrl();
      return this.axios
        .get(`${promise}0`)
        .then(res => res.data.Total)
        .catch(err => null);
    },
    downloadList() {
      download(this.filteredItems, this.selected);
    },
    editHouse(house, index, evt) {
      evt.preventDefault();
      this.$refs.rightMenu.showMenu(evt, house);
    },
    showContextMenu(data) {
      if (data.option.slug == "edit") {
        this.updateModal.item = data.item ? data.item : {};
        this.updateModal.option = data.option ? data.option : {};
        this.updateModal.show = data.item ? true : false;
      } else if (data.option.slug == "delete") {
        this.deleteHouse(data.item);
      } else if (data.option.slug == "send") {
        this.message.sendTo.push(data.item.owner.phone);
        this.message.show = true;
      }
    },
    deleteHouse(house) {
      const message = `Do you want to delete this house? Names: ${house.owner.fname} ${house.owner.lname}  ID: ${house.id}`;

      this.$bvModal
        .msgBoxConfirm(message, {
          title: "Delete Property?",
          size: "sm",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "Delete",
          cancelTitle: "NO",
          footerClass: "p-2",
          hideHeaderClose: false,
          centered: true
        })
        .then(value => {
          if (value === true) {
            this.loading.request = true;
            this.axios
              .delete("/properties/" + house.id)
              .then(res => {
                console.log(res.data);
                this.$snotify.info("House deleted successfully");
                this.loadData();
              })
              .catch(err => {
                const error = err.response
                  ? err.response.data.error || err.response.data
                  : null;
                this.loading.request = false;
                if (error) this.$snotify.error(error);
              });
          }
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        });
    },
    closeUpdateModal() {
      this.loadData();
      this.updateModal.show = false;
    },
    messageSent() {
      this.message.show = false;
      this.message.sendTo = [];
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
    disableColumns() {
      //disabling some of the columns
      this.loading.request = true;
      if (this.$refs.dropdown) this.$refs.dropdown.hide(true);
      this.select.shownColumn = this.select.postColumns;
      this.fields.map(value => {
        if (!this.select.shownColumn.includes(value.label)) {
          value.tdClass = "d-none";
          value.thClass = "d-none";
        } else {
          delete value.tdClass;
          delete value.thClass;
        }
      });
      this.key++;
      this.loading.request = false;
    },
    clearFilter() {
      this.select.cell = null;
      this.select.village = null;
      this.select.shownColumn = this.columns;
      this.key = 1;
      if (this.user.role.toLowerCase() == "basic") {
        this.selected = this.activeCell;
      } else {
        this.selected = this.activeSector;
      }
      this.$refs.dropdown.hide(true);
    },
    capitalize(string) {
      string.toLowerCase();
      return string.charAt(0).toUpperCase() + string.slice(1);
    },
    lc(a) {
      return a.toLowerCase();
    },
    getUrl() {
      if (this.user.role.toLowerCase() == "basic") {
        return `/properties?cell=${this.activeCell}&offset=0&limit=`;
      } else {
        return `/properties?sector=${this.activeSector}&offset=0&limit=`;
      }
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
<style lang="scss">
@import "../assets/css/properties.scss";
</style>
