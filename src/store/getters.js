const getters = {
  getEndpoint: state => state.endpoint,
  getSectorArray: state => state.sector,
  getCellsArray: state => state.cells_array,
  getActiveCell: state => state.active_cell,
  getActiveVillage: state => state.active_village,
  getVillageArray: state => state.village_array,
  getActiveSector: state => state.active_sector,
  userDetails: state => state.user,
  getMonths: state => state.months
}

export default getters
