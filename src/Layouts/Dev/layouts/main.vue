<template>
  <div class="dev-dashboard">
    <vue-title title="Paypack | Dev" />
    <nav class="dev-nav">
      <h1 class="logo">
        paypack
        <code>{dev}</code>
      </h1>
      <ul class="dev-links">
        <router-link to="/dev" exact-active-class="selected">
          <li>Accounts</li>
        </router-link>
        <router-link to="/dev/developers" exact-active-class="selected">
          <li>Developers</li>
        </router-link>
        <router-link to="/dev/admins" exact-active-class="selected">
          <li>Admins</li>
        </router-link>
        <router-link to="/dev/managers" exact-active-class="selected">
          <li>Managers</li>
        </router-link>
      </ul>
      <b-button class="custom-button" @click="logout">Logout</b-button>
    </nav>
    <div class="dev-body">
      <b-container>
        <router-view />
      </b-container>
    </div>
  </div>
</template>

<script>
export default {
  name: "dev-dashboard",
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
  width: 100%;
  display: grid;
  grid-template-rows: 60px minmax(calc(100vh - 60px), auto);
  min-width: 880px;

  nav {
    height: 60px;
    display: flex;
    -webkit-box-pack: justify;
    justify-content: space-between;
    background: #0382b9;
    color: white;
    align-items: center;
    padding: 0 1rem;
    position: sticky;
    top: 0;
    box-shadow: 0 1px 5px 0 rgba(32, 33, 36, 0.32);

    .logo {
      text-transform: uppercase;
      margin: 0;
      font-family: "Montserrat", Arial, sans-serif;
      font-size: 2.2rem;
      display: flex;
      flex-wrap: nowrap;
      code {
        color: #bdf3ff;
        font-size: 1.2rem;
      }
    }

    .dev-links {
      margin: 0 2rem;
      padding: 0;
      display: flex;
      flex-direction: row;
      flex-wrap: nowrap;
      flex: 1;
      justify-content: center;
      align-items: flex-end;
      align-self: flex-end;

      a {
        list-style: none;
        margin: 0 1px;
        color: white;
        text-transform: uppercase;
        font-family: "Montserrat", Arial, Helvetica, -apple-system,
          BlinkMacSystemFont, "Helvetica Neue", "Noto Sans", sans-serif;
        font-weight: 500;
        font-size: 1rem;
        padding: 0.5rem 1.5rem;
        text-decoration: none;
        border-radius: 3px 3px 0 0;

        &.selected,
        &:hover {
          background: white;
          color: #0382b9;
        }
      }
    }

    .custom-button {
      border-radius: 3px;
      height: fit-content;
      background: transparent;
      color: white;
      border: 1px solid white;
    }
  }
  .dev-body {
    width: 100%;
    height: calc(100vh - 60px);
    padding: 0;
    margin: 0;
    overflow: auto;
  }
}
</style>