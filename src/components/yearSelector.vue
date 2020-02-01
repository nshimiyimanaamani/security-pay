<template>
  <b-dropdown
    menu-class="bg-light"
    id="config-dropdown-form"
    size="sm"
    ref="configure"
    no-caret
    right
  >
    <template v-slot:button-content>
      <i
        class="fa fa-cog text-white font-20"
        @click="object.configuring=true"
        :class="{'fa-spin':object.configuring}"
      />
    </template>
    <b-dropdown-form form-class="p-2">
      <b-form-group label="Year" label-for="dropdown-year" @submit.stop.prevent>
        <b-form-select id="dropdown-year" v-model="object.year" class="bg-light" size="sm">
          <option
            v-for="(year,i) in (currentYear-2019)"
            :value="currentYear-i"
            :key="id+''+year"
          >{{currentYear-i}}</option>
        </b-form-select>
      </b-form-group>

      <b-form-group label="Month" label-for="dropdown-month">
        <b-form-select id="dropdown-month" v-model="object.month" class="bg-light" size="sm">
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
    }
  },
  watch: {
    "object.config.year"() {
      handler: {
        if (this.currentYear < this.object.year) {
          this.$nextTick(() => {
            this.$set(this.object, "year", this.currentYear);
            this.object.month = 1;
          });
        }
      }
    },
    "object.config.month"() {
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
  mounted() {
    this.$root.$on("bv::dropdown::hide", bvEvent => {
      if (bvEvent.componentId == "config-dropdown-form") {
        this.object.configuring = false;
      }
    });
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