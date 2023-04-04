<template>
  <b-container class="table-container px-5 py-3" fluid>
    <vue-title title="Paypack | Properties" />
    <header>
      <h4 class="primary-font">List of properties in {{ selected }}</h4>
      <b-button
        class="br-2 primary-font"
        variant="info"
        @click="addProperty.show = true"
      >
        <i class="fa fa-plus-circle mr-1" />Property
      </b-button>
    </header>
    <hr class="my-1" />

    <!-- filters -->
    <b-row class="mt-3 mb-2 m-0 flex-nowrap">
      <b-dropdown
        id="dropdown-dropright"
        variant="info"
        ref="dropdown"
        toggle-class="br-2 primary-font"
        class="filter-dropdown mr-auto primary-font"
        text="Filter By"
      >
        <b-dropdown-form>
          <!-- location Filters -->
          <b-form-group label="Location Filters" class="flex-grow-1 px-2">
            <b-form-group
              label="Sector"
              label-class="text-muted p-0"
              class="mb-3"
            >
              <b-form-select
                size="sm"
                class="br-2"
                v-model="select.sector"
                :options="sectorOptions"
                :disabled="sectorOptions.length < 1 || isManager"
              >
                <template v-slot:first>
                  <option :value="null">All</option>
                </template>
              </b-form-select>
            </b-form-group>
            <b-form-group
              label="Cell"
              label-class="text-muted p-0"
              class="mb-3"
            >
              <b-form-select
                size="sm"
                class="br-2"
                v-model="select.cell"
                :options="cellOptions"
                :disabled="cellOptions.length < 1 || isManager"
              >
                <template v-slot:first>
                  <option :value="null">All</option>
                </template>
              </b-form-select>
            </b-form-group>
            <b-form-group label="Village" label-class="text-muted p-0">
              <b-form-select
                size="sm"
                class="br-2"
                v-model="select.village"
                :options="villageOptions"
                :disabled="villageOptions.length < 1"
              >
                <template v-slot:first>
                  <option :value="null">All</option>
                </template>
              </b-form-select>
            </b-form-group>
          </b-form-group>
          <b-form-group label="Columns Filters" class="flex-grow-1 px-2">
            <b-form-checkbox-group
              id="columns"
              v-model="select.postColumns"
              :options="columns"
              size="sm"
              name="columns"
              stacked
            ></b-form-checkbox-group>
          </b-form-group>
        </b-dropdown-form>
        <b-button variant="info" class="br-2" @click.prevent="filterByLocation"
          >Filter</b-button
        >
        <b-button variant="danger" class="br-2" @click.prevent="clearFilter"
          >Reset</b-button
        >
      </b-dropdown>

      <div class="ml-2">
        <b-form-input
          placeholder="search user..."
          class="br-2 primary-font"
          type="search"
          v-model="search.name"
          @keypress.enter="searchProperties"
        ></b-form-input>
      </div>
      <b-button variant="info" class="ml-2 br-2" @click="loadData">
        <i class="fa fa-sync-alt" :class="{ 'fa-spin': loading.request }"></i>
      </b-button>
      <b-button
        @click="downloadList"
        :disabled="filteredData.length ? false : true"
        variant="info"
        class="ml-2 br-2 primary-font"
      >
        Download
        <i class="fa fa-download" />
      </b-button>
    </b-row>
    <!-- end of filters -->

    <b-table
      hover
      small
      bordered
      responsive
      id="data-table"
      :fields="fields"
      :items="filteredData"
      head-variant="light"
      :sort-by.sync="sortBy"
      :busy="loading.request"
      @row-contextmenu="editHouse"
      thead-class="primary-font"
      tbody-class="secondary-font"
      :show-empty="!loading.request"
    >
      <template v-slot:cell(due)="data"
        >{{ data.item.due | number }} Rwf</template
      >
      <template v-slot:cell(owner)="data">{{
        data.item.owner.fname + " " + data.item.owner.lname
      }}</template>
      <template v-slot:cell(occupied)="data">
        <b-card-text class="text-right m-0 text-center">{{
          data.item.occupied ? (data.item.occupied ? "Yes" : "No") : "No"
        }}</b-card-text>
      </template>
      <template v-slot:cell(index)="data">
        <article class="text-center">{{ data.index + 1 }}</article>
      </template>
      <template v-slot:table-busy>
        <vue-load class="primary-font" />
      </template>
      <template v-slot:empty>
        <h6 class="text-center font-weight-bold p-5 w-100 primary-font">
          {{
            search.name
              ? search.name + ' "is not available in this property list"'
              : "There are no properties available to show at the moment!"
          }}
        </h6>
      </template>
      <template v-slot:custom-foot v-if="!loading.request">
        <b-tr v-if="select.shownColumn.includes('Amount')">
          <b-td
            v-for="(index, i) in select.shownColumn"
            :key="i"
            variant="primary"
          >
            <div
              class="text-danger text-right"
              v-if="
                select.shownColumn.length > 1 &&
                i === select.shownColumn.length - 2
              "
            >
              <strong>Total:</strong>
            </div>
            <div v-if="i === select.shownColumn.length - 1">
              <strong>{{ totalAmount | number }} Rwf</strong>
            </div>
          </b-td>
        </b-tr>
      </template>
    </b-table>

    <b-pagination
      class="my-0"
      align="center"
      v-if="showPagination"
      :per-page="pagination.perPage"
      v-model="pagination.currentPage"
      :total-rows="pagination.totalRows"
      @input="loadData"
    ></b-pagination>
    <add-property
      :show="addProperty.show"
      @closeModal="addProperty.show = false"
      @added="propertyAdded"
    />
    <b-modal
      id="updateModal"
      v-model="updateModal.show"
      hide-footer
      body-class="px-4 py-3"
    >
      <template v-slot:modal-title>
        <header class="primary-font">Modify House</header>
      </template>
      <update-house
        v-if="updateModal.show"
        :item="updateModal.item"
        :option="updateModal.option"
        @closeModal="closeUpdateModal"
        @updated="state.reloadData = true"
      />
    </b-modal>
    <vue-menu
      :elementId="'rightmenu'"
      :options="rightMenu.options"
      :ref="'rightMenu'"
      @option-clicked="showContextMenu"
    />

    <message
      :phones="message.sendTo"
      v-if="message.show"
      @modal-closed="closeModal('message')"
      @sent="messageSent"
    />
  </b-container>
</template>
<script>
import download from "../components/download scripts/downloadProperties";
export default {
  name: "properties",
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
      import(/* webpackChunkName: "message" */ "../components/modals/message"),
  },
  data() {
    return {
      development: true,
      originalData: [],
      filteredData: [],
      totalAmount: 0,
      addProperty: { show: false },
      state: {
        changedLocation: false,
        reloadData: false,
      },
      selected: null,
      width: 0,
      options: [],
      loading: {
        progress: false,
        request: false,
      },
      rightMenu: {
        options: [
          { name: "Edit House", slug: "edit" },
          { name: "Delete House", slug: "delete" },
          { name: "Send Message", slug: "send" },
        ],
      },
      updateModal: {
        show: false,
        item: [],
        option: [],
      },
      message: {
        sendTo: [],
        show: false,
      },
      search: { name: "" },
      select: {
        sector: null,
        cell: null,
        village: null,
        shownColumn: [],
        postColumns: [],
      },
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
        { key: "due", label: "Amount" },
      ],
      pagination: {
        perPage: 20,
        currentPage: 1,
        totalRows: 1,
      },
    };
  },
  computed: {
    showPagination() {
      if (this.loading.request) return false;
      if (this.pagination.totalRows < this.pagination.perPage) return false;
      return true;
    },
    sectorOptions() {
      return [this.activeSector] || [];
    },
    cellOptions() {
      const { province, district, sector } = this.location;
      if (this.select.sector)
        return this.$cells(province, district, this.select.sector);
      return [];
    },
    villageOptions() {
      const { province, district } = this.location;
      if (this.select.cell)
        return this.$villages(
          province,
          district,
          this.select.sector,
          this.select.cell
        );
      return [];
    },
    activeSector() {
      return this.$capitalize(this.$store.getters.getActiveSector);
    },
    activeCell() {
      return this.$capitalize(this.$store.getters.getActiveCell);
    },
    columns() {
      return this.fields.map((i) => i.label);
    },
    totals() {
      if (this.filteredData.length > 0)
        return this.filteredData.reduce(
          (a, b) => Number(a.due || a) + Number(b.due)
        );
      return 0;
    },
    user() {
      return this.$store.getters.userDetails;
    },
    location() {
      return this.$store.getters.location;
    },
    role() {
      return this.user.role;
    },
    isManager() {
      return this.user.role === "basic";
    },
  },
  mounted() {
    this.loadData();
    this.select.postColumns = this.columns;
    this.select.shownColumn = this.columns;
    if (this.isManager) {
      this.selected = this.activeCell;
      this.select.sector = this.activeSector;
      this.select.cell = this.activeCell;
    } else {
      this.selected = this.activeSector;
    }
  },
  watch: {
    "search.name"() {
      handler: {
        this.filteredData = this.filterByName(this.search.name);
      }
    },
    "select.sector"() {
      handler: {
        this.state.changedLocation = true;
        if (this.isManager) this.select.cell = this.activeCell;
        else this.select.cell = null;
      }
    },
    "select.cell"() {
      handler: {
        this.state.changedLocation = true;
        this.select.village = null;
      }
    },
    "select.village"() {
      handler: {
        this.state.changedLocation = true;
      }
    },
  },
  methods: {
    pageChanged() {
      this.loadData();
    },
    async loadData() {
      this.loading.request = true;
      var promise = this.getUrl();

      this.axios
        .get(promise)
        .then((res) => {
          this.filteredData = res.data.Properties;
          this.originalData = res.data.Properties;
          this.totalAmount = res.data.amount;
          this.pagination.totalRows = res.data.Total;
          this.filterByLocation();
        })
        .catch((err) => {
          console.log(err);
          try {
            this.$snotify.error(err.response.data.error || err.response.data);
          } catch {
            this.$snotify.error(
              "Failed to load data from database! check your connectivity and try again"
            );
          }
        })
        .finally(() => {
          this.loading.request = false;
          this.state.reloadData = false;
        });
    },
    propertyAdded() {
      this.state.reloadData = true;
      this.loadData();
    },
    filterByName(name) {
      if (!name) return this.originalData;
      return this.originalData.filter((item) =>
        `${item.owner.fname} ${item.owner.lname}`
          .toLowerCase()
          .includes(name.toLowerCase())
      );
    },
    async filterByLocation() {
      if (this.$refs.dropdown) await this.$refs.dropdown.hide(true);
      this.disableColumns();
      const { sector, cell, village } = this.select;
      if (!sector && !cell && !village) return this.originalData;
      this.selected =
        village || cell || sector || this.activeSector || this.activeCell;

      this.filteredData = this.originalData.filter(
        (item) =>
          item.address.sector
            .toLowerCase()
            .endsWith(String(sector || "").toLowerCase()) &&
          item.address.cell
            .toLowerCase()
            .endsWith(String(cell || "").toLowerCase()) &&
          item.address.village
            .toLowerCase()
            .endsWith(String(village || "").toLowerCase())
      );
    },
    disableColumns() {
      //disabling unsellected columns by applying the display class on both thead and td
      if (this.select.postColumns.length === this.fields.length) return;
      this.select.shownColumn = this.select.postColumns;
      this.fields.map((value) => {
        if (!this.select.shownColumn.includes(value.label)) {
          value.tdClass = "d-none";
          value.thClass = "d-none";
        } else {
          delete value.tdClass;
          delete value.thClass;
        }
      });
    },
    downloadList() {
      if (this.filteredData.length > 0) {
        var today = new Date().toLocaleDateString("en-EN", {
          year: "numeric",
          month: "long",
        });
        var Headers = [
          "ID",
          "Names",
          "House Code",
          "Phone Number",
          "Sector",
          "Cell",
          "Village",
          "Payment Amount",
        ];
        var Body = this.filteredData.map((item, i) => {
          var result = Headers.map((i) => "");
          result[Headers.indexOf("ID")] = i;
          result[
            Headers.indexOf("Names")
          ] = `${item.owner.fname} ${item.owner.lname}`;
          result[Headers.indexOf("House Code")] = item.id;
          result[Headers.indexOf("Phone Number")] = item.owner.phone;
          result[Headers.indexOf("Sector")] = item.address.sector;
          result[Headers.indexOf("Cell")] = item.address.cell;
          result[Headers.indexOf("Village")] = item.address.village;
          result[Headers.indexOf("Payment Amount")] = `${Number(
            item.due
          ).toLocaleString()} Rwf`;
          return result;
        });
        var data = {
          title: String(`List of Properties in ${this.selected}`).toUpperCase(),
          name: `List of Properties in ${this.selected} in ${today}`,
          data: {
            Headers: Headers,
            Body: Body,
          },
        };
        download(data);
      } else this.$snotify.error("There are no Data available to download!");
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
        this.message.sendTo = [data.item.owner.phone];
        this.message.show = true;
      }
    },
    deleteHouse(house) {
      const el = this.$createElement;
      const messageNode = el(
        "ul",
        { class: "list-style-none p-0 text-center" },
        [
          el("li", { class: "mb-2" }, [
            `Names: ${house.owner.fname} ${house.owner.lname}`,
          ]),
          el("li", { class: "mb-2" }, [`HouseId: ${house.id}`]),
        ]
      );
      const message = `Names: ${house.owner.fname} ${house.owner.lname}  ID: ${house.id}`;

      this.$bvModal
        .msgBoxConfirm(messageNode, {
          title: "Are you sure to delete this house?",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "yes! delete",
          cancelTitle: "NO",
          footerClass: "p-2",
          contentClass: "primary-font",
          hideHeaderClose: false,
        })
        .then((value) => {
          if (value === true) {
            this.loading.request = true;
            this.axios
              .delete("/properties/" + house.id)
              .then((res) => {
                this.state.reloadData = true;
                this.$snotify.info("House deleted successfully");
                this.loadData();
              })
              .catch((err) => {
                const error = err.response
                  ? err.response.data.error || err.response.data
                  : null;
                this.loading.request = false;
                try {
                  this.$snotify.error(error);
                } catch {
                  this.$snotify.error("Failed to delete House");
                }
              });
          }
        })
        .catch((err) => {
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
    closeModal(name) {
      if (name == "message") {
        this.message.sendTo = [];
        this.message.show = false;
      }
    },
    messageSent() {
      this.message.show = false;
      this.message.sendTo = [];
    },
    async clearFilter() {
      await this.$refs.dropdown.hide(true);
      this.select.village = null;
      this.select.shownColumn = this.columns;
      this.select.postColumns = this.columns;
      if (this.isManager) {
        this.select.sector = this.activeSector;
        this.select.cell = this.activeCell;
      } else {
        this.select.sector = null;
        this.select.cell = null;
      }
      this.filterByLocation();
    },
    getUrl() {
      const offset =
        (this.pagination.currentPage - 1) * this.pagination.perPage;
      const limit = this.pagination.perPage;
      if (this.isManager)
        return `/properties?cell=${this.activeCell}&offset=${offset}&limit=${limit}`;
      return `/properties?sector=${this.activeSector}&offset=${offset}&limit=${limit}`;
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
        centered: true,
      });
    },
    searchProperties() {
      console.log("searching");
    },
  },
};
</script>
<style lang="scss">
@import "../assets/css/properties.scss";
</style>