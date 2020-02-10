<template>
  <div>
    <b-card-header class="d-flex bg-white justify-content-between py-1">
      <p>{{house.owner.fname +' '+ house.owner.lname}}</p>
      <p class="d-block text-truncate">{{house.id}}</p>
    </b-card-header>
    <b-card-footer class="d-flex justify-content-between py-2 align-items-end bg-white">
      <article class="d-flex flex-column w-100 pr-1">
        <span v-show="house.percentage" class="completed font-13">completed</span>
        <span class="details font-13">{{house.due}} Rwf per month</span>
        <b-progress class="w-100" :value="0" :max="100"></b-progress>
      </article>
      <i
        class="fa fa-ellipsis-v cursor-pointer"
        :class="state.show ? null : 'collapsed'"
        :aria-expanded="state.show ? 'true' : 'false'"
        :aria-controls="''+index"
        @click="showCollapse()"
      ></i>
    </b-card-footer>
    <b-collapse :id="''+index" v-model="state.show" class="my-2 px-2">
      <b-card no-body class="mb-1 border" v-for="(item,i) in availableMonths/12" :key="i">
        <b-card no-body class="bg-white border-0 rounded-0 border-bottom-1 p-0 m-0">
          <b-button
            block
            v-b-toggle="'accordion-'+i+house.id"
            variant="light"
            class="font-13"
          >Year - {{currentYear-i}}</b-button>
        </b-card>
        <b-collapse :id="'accordion-'+i+house.id" accordion="my-accordion" role="tabpanel">
          <b-card-body>
            <b-table
              id="user-table"
              striped
              hover
              small
              responsive
              :items="invoices"
              :fields="fields"
              :busy="state.loading"
              :tbody-tr-class="rowClass"
              show-empty
            >
              <template
                v-slot:cell(created_at)="data"
              >{{new Date(data.value).toLocaleString('en-EN', {month: 'long'})}}</template>
              <template v-slot:cell(amount)="data">
                {{data.item.status=="pending"?'-':''}}
                {{Number(data.item.amount).toLocaleString()}} Rwf
              </template>
              <template v-slot:table-busy>
                <div class="text-center my-2">
                  <loader />
                </div>
              </template>
            </b-table>
          </b-card-body>
        </b-collapse>
      </b-card>
    </b-collapse>
  </div>
</template>

<script>
import loader from "../components/loader.vue";
export default {
  name: "usercard",
  props: {
    house: Object,
    index: Number
  },
  components: {
    loader
  },
  data() {
    return {
      state: {
        show: false,
        loading: false
      },
      invoices: null,
      fields: [
        { key: "created_at", label: "Month" },
        { key: "property", label: "House ID" },
        { key: "status", label: "Status" },
        { key: "amount", label: "Amount" }
      ]
    };
  },
  mounted() {
    this.$root.$on("bv::collapse::state", (collapseId, isJustShown) => {
      if (collapseId.includes("accordion") && isJustShown) {
        this.loadData();
      }
    });
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    months() {
      return this.$store.getters.getMonths;
    },
    currentYear() {
      return new Date().getFullYear();
    },
    availableMonths() {
      const createdYear = new Date(this.house.created_at);
      return 12 * (this.currentYear - createdYear.getFullYear() + 1);
    }
  },
  methods: {
    showCollapse() {
      this.state.show = !this.state.show;
    },
    rowClass(item, type) {
      if (!item || type !== "row") return;
      if (item.status === "pending") return "table-danger";
    },
    loadData() {
      if (this.state.show) {
        this.state.loading = true;
        this.axios
          .get(
            `/billing/invoices?property=${this.house.id}&months=${this.availableMonths}`
          )
          .then(res => {
            this.state.show = true;
            this.invoices = res.data.Invoices;
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
    }
  }
};
</script>

<style>
</style>