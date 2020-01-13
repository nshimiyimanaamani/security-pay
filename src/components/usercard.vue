<template>
  <div>
    <b-card-header>
      <p>{{house.owner.fname +' '+ house.owner.lname}}</p>
      <p>{{house.id}}</p>
    </b-card-header>
    <b-card-footer>
      <article>
        <span v-show="house.percentage" class="completed">completed</span>
        <span class="details">{{house.due}} /12 last months</span>
        <b-progress :value="60" :max="100"></b-progress>
      </article>
      <i
        class="fa fa-ellipsis-v"
        :class="state.show ? null : 'collapsed'"
        :aria-expanded="state.show ? 'true' : 'false'"
        :aria-controls="''+index"
        @click="showCollapse()"
      ></i>
    </b-card-footer>
    <b-collapse :id="''+index" v-model="state.show" class="mt-3">
      <b-card no-body class="mb-1 border" v-for="(item,i) in availableMonths/12" :key="i">
        <b-card no-body class="border-0 rounded-0 border-bottom-1 p-0 m-0">
          <b-button
            block
            v-b-toggle="'accordion-'+i+house.id"
            variant="light"
            style="font-size: 15px"
          >Year - {{currentYear-i}}</b-button>
        </b-card>
        <b-collapse :id="'accordion-'+i+house.id" accordion="my-accordion" role="tabpanel">
          <b-card-body>
            <b-table
              id="user-table"
              striped
              hover
              small
              :items="invoices"
              :fields="fields"
              :busy="state.loading"
              :tbody-tr-class="rowClass"
              show-empty
            >
              <template v-slot:cell(id)="data">{{months[data.index]}}</template>
              <template v-slot:cell(amount)="data">
                {{data.item.status=="pending"?'-':''}}
                {{Number(data.item.amount).toLocaleString()}} Rwf
              </template>
              <template v-slot:table-busy>
                <div class="text-center my-2">
                  <b-spinner class="align-middle"></b-spinner>
                  <strong>Loading...</strong>
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
export default {
  name: "usercard",
  props: {
    house: Object,
    index: Number
  },
  data() {
    return {
      state: {
        show: false,
        loading: false
      },
      invoices: null,
      fields: [
        { key: "id", label: "Month" },
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
            this.endpoint +
              `/billing/invoices?property=${this.house.id}&months=${this.availableMonths}`
          )
          .then(res => {
            this.state.show = true;
            this.invoices = res.data.Invoices;
            this.state.loading = false;
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
            this.state.loading = false;
          });
      }
    }
  }
};
</script>

<style>
</style>