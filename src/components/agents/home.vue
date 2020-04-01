<template>
  <b-container class="h-auto mw-100 agent-home">
    <b-row class="pt-3 px-3 justify-content-center">
      <b-button-group class="agent-nav-btn">
        <b-button
          variant="info"
          class="px-4"
          :class="{'active':state.btn1}"
          @click="toggleClass('btn1')"
        >Lists</b-button>
        <b-button
          variant="info"
          class="px-4"
          :class="{'active':state.btn2}"
          @click="toggleClass('btn2')"
        >Payments</b-button>
      </b-button-group>
    </b-row>
    <hr />
    <transition-group name="fade" :duration="300">
      <div v-show="state.btn1" key="lists">
        <b-row class="px-4 flex-nowrap">
          <controller :user="user" v-on:refresh="key++" />
          <b-button size="sm" class="my-2" variant="info" @click="state.search = true">search</b-button>
        </b-row>
        <transition name="fade" :duration="500">
          <b-row
            class="flex-row justify-content-end w-100 m-0 flex-nowrap align-items-center px-2"
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
        <b-row class="px-4">
          <user-table v-on:getInfo="getInfo" :user="user" :key="key" :searchItem="searchItem" />
        </b-row>
      </div>
      <div v-show="state.btn2" key="payment">
        <agent-payment-view :user="user" />
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
    "agent-payment-view": () => import("./agentPaymentView")
  },
  data() {
    return {
      searchItem: "",
      state: {
        btn1: true,
        btn2: false,
        search: false
      },
      user: null,
      key: 0
    };
  },
  computed: {
    userDetails() {
      return this.$store.getters.userDetails;
    }
  },
  created() {
    this.getInfo();
  },
  methods: {
    getInfo() {
      this.axios
        .get("/accounts/agents/" + this.userDetails.username)
        .then(res => {
          this.user = { ...res.data };
          this.key++;
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        });
    },
    toggleClass(key) {
      if (this.state[key] === false) {
        Object.keys(this.state)
          .filter(item => item !== key)
          .map(item => {
            this.state[item] = !this.state[item];
          });

        setTimeout(() => {
          this.state[key] = !this.state[key];
        }, 300);
      }
    },
    closeSearch() {
      this.state.search = false;
      this.searchItem = "";
    }
  }
};
</script>

<style lang='scss'>
.agent-home {
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
  .agent-nav-btn {
    .btn {
      border: 0;
      border-left: 1.5px solid white !important ;
      border-right: 1.5px solid white !important ;
      border-radius: 15px !important;
      &.active {
        background: #37505a !important;
      }
    }
  }
  .table-control {
    overflow: auto;
    width: 100%;
  }
}
</style>