<template>
  <b-table small striped bordered hover responsive show-empty :busy="state.loading" :items="shownItems"
    :fields="table.fields" :sort-by.sync="table.sortBy">
    <template v-slot:cell(index)="data">
      <article class="text-center">{{ data.index + 1 }}</article>
    </template>
    <template v-slot:cell(due)="data">{{ Number(data.item.due).toLocaleString() }} Rwf</template>
    <template v-slot:cell(owner)="data">
      {{ data.item.owner.fname + " " + data.item.owner.lname }}
    </template>
    <template v-slot:table-busy>
      <div class="text-center my-2">
        <loader />
      </div>
    </template>
    <template v-slot:custom-foot>
      <b-tr class="total">
        <b-td colspan="10">
          <b-pagination class="my-0" align="center" :per-page="pagination.perPage" v-model="pagination.page"
            :total-rows="pagination.total" @input="loadData"></b-pagination>
        </b-td>
      </b-tr>
    </template>
  </b-table>
</template>

<script>
import loader from "../loader.vue";
export default {
  name: "userTable",
  components: {
    loader,
  },
  props: {
    user: Object,
    searchItem: String,
  },
  data() {
    return {
      state: {
        loading: false,
        info: false,
      },
      table: {
        fields: [
          { key: "index", label: "NO" },
          { key: "owner", label: "Names", sortable: true },
          { key: "id", label: "House Code", sortable: false },
          { key: "owner.phone", label: "Phone Number", sortable: false },
          { key: "address.sector", label: "Sector", sortable: true },
          { key: "address.cell", label: "Cell", sortable: true },
          { key: "address.village", label: "Village", sortable: true },
          { key: "due", label: "Amount", sortable: false },
        ],
        items: null,
        sortBy: "owner",
        sortDesc: false,
      },
      pagination: {
        perPage: 20,
        total: 0,
        page: 1,
      },
    };
  },
  mounted() {
    this.loadData();
  },
  computed: {
    shownItems() {
      if (this.table.items) {
        return this.table.items.filter((item) => {
          return (item.owner.fname + " " + item.owner.lname)
            .toLowerCase()
            .includes(this.searchItem.toLowerCase());
        });
      } else[];
    },
    offset() {
      return (this.pagination.page - 1) * this.pagination.perPage;
    },
    showPagination() {
      if (this.state.loading) return false;
      if (this.pagination.total / this.pagination.perPage < 2) return false;
      return true;
    },
  },
  methods: {
    async loadData() {
      this.state.loading = true;
      this.axios
        .get(
          `/properties?village=${this.user.village}&offset=${this.offset}&limit=${this.pagination.perPage}&names=`
        )
        .then((res) => {
          this.table.items = res.data.Properties.filter(
            (item) => item.address.cell == this.user.cell
          );
          this.pagination.total = res.data.Total;
        })
        .catch((err) => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.state.loading = false;
        });
    },
  },
};
</script>

<style lang="scss">
.table {
  min-width: max-content;

  thead,
  tbody {
    font-size: 14px;
    min-width: fit-content;
  }
}
</style>
