<template>
  <b-dropdown
    menu-class="bg-light w-100 font-13"
    size="sm"
    :id="'report-dropdown-form_'+random"
    ref="configureReports"
    no-caret
    :disabled.sync="disabled"
    dropright
    toggle-class="h-f-content border-0 my-2 py-2 app-color"
  >
    <template v-slot:button-content>
      <p class="m-0 font-13">{{title}}</p>
    </template>
    <b-dropdown-form form-class="p-2">
      <b-form-group label="Year" :label-for="'dropdown-year_'+random" @submit.stop.prevent>
        <b-form-select
          :id="'dropdown-year_'+random"
          v-model="object.year"
          class="bg-light"
          size="sm"
        >
          <option
            v-for="(year,i) in (currentYear-2019)"
            :value="currentYear-i"
            :key="id+''+year"
          >{{currentYear-i}}</option>
        </b-form-select>
      </b-form-group>

      <b-form-group label="Month" :label-for="'dropdown-month_'+random">
        <b-form-select
          :id="'dropdown-month_'+random"
          v-model="object.month"
          class="bg-light"
          size="sm"
        >
          <option v-for="i in 12" :value="i" :key="id+''+i">{{months[i-1]}}</option>
        </b-form-select>
      </b-form-group>
      <b-button variant="primary" class="w-100 app-color" size="sm" @click="handleOk">OK</b-button>
    </b-dropdown-form>
  </b-dropdown>
</template>

<script>
export default {
  props: {
    object: Object,
    id: String,
    title: String,
    disabled: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    currentYear() {
      return new Date().getFullYear();
    },
    currentMonth() {
      return new Date().getMonth() + 1;
    },
    months() {
      return this.$store.getters.getMonths;
    },
    random() {
      return Math.floor(Math.random() * 101);
    }
  },
  watch: {
    "object.year"() {
      handler: {
        if (this.currentYear < this.object.year) {
          this.$nextTick(() => {
            this.$set(this.object, "year", this.currentYear);
            this.object.month = 1;
          });
          padding;
        }
      }
    },
    "object.month"() {
      handler: {
        if (this.currentYear == this.object.year) {
          if (this.currentMonth < this.object.month) {
            this.$nextTick(() => {
              this.$set(this.object, "month", this.currentMonth);
            });
          }
        }
      }
    }
  },
  methods: {
    handleOk() {
      this.$refs.configureReports.hide(true);
      this.$emit("ok", this.object);
    }
  }
};
</script>

<style>
#config-dropdown-form > button {
  padding: 0;
  background: transparent !important;
  border: 0;
}
</style>