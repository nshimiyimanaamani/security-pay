<template>
  <div class="agent-payment-view">
    <b-row class="justify-content-between px-3">
      <b-button size="sm" variant="info" @click="loadData">Refresh</b-button>
      <b-button-group>
        <b-button size="sm" variant="info" @click="state.search = true">search</b-button>
        <b-button size="sm" variant="info">
          settings
          <selector :object="date" :id="'agent-date-selector'" @ok="loadData" />
        </b-button>
      </b-button-group>
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
        <i class="fa fa-times ml-2 p-2 app-color text-white" @click="closeSearch" />
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
          <article class="text-center">{{ data.index + 1}}</article>
        </template>
        <template v-slot:empty>
          <article class="text-center p-3">No payment records found.</article>
        </template>
        <template v-slot:cell(amount)="data">{{ Number(data.item.amount).toLocaleString() }} Rwf</template>
        <template v-slot:cell(month)="data">{{months[new Date(data.item.date_recorded).getMonth()]}}</template>
        <template v-slot:cell(year)="data">{{new Date(data.item.date_recorded).getFullYear()}}</template>
        <template v-slot:cell(owner)="data">
          {{
          data.item.owner_firstname+ " " + data.item.owner_lastname
          }}
        </template>
        <template v-slot:cell(method)="data">
          <article :class="data.value">{{ data.value.includes('mtn')?'mtn':data.value }}</article>
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
  name: "agentPaymentView",
  components: {
    loader: () => import("../loader"),
    selector: () => import("../yearSelector")
  },
  props: {
    user: Object
  },
  data() {
    return {
      searchItem: "",
      state: {
        search: false,
        loading: false
      },
      date: {
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1
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
          { key: "amount", label: "Amount", sortable: false }
        ],
        sortBy: "owner"
      }
    };
  },
  computed: {
    months() {
      return this.$store.getters.getMonths;
    },
    shownItems() {
      if (this.table.items) {
        return this.table.items.filter(item => {
          return (item.owner_firstname + " " + item.owner_lastname)
            .toLowerCase()
            .includes(this.searchItem.toLowerCase());
        });
      } else [];
    }
  },
  mounted() {
    this.state.loader = false;
    console.log(this.user);
    this.loadData();
  },
  methods: {
    async loadData() {
      this.state.loading = true;
      const url = "/transactions?offset=0&limit=";
      let total = await this.getTotal(url + "0");
      if (total && this.user.village) {
        this.axios
          .get(url + total)
          .then(res => {
            this.table.items = res.data.Transactions.filter(item => {
              return (
                item.village === "Gishushu" &&
                new Date(item.date_recorded).getFullYear() === this.date.year &&
                new Date(item.date_recorded).getMonth() + 1 === this.date.month
              );
            });
            console.log(res.data);
          })
          .catch(err => {
            console.log(err, err.response);
          })
          .finally(() => {
            this.state.loading = false;
          });
      }
    },
    getTotal(url) {
      return this.axios
        .get(url)
        .then(res => res.data.Total)
        .catch(err => null);
    },
    closeSearch(){
      this.searchItem=''
      this.state.search = false
    }
  }
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