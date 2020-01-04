<template>
  <b-table small striped bordered hover :items="table.items" :fields="table.fields">
    <template v-slot:cell(due)="data">{{Number(data.item.due).toLocaleString()}} Rwf</template>
  </b-table>
</template>

<script>
export default {
  name: "userTable",
  data() {
    return {
      table: {
        fields: [
          { key: "name", label: "Names", sortable: true },
          { key: "id", label: "House Code", sortable: false },
          { key: "phone", label: "Phone Number", sortable: false },
          { key: "sector", label: "Sector", sortable: true },
          { key: "cell", label: "Cell", sortable: true },
          { key: "village", label: "Village", sortable: true },
          { key: "due", label: "Amount", sortable: false }
        ],
        items: [
          {
            name: "shami martin",
            id: "5FDDFE3",
            phone: "123456",
            sector: "remera",
            cell: "nyabisindu",
            village: "nyabisindu",
            due: "1000"
          }
        ]
      }
    };
  }
};
</script>

<style lang="scss">
.table {
  min-width: max-content;
  thead,
  tbody {
    font-size: 14px;
    min-width: fit-content;
  }
}
</style>