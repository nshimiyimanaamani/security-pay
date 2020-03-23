<template>
  <div class="message-container px-5 py-3">
    <h4 align="center" class="align-self-center text-uppercase mb-4">Send Text Messages</h4>
    <div class="controls">
      <b-form-radio-group v-model="selected" :options="options" name="message-controls"></b-form-radio-group>
    </div>
    <div class="message-select">
      <b-select v-model="select.cell" :options="cellOptions" v-if="showCell || showVillage">
        <template v-slot:first>
          <option :value="null" disabled>Select cell</option>
        </template>
      </b-select>
      <b-select v-model="select.village" :options="villageOptions" v-if="showVillage">
        <template v-slot:first>
          <option
            :value="null"
            disabled
          >{{villageOptions.length?'Select village':'Select cell first'}}</option>
        </template>
      </b-select>
    </div>
    <div class="text-wrapper w-100">
      <textarea
        name="message"
        id="message"
        placeholder="Write your message..."
        cols="30"
        rows="10"
        class="message-textarea"
        v-model="message"
      />
      <b-button
        variant="info"
        :disabled="message && selected?false:true"
        class="float-right"
        @click="send"
      >
        {{state.sending ? 'Sending ':'Send '}}
        <b-spinner v-if="state.sending" variant="black" small />
        <i v-else class="fa fa-paper-plane" />
      </b-button>
    </div>
  </div>
</template>

<script>
import { Village } from "rwanda";
export default {
  data() {
    return {
      message: null,
      selected: null,
      state: {
        sending: false
      },
      select: {
        cell: null,
        village: null
      },
      options: [
        { text: "Send to sector", value: "sector" },
        { text: "Send to cell", value: "cell" },
        { text: "Send to village", value: "village" }
      ]
    };
  },
  computed: {
    showVillage() {
      return this.selected == "village" ? true : false;
    },
    showCell() {
      return this.selected == "cell" ? true : false;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      const cell = this.select.cell;
      if (cell) {
        return Village("Kigali", "Gasabo", "Remera", cell);
      } else {
        return [];
      }
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    }
  },
  methods: {
    send() {
      const selected = this.selected;
      if (selected == "sector") this.sendToSector();
      else if (selected == "cell") this.sendToCell();
      else if (selected == "village") this.sendToVillage();
    },
    async sendToSector() {
      this.state.sending = true;
      let request = `/properties?sector=${this.activeSector}&offset=0&limit=`;
      let message = this.message;
      let recipients = await this.getPhoneArray(request);
      if (recipients) {
        this.axios
          .get(`/notifications/send`, {
            message: message,
            recipients: recipients
          })
          .then(res => {
            this.$snotify.info(
              `Message sent to all properties in ${this.activeSector} sector! Message sent successfully`
            );
          })
          .catch(err => {
            console.log(err.response);
            this.$snotify.error("Unable to send message! Message not sent");
          })
          .finally(() => {
            this.state.sending = false;
          });
      } else {
        this.$snotify.error("An error occured");
      }
    },
    async sendToCell() {
      this.state.sending = true;
      let cell = this.select.cell;
      let request = `/properties?cell=${cell}&offset=0&limit=`;
      let message = this.message;
      let recipients = await this.getPhoneArray(request);
      if (recipients) {
        this.axios
          .get(`/notifications/send`, {
            message: message,
            recipients: recipients
          })
          .then(res => {
            this.$snotify.info(
              `Message sent to all properties in ${cell} cell! Message sent successfully`
            );
          })
          .catch(err => {
            console.log(err.response);
            this.$snotify.error("Unable to send message! Message not sent");
          })
          .finally(() => {
            this.state.sending = false;
          });
      } else {
        this.$snotify.error("An error occured!");
      }
    },
    async sendToVillage() {
      this.state.sending = true;
      let village = this.select.village;
      let request = `/properties?cell=${village}&offset=0&limit=`;
      let message = this.message;
      let recipients = await this.getPhoneArray(request);
      if (recipients) {
        this.axios
          .get(`/notifications/send`, {
            message: message,
            recipients: recipients
          })
          .then(res => {
            this.$snotify.info(
              `Message sent to all properties in ${village} village! Message sent successfully`
            );
          })
          .catch(err => {
            console.log(err.response);
            this.$snotify.error("Unable to send message! Message not sent");
          })
          .finally(() => {
            this.state.sending = false;
          });
      } else {
        this.$snotify.error("An error occured!");
      }
    },
    async getPhoneArray(request) {
      let total = await this.getTotal(request + "0");
      return this.axios
        .get(request + `${total}`)
        .then(res => res.data.Properties.map(item => item.owner.phone))
        .catch(err => null);
    },
    getTotal(request) {
      return this.axios
        .get(request)
        .then(res => res.data.Total)
        .catch(err => null);
    }
  }
};
</script>

<style lang="scss">
@import "../assets/css/messages.scss";
</style>