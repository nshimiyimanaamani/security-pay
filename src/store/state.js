const state = {
  endpoint: process.env.VUE_APP_PAYPACK_API,
  user: null,
  active_sector: "Remera",
  active_cell: null,
  cells_array: null,
  active_village: null,
  village_array: null,
  villageByCell: null,
  months: [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December"
  ]
};
export default state;
