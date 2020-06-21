<template>
  <b-dropdown
    menu-class="bg-light"
    :id="'config-dropdown-form_'+random"
    ref="configure"
    no-caret
    right
    toggle-class="m-0 border-0 bg-transparent"
  >
    <template v-slot:button-content>
      <p class="m-0 d-flex align-items-center">
        settings
        <i class="fa fa-cog text-white fsize-xl ml-1" />
      </p>
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
          :id="'dropdown-month'+random"
          v-model="object.month"
          class="bg-light"
          size="sm"
        >
          <option :value="null">All</option>
          <option v-for="i in 12" :value="i" :key="id+''+i">{{months[i-1]}}</option>
        </b-form-select>
      </b-form-group>
      <b-button variant="primary" class="w-100" size="sm" @click="handleOk">OK</b-button>
    </b-dropdown-form>
  </b-dropdown>
</template>

<script>
export default {
  props: {
    object: Object,
    id: String
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
            this.object.month = null;
          });
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
      this.$refs.configure.hide(true);
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