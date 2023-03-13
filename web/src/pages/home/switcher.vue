<script>
/**
 * Switcher component
 */
export default {
  name: "Switcher",
  data() {
    return {
      isVisible: false
    };
  },
  methods: {
    // eslint-disable-next-line no-unused-vars
    setTheme(theme) {
      document
        .getElementById("color-opt")
        .setAttribute("href", "./css/colors/" + theme + ".css");
    },
    toggleSwitcher() {
      this.isVisible = !this.isVisible;
    }
  }
};
</script>
<template>
  <!-- Style switcher -->
  <div id="style-switcher" v-bind:style="isVisible ? 'left: 0px' : 'left: -189px' ">
    <div>
      <h3>Select your color</h3>
      <ul class="pattern">
        <li>
          <a class="color1" href="javascript: void(0);" v-on:click="setTheme('cyan')"></a>
        </li>
        <li>
          <a class="color2" href="javascript: void(0);" v-on:click="setTheme('red')"></a>
        </li>
        <li>
          <a class="color3" href="javascript: void(0);" v-on:click="setTheme('green')"></a>
        </li>
        <li>
          <a class="color4" href="javascript: void(0);" v-on:click="setTheme('pink')"></a>
        </li>
        <li>
          <a class="color5" href="javascript: void(0);" v-on:click="setTheme('blue')"></a>
        </li>
        <li>
          <a class="color6" href="javascript: void(0);" v-on:click="setTheme('purple')"></a>
        </li>
        <li>
          <a class="color7" href="javascript: void(0);" v-on:click="setTheme('orange')"></a>
        </li>
        <li>
          <a class="color8" href="javascript: void(0);" v-on:click="setTheme('yellow')"></a>
        </li>
      </ul>
    </div>
    <div class="bottom">
      <a href="javascript: void(0);" class="settings rounded-right" @click="toggleSwitcher">
        <i class="mdi mdi-settings mdi-spin"></i>
      </a>
    </div>
  </div>
  <!-- end Style switcher -->
</template>