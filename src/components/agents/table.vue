<template>
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
    <template v-slot:cell(due)="data">{{ Number(data.item.due).toLocaleString() }} Rwf</template>
    <template v-slot:cell(owner)="data">
      {{
      data.item.owner.fname + " " + data.item.owner.lname
      }}
    </template>
    <template v-slot:table-busy>
      <div class="text-center my-2">
        <loader />
      </div>
    </template>
  </b-table>
</template>

<script>
import loader from "../loader.vue";
export default {
  name: "userTable",
  components: {
    loader
  },
  props: {
    user: Object,
    searchItem: String
  },
  data() {
    return {
      state: {
        loading: false,
        info: false
      },
      agent: null,
      table: {
        fields: [
          { key: "index", label: "NO" },
          { key: "owner", label: "Names", sortable: true },
          { key: "id", label: "House Code", sortable: false },
          { key: "owner.phone", label: "Phone Number", sortable: false },
          { key: "address.sector", label: "Sector", sortable: true },
          { key: "address.cell", label: "Cell", sortable: true },
          { key: "address.village", label: "Village", sortable: true },
          { key: "due", label: "Amount", sortable: false }
        ],
        items: null,
        sortBy: "owner",
        sortDesc: false
      }
    };
  },
  mounted() {
    this.loadData();
  },
  computed: {
    shownItems() {
      if (this.table.items) {
        return this.table.items.filter(item => {
          return (item.owner.fname + " " + item.owner.lname)
            .toLowerCase()
            .includes(this.searchItem.toLowerCase());
        });
      } else [];
    }
  },
  methods: {
    async loadData() {
      const agent = this.user;
      if (!agent) {
        this.$emit("getInfo");
      } else {
        this.state.loading = true;
        const total = await this.getTotal();
        this.axios
          .get(
            `/properties?village=${this.user.village}&offset=0&limit=${total}`
          )
          .then(res => {
            this.table.items = [];
            this.table.items = res.data.Properties.filter(
              item => item.address.cell == this.user.cell
            );
          })
          .catch(err => {
            const error = err.response
              ? err.response.data.error || err.response.data
              : null;
            if (error) this.$snotify.error(error);
          })
          .finally(() => {
            this.state.loading = false;
          });
      }
    },
    getTotal() {
      return this.axios
        .get(`/properties?village=${this.user.village}&offset=0&limit=10`)
        .then(res => {
          return res.data.Total;
        });
    }
  }
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
