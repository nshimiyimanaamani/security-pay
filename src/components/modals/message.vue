<template>
  <b-modal @hidden="modalClosed" v-model="show" size="sm" class="mt-3" hide-footer>
    <template v-slot:modal-title>Send Text Message</template>
    <div class="message-wrapper px-2">
      <b-row class="align-items-baseline">
        <p>Phone Number:</p>
        <b-badge v-for="(phone,i) in [phones]" :key="i" class="align-items-baseline d-flex p-1 m-1">
          <b-card-text class="m-0">{{phone}}</b-card-text>
          <!-- <b-button size="sm" class="ml-2 p-1 rounded-circle bg-white close-button">
            <i class="fa fa-times text-black-50" />
          </b-button>-->
        </b-badge>
      </b-row>
      <b-row class="mt-3">
        <textarea
          v-model="message"
          class="rounded w-100"
          rows="5"
          placeholder="Enter your message..."
        ></textarea>
      </b-row>
      <b-row>
        <b-button
          @click="send"
          size="sm"
          class="mt-3 ml-auto"
          variant="info"
          :disabled="message? false : true"
        >
          Send
          <i class="fa fa-paper-plane" />
        </b-button>
      </b-row>
    </div>
  </b-modal>
</template>

<script>
export default {
  name: "message-component",
  props: {
    phones: String
  },
  data() {
    return {
      show: true,
      message: null
    };
  },
  methods: {
    send() {
      this.$emit("sent");
      alert("message feature is not yet finished!");
    },
    modalClosed() {
      this.$emit("modal-closed");
      this.message = null;
    }
  }
};
</script>

<style lang="scss">
.message-wrapper {
  .close-button {
    width: 2em;
    height: 2em;
  }
}
</style>