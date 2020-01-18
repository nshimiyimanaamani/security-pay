<template>
  <b-table
    small
    striped
    bordered
    hover
    responsive
    show-empty
    :busy="state.loading"
    :items="table.items"
    :fields="table.fields"
  >
    <template v-slot:cell(due)="data">{{Number(data.item.due).toLocaleString()}} Rwf</template>
    <template v-slot:cell(owner)="data">{{data.item.owner.fname +" "+ data.item.owner.lname}}</template>
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
    user: Object
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
          { key: "owner", label: "Names", sortable: true },
          { key: "id", label: "House Code", sortable: false },
          { key: "owner.phone", label: "Phone Number", sortable: false },
          { key: "address.sector", label: "Sector", sortable: true },
          { key: "address.cell", label: "Cell", sortable: true },
          { key: "address.village", label: "Village", sortable: true },
          { key: "due", label: "Amount", sortable: false }
        ],
        items: null
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    userDetails() {
      return this.$store.getters.userDetails;
    }
  },
  created() {
    this.loadData();
  },
  methods: {
    loadData() {
      const agent = this.user;
      if (!agent) {
        this.$emit("getInfo");
      } else {
        this.state.loading = true;
        this.axios
          .get(
            this.endpoint +
              `/properties?sector=${agent.sector}&offset=0&limit=1000`
          )
          .then(res => {
            this.table.items = [];
            this.table.items = res.data.Properties.filter(this.isForAgent);
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
            this.state.loading = false;
          });
      }
    },
    isForAgent(value) {
      return value.recorded_by == this.userDetails.username;
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