<template>
  <b-container class="h-auto mw-100 agent-home">
    <b-row class="tab-buttons">
      <b-button
        variant="info"
        :class="{ active: state.selectedBtn == 1 }"
        @click="state.selectedBtn = 1"
        >Lists</b-button
      >
      <b-button
        variant="info"
        :class="{ active: state.selectedBtn == 2 }"
        @click="state.selectedBtn = 2"
        >Payments</b-button
      >
      <b-button
        variant="info"
        :class="{ active: state.selectedBtn == 3 }"
        @click="state.selectedBtn = 3"
        >UnPaid</b-button
      >
    </b-row>
    <hr />
    <transition-group name="fade" :duration="300">
      <div v-if="state.selectedBtn == 1" key="lists">
        <b-row class="flex-nowrap m-0 my-2">
          <controller :user="user" v-on:refresh="key++" />

          <b-button
            class="ml-2 d-flex align-items-center"
            variant="info"
            @click="state.search = true"
          >
            search
            <i class="fa fa-search ml-1" />
          </b-button>
        </b-row>
        <transition name="fade" :duration="500">
          <b-row class="w-100 m-0 my-3 search" v-if="state.search">
            <input
              type="search"
              id="agent-search-user"
              placeholder="keyword to search..."
              v-model="searchItem"
            />
            <i class="fa fa-times" @click="closeSearch" />
          </b-row>
        </transition>
        <b-row class="m-0">
          <user-table
            v-if="user"
            :user="user"
            :key="key"
            :searchItem="searchItem"
          />
        </b-row>
      </div>
      <div v-if="state.selectedBtn == 2" key="payment">
        <agent-payment-view v-if="user" :user="user" :key="key" />
      </div>
      <div v-if="state.selectedBtn == 3" key="unpaid">
        <agent-unpaid-view v-if="user" :user="user" :key="key" />
      </div>
    </transition-group>
  </b-container>
</template>

<script>
export default {
  name: "home",
  components: {
    "user-table": () => import("./table.vue"),
    controller: () => import("./controllers.vue"),
    "agent-payment-view": () => import("./agentPaymentView"),
    "agent-unpaid-view": () => import("./agentUnpaidView"),
  },
  data() {
    return {
      searchItem: "",
      state: {
        selectedBtn: 1,
        search: false,
      },
      user: null,
      key: 0,
    };
  },
  computed: {
    userDetails() {
      return this.$store.getters.userDetails;
    },
  },
  mounted() {
    this.getInfo();
  },
  methods: {
    getInfo() {
      this.axios
        .get("/accounts/agents/" + this.userDetails.username)
        .then((res) => {
          this.user = res.data;
        })
        .catch((err) => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        });
    },
    closeSearch() {
      this.state.search = false;
      this.searchItem = "";
    },
  },
};
</script>

<style lang="scss">
.agent-home {
  .search {
    input {
      flex: 1;
      border-radius: 3px;
      border: 1px solid #d6d6d6;
      color: #233a44;
      padding: 0.5rem;
      margin-right: 0.5rem;
      &::placeholder {
        color: #717679;
      }
    }
    i {
      background: #017db3;
      color: #ffffff;
      display: flex;
      justify-content: center;
      -webkit-box-align: center;
      align-items: center;
      width: 2.5rem;
      font-size: 1rem;
      border-radius: 3px;
      cursor: pointer;
    }
  }

  hr {
    margin-top: 0;
  }
  .tab-buttons {
    display: flex;
    flex-wrap: nowrap;
    justify-content: center;
    align-items: center;
    padding: 0.5rem;

    button {
      margin: 0 0.5rem;
      background: white !important;
      color: #017db3 !important;
      &.active {
        background: #017db3 !important;
        color: white !important;
      }
    }
  }
  .table-control {
    overflow: auto;
    width: 100%;
  }
}
</style>
