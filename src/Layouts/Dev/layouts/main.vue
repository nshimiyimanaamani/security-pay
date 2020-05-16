<template>
  <div class="dev-dashboard d-flex">
    <div class="dev-sidebar app-color" :class="{'active':active}">
      <h1>
        paypack
        <code class="font-13 text-white">{dev}</code>
      </h1>

      <ul class="dev-sidebar-links">
        <router-link
          to="/dev"
          class="text-white text-decoration-none"
          exact-active-class="selected"
        >
          <li class="text-uppercase">statistics</li>
        </router-link>
        <router-link
          to="/dev/account"
          class="text-white text-decoration-none"
          exact-active-class="selected"
        >
          <li class="text-uppercase">accounts</li>
        </router-link>
      </ul>
    </div>
    <div class="dev-content" :class="{'active':active}">
      <nav class="dev-nav navbar navbar-expand-lg navbar-light bg-light border-bottom">
        <b-button variant="info" @click="active =!active">
          <i class="fa fa-align-left" />
        </b-button>
        <b-button variant="info" @click="logout">Logout</b-button>
      </nav>
      <b-container class="dev-body">
        <router-view />
      </b-container>
    </div>
  </div>
</template>

<script>
export default {
  name: "dev-dashboard",
  data() {
    return {
      active: false
    };
  },
  methods: {
    logout() {
      this.$store.dispatch("logout");
      this.$destroy;
    }
  }
};
</script>

<style lang="scss">
.dev-dashboard {
  width: 100vw;
  .dev-sidebar {
    min-width: 200px;
    max-width: 200px;
    transition: all 0.5s;
    display: flex;
    flex-direction: column;

    &.active {
      margin-left: -200px;
    }
    h1 {
      font-size: 27px;
      height: 60px;
      letter-spacing: 5px;
      white-space: nowrap;
      display: flex;
      user-select: none;
      justify-content: center;
      align-items: center;
      text-transform: uppercase;
      color: white;
      border-bottom: 1px solid #dee2e6;
    }
    .dev-sidebar-links {
      list-style: none;
      padding: 0;
      width: 100%;
      margin: 0;
      margin-top: 3rem;
      .selected li {
        background: white;
        color: #017db3;
      }

      li {
        padding: 1rem;
        margin: 1px 0;

        &:hover {
          background: white;
          color: #017db3;
        }
      }
    }
  }
  .dev-content {
    width: 100%;
    min-height: 100vh;
    transition: all 0.5s;

    nav {
      height: 60px;
      max-width: calc(100vw - 200px);
      transition: all 0.5s;
      display: flex;
      justify-content: space-between;

      button {
        border-radius: 3px;
      }
    }
    .dev-body {
      height: calc(100vh - 60px);
      max-width: calc(100vw - 200px);
      transition: all 0.5s;
    }
    &.active {
      .dev-body {
        max-width: 100vw;
      }
      nav {
        max-width: 100vw;
      }
    }
  }
}
</style>