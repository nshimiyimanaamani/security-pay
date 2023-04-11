<template>
  <div class="agent-payment-view">
    <b-row class="justify-content-between px-3">
      <b-button variant="info" @click="loadData">Refresh</b-button>
      <div class="button-group">
        <b-button variant="info" @click="state.search = true">
          search
          <i class="fa fa-search" />
        </b-button>
        <b-button variant="info" class="settings p-0">
          <selector :object="date" :id="'agent-date-selector'" @ok="loadData" />
        </b-button>
      </div>
    </b-row>
    <transition name="fade" :duration="500">
      <b-row
        class="flex-row justify-content-end w-100 m-0 flex-nowrap align-items-center"
        v-if="state.search"
      >
        <input
          type="search"
          name="search"
          id="agent-search-user"
          v-model="searchItem"
          class="search w-100 my-2 py-1 px-3"
        />
        <i
          class="fa fa-times ml-2 p-2 app-color text-white"
          @click="closeSearch"
        />
      </b-row>
    </transition>
    <b-row class="m-0 mt-2">
      <b-table
        small
        striped
        bordered
        hover
        responsive
        show-empty
        :busy="state.loading"
        :items="shownItems"
        :fields="table.fields"
        :sort-by.sync="table.sortBy"
      >
        <template v-slot:cell(index)="data">
          <article class="text-center">{{ data.index + 1 }}</article>
        </template>
        <template v-slot:empty>
          <article class="text-center p-3">No payment records found.</article>
        </template>
        <template v-slot:cell(amount)="data"
          >{{ data.item.amount | number }} Rwf</template
        >
        <template v-slot:cell(month)="data">{{
          months[new Date(data.item.date_recorded).getMonth()]
        }}</template>
        <template v-slot:cell(year)="data">{{
          new Date(data.item.date_recorded).getFullYear()
        }}</template>
        <template v-slot:cell(owner)="data">
          {{ data.item.owner_firstname + " " + data.item.owner_lastname }}
        </template>
        <template v-slot:cell(method)="data">
          <article :class="data.value">
            {{ data.value.includes("mtn") ? "mtn" : data.value }}
          </article>
        </template>
        <template v-slot:custom-foot>
          <b-tr class="total">
            <b-td></b-td>
            <b-td></b-td>
            <b-td></b-td>
            <b-td></b-td>
            <b-td></b-td>
            <b-td></b-td>
            <b-td></b-td>
            <b-td></b-td>
            <b-td></b-td>

            <b-td class="py-2">
              <small
                ><strong style=""
                  ><span style="color: #dc3545">Total </span>:
                  {{ totalAmount | number }} Rwf</strong
                ></small
              >
            </b-td>
          </b-tr>
        </template>
        <template v-slot:table-busy>
          <div class="text-center my-2">
            <loader />
          </div>
        </template>
      </b-table>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "agentUnpaidView",
  components: {
    loader: () => import("../loader"),
    selector: () => import("../agent-yearSelector"),
  },
  props: {
    user: Object,
  },
  data() {
    return {
      searchItem: "",
      state: {
        search: false,
        loading: false,
      },
      totalAmount: 0,
      date: {
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1,
      },
      table: {
        items: null,
        fields: [
          { key: "index", label: "NO" },
          { key: "owner", label: "Name", sortable: true },
          { key: "madefor", label: "House Code", sortable: false },
          { key: "year", label: "Year", sortable: true },
          { key: "month", label: "Month", sortable: true },
          { key: "sector", label: "Sector", sortable: true },
          { key: "cell", label: "Cell", sortable: true },
          { key: "village", label: "Village", sortable: true },
          { key: "method", label: "Method", sortable: true },
          { key: "amount", label: "Amount", sortable: false },
        ],
        sortBy: "owner",
      },
      pagination: {
        perPage: 15,
        total: 0,
        page: 1,
      },
    };
  },
  computed: {
    months() {
      return this.$store.getters.getMonths;
    },
    shownItems() {
      if (this.table.items) {
        return this.table.items.filter((item) => {
          return (item.owner_firstname + " " + item.owner_lastname)
            .toLowerCase()
            .includes(this.searchItem.toLowerCase());
        });
      } else [];
    },
    offset() {
      return (this.pagination.page - 1) * this.pagination.perPage;
    },
  },
  mounted() {
    this.loadData();
  },
  methods: {
    async loadData() {
      this.state.loading = true;
      if (this.user.village) {
        this.axios
          .get(
            `/transactions?offset=${this.offset}&limit=10000`
            // `/transactions?offset=${this.offset}&limit=${this.pagination.perPage}`
          )
          .then((res) => {
            let filteredItems = res.data.Transactions;
            filteredItems = filteredItems.filter(
              (item) => item.village === this.user.village
            );
            if (this.date.year) {
              filteredItems = filteredItems.filter(
                (item) =>
                  new Date(item.date_recorded).getFullYear() === this.date.year
              );
            }
            if (this.date.month) {
              filteredItems = filteredItems.filter(
                (item) =>
                  new Date(item.date_recorded).getMonth() + 1 == this.date.month
              );
            }
            this.table.items = filteredItems;
            this.totalAmount = 0;
            for (let i = 0; i < this.table.items.length; i++) {
              this.totalAmount += this.table.items[i].amount;
            }
          })
          .catch((err) => {
            console.log(err, err.response);
          })
          .finally(() => {
            this.state.loading = false;
          });
      }
    },
    closeSearch() {
      this.searchItem = "";
      this.state.search = false;
    },
  },
};
</script>

<style lang="scss">
.agent-payment-view {
  .search {
    border-radius: 2px;
    border: 1px solid #909090;
    color: #212121;
    max-width: 300px;
    &::placeholder {
      color: #212121;
    }
  }
  .button-group {
    button {
      margin: 0 0.5rem;
      border-radius: 3px;
      padding: 0.5rem 1.5rem;
    }
  }
  .fa-times {
    border-radius: 2px;
  }
  .momo-mtn-rw {
    text-align: center;
    border: 1px solid white;
    background: #f7c223;
    font-weight: 600;
    text-transform: uppercase;
  }
}
</style>
