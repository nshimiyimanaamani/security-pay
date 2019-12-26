const state = {
  token: sessionStorage.token ?
    sessionStorage.token : null,
  endpoint: process.env.VUE_APP_PAYPACK_API,
  user: null,
  active_sector: "remera",
  active_cell: "",
  cells_array: [],
  active_village: "",
  village_array: [],
  villageByCell: []
}
export default state
